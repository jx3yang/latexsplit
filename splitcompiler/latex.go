package splitcompiler

import (
	"bufio"
	"path/filepath"

	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jx3yang/latexsplit/argparse"
	pdfcup "github.com/pdfcpu/pdfcpu/pkg/api"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(length int) string {
	r := make([]rune, length)
	for i := range r {
		r[i] = letters[rand.Intn(len(letters))]
	}
	return string(r)
}

func getRandFileId(length int) string {
	return randString(length) + "_latexsplit"
}

func genNewFileContent(fileName, splitLine, id string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		currLine := scanner.Text()
		if currLine == splitLine {
			lines = append(lines, fmt.Sprintf("\\typeout{%v \\thepage}", id))
		} else {
			lines = append(lines, currLine)
		}
	}
	return lines, nil
}

func compile(lines []string, compiler, args string) error {
	cmd := exec.Command(compiler, args)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, strings.Join(lines, "\n"))
	}()
	return cmd.Run()
}

func parseLogFile(logFile, id string) ([]int, error) {
	f, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var pageNums []int
	pattern := fmt.Sprintf("^%v \\d+$", id)
	for scanner.Scan() {
		b := scanner.Bytes()
		if matched, _ := regexp.Match(pattern, b); matched {
			i, _ := strconv.Atoi(string(b)[len(id)+1:])
			pageNums = append(pageNums, i)
		}
	}
	return pageNums, nil
}

func splitPDF(fileStem, id string) (string, error) {
	tempDir, err := ioutil.TempDir("", id)
	if err != nil {
		return "", err
	}
	os.Mkdir(tempDir, 0755)
	err = pdfcup.SplitFile(fileStem+".pdf", tempDir, 1, nil)
	if err != nil {
		return "", err
	}
	return tempDir, nil
}

func mergePDF(dir, fileStem string, start, end, part int) (string, error) {
	outFile := fmt.Sprintf("%v_%v.pdf", part, fileStem)
	inFiles := make([]string, end-start+1)
	for i := start; i <= end; i++ {
		inFiles[i-start] = filepath.Join(dir, fmt.Sprintf("%v_%v.pdf", fileStem, i))
	}
	err := pdfcup.MergeCreateFile(inFiles, outFile, nil)
	return outFile, err
}

func worker(wg *sync.WaitGroup, results chan struct {
	string
	error
}, dir, fileStem string, start, end, part int) {
	defer wg.Done()

	outFile, err := mergePDF(dir, fileStem, start, end, part)
	results <- struct {
		string
		error
	}{outFile, err}
}

func LatexSplitCompiler(args *argparse.Arguments) ([]string, error) {
	rand.Seed(time.Now().UnixNano())
	id := getRandFileId(8)

	lines, err := genNewFileContent(args.FileName, args.SplitLine, id)
	if err != nil {
		return nil, err
	}

	err = compile(lines, args.Compiler, args.CompilerArgs)
	if err != nil {
		return nil, err
	}

	logFile := args.FileStem + ".log"
	pageNums, err := parseLogFile(logFile, id)
	if err != nil {
		return nil, err
	}

	numPages, err := pdfcup.PageCountFile(args.FileStem + ".pdf")
	if err != nil {
		return nil, err
	}
	pageNums = append([]int{0}, pageNums...)
	pageNums = append(pageNums, numPages)

	tempDir, err := splitPDF(args.FileStem, id)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		return nil, err
	}

	results := make(chan struct {
		string
		error
	})

	wg := new(sync.WaitGroup)
	wg.Add(len(pageNums) - 1)

	for i := 0; i < len(pageNums)-1; i++ {
		go worker(wg, results, tempDir, args.FileStem, pageNums[i]+1, pageNums[i+1], i+1)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	successFiles := make([]string, 0, len(pageNums)-1)
	failedFiles := make([]string, 0, len(pageNums)-1)
	for res := range results {
		if res.error != nil {
			failedFiles = append(failedFiles, res.string)
		} else {
			successFiles = append(successFiles, res.string)
		}
	}

	if len(failedFiles) > 0 {
		return successFiles, errors.New(fmt.Sprintf("Failed to generate the following files:\n%v", strings.Join(failedFiles, "\n")))
	}
	return successFiles, nil
}

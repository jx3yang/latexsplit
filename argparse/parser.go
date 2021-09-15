package argparse

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Arguments struct {
	FileName     string
	FileStem     string
	OutputDir    string
	Compiler     string
	CompilerArgs string
	SplitLine    string
}

func Parse() (*Arguments, error) {
	args := &Arguments{}
	flag.StringVar(&args.FileName, "file", "", "TeX file to compile")
	flag.StringVar(&args.OutputDir, "outdir", "", "Output directory")
	flag.StringVar(&args.Compiler, "c", "pdflatex", "Compiler")
	flag.StringVar(&args.CompilerArgs, "args", "", "Arguments to the compiler")
	flag.StringVar(&args.SplitLine, "split", "% latexsplit", "Line to split")

	flag.Parse()

	if args.FileName == "" {
		return nil, errors.New("usage: latexsplit -file <filename> [-outdir \"\"] [-c pdflatex] [-args \"-jobname <filestem>\"] [-split \"% latexsplit\"]")
	}

	_, err := os.Stat(args.FileName)
	if os.IsNotExist(err) {
		return nil, errors.New("Could not find file " + args.FileName)
	}
	fileBaseName := filepath.Base(args.FileName)
	args.FileStem = fileBaseName[0 : len(fileBaseName)-len(filepath.Ext(fileBaseName))]

	if len(args.CompilerArgs) > 0 {
		args.CompilerArgs = fmt.Sprintf("%v -jobname %v", args.CompilerArgs, args.FileStem)
	} else {
		args.CompilerArgs = fmt.Sprintf("-jobname %v", args.FileStem)
	}

	if len(args.OutputDir) > 0 {
		args.CompilerArgs = fmt.Sprintf("%v -output-directory %v", args.CompilerArgs, args.OutputDir)
	}
	return args, nil
}

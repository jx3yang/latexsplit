package main

import (
	"fmt"
	"os"
	"strings"

	argparse "github.com/jx3yang/latexsplit/argparse"
	"github.com/jx3yang/latexsplit/splitcompiler"
)

func main() {
	args, err := argparse.Parse()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	successFiles, err := splitcompiler.LatexSplitCompiler(args)
	if len(successFiles) > 0 {
		fmt.Printf("Successfully created:\n%v\n", strings.Join(successFiles, "\n"))
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

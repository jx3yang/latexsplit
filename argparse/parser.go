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
	Compiler     string
	CompilerArgs string
	SplitLine    string
}

func Parse() (*Arguments, error) {
	args := &Arguments{}
	flag.StringVar(&args.FileName, "file", "", "TeX file to compile")
	flag.StringVar(&args.Compiler, "c", "pdflatex", "Compiler")
	flag.StringVar(&args.CompilerArgs, "args", "", "Arguments to the compiler")
	flag.StringVar(&args.SplitLine, "split", "% latexsplit", "Line to split")

	flag.Parse()

	if args.FileName == "" {
		return nil, errors.New("usage: latexsplit -file <filename> [-c pdflatex] [-args \"-jobname=<filestem>\"] [-split \"% latexsplit\"]")
	}

	_, err := os.Stat(args.FileName)
	if os.IsNotExist(err) {
		return nil, errors.New("Could not find file " + args.FileName)
	}
	args.FileStem = filepath.Base(args.FileName)[0 : len(args.FileName)-len(filepath.Ext(args.FileName))]

	if len(args.CompilerArgs) > 0 {
		args.CompilerArgs = fmt.Sprintf("-jobname=%v %v", args.FileStem, args.CompilerArgs)
	} else {
		args.CompilerArgs = fmt.Sprintf("-jobname=%v", args.FileStem)
	}
	return args, nil
}

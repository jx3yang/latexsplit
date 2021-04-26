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
	if len(os.Args) < 2 {
		return nil, errors.New("usage: latexsplit <filename> [-c=pdflatex] [-args=\"-jobname=<filestem>\"] [-split=\"% latexsplit\"]")
	}

	args.FileName = os.Args[1]
	_, err := os.Stat(args.FileName)
	if os.IsNotExist(err) {
		return nil, errors.New("Could not find file " + args.FileName)
	}
	args.FileStem = filepath.Base(args.FileName)[0 : len(args.FileName)-len(filepath.Ext(args.FileName))]

	flag.StringVar(&args.Compiler, "Compiler name", "pdflatex", "")
	flag.StringVar(&args.CompilerArgs, "Compiler arguments", "", "")
	flag.StringVar(&args.SplitLine, "Line to split", "% latexsplit", "")

	flag.Parse()
	if len(args.CompilerArgs) > 0 {
		args.CompilerArgs = fmt.Sprintf("-jobname=%v %v", args.FileStem, args.CompilerArgs)
	} else {
		args.CompilerArgs = fmt.Sprintf("-jobname=%v", args.FileStem)
	}
	return args, nil
}

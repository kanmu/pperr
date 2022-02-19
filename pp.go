package pperr

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type Printer func(w io.Writer, e error, st, leaf errors.StackTrace)

var DefaultIndent = "\t"

var DefaultPrinter Printer = func(w io.Writer, err error, st, leaf errors.StackTrace) {
	fmt.Fprintf(w, "%T: %s\n", err, err.Error())

	if st != nil {
		var frames []*Frame

		if leaf == nil {
			frames = ExtractFrames(st)
		} else {
			frames = ExtractFramesN(st, 1)
		}

		for _, f := range frames {
			fmt.Fprintln(w, DefaultIndent+f.Name)
			fmt.Fprintf(w, "%s%s%s:%d\n", DefaultIndent, DefaultIndent, f.File, f.Line)
		}
	} else {
		fmt.Fprintln(w, DefaultIndent+"(no stack trace available)")
	}
}

func Fprint(w io.Writer, err error) {
	FprintFunc(w, err, DefaultPrinter)
}

func FprintFunc(w io.Writer, err error, puts Printer) {
	FprintFuncWithLeaf(w, err, puts, nil)
}

func FprintFuncWithLeaf(w io.Writer, err error, puts Printer, leaf errors.StackTrace) {
	if err == nil {
		return
	}

	if withStack, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
		puts(w, err, withStack.StackTrace(), leaf)

		if leaf == nil {
			leaf = withStack.StackTrace()
		}

		if withCause, ok := withStack.(interface{ Unwrap() error }); ok {
			err = withCause.Unwrap()
		}
	} else {
		puts(w, err, nil, leaf)
	}

	withCause, ok := err.(interface{ Unwrap() error })

	if !ok {
		return
	}

	FprintFuncWithLeaf(w, withCause.Unwrap(), puts, leaf)
}

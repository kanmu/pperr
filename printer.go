package pperr

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type Printer func(w io.Writer, err error, st, leaf errors.StackTrace)

var DefaultIndent = "\t"

var DefaultPrinterWithIndent = func(w io.Writer, err error, st, leaf errors.StackTrace, indent string) {
	fmt.Fprintf(w, "%T: %s\n", err, err.Error())

	if st != nil {
		var frames []*Frame

		if leaf == nil {
			frames = ExtractFrames(st)
		} else {
			frames = ExtractFramesN(st, 1)
		}

		for _, f := range frames {
			fmt.Fprintln(w, indent+f.Name)
			fmt.Fprintf(w, "%s%s%s:%d\n", indent, indent, f.File, f.Line)
		}
	} else {
		fmt.Fprintln(w, indent+"(no stack trace available)")
	}
}

var DefaultPrinter Printer = func(w io.Writer, err error, st, leaf errors.StackTrace) {
	DefaultPrinterWithIndent(w, err, st, leaf, DefaultIndent)
}

func NewPrinterWithIndent(indent string) Printer {
	return func(w io.Writer, err error, st, leaf errors.StackTrace) {
		DefaultPrinterWithIndent(w, err, st, leaf, indent)
	}
}

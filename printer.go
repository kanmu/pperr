package pperr

import (
	"fmt"
	"io"
)

type Printer func(w io.Writer, err error, frames, parent Frames)

var DefaultIndent = "\t"

var DefaultPrinterWithIndent = func(w io.Writer, err error, frames, parent Frames, indent string) {
	fmt.Fprintf(w, "%T: %s\n", err, err.Error())

	if frames != nil {
		if parent != nil {
			frames = frames.Exclude(parent)
		}

		for _, f := range frames {
			fmt.Fprintln(w, indent+f.Name)
			fmt.Fprintf(w, "%s%s%s:%d\n", indent, indent, f.File, f.Line)
		}
	}
}

var DefaultPrinter Printer = func(w io.Writer, err error, frames, parent Frames) {
	DefaultPrinterWithIndent(w, err, frames, parent, DefaultIndent)
}

func NewPrinterWithIndent(indent string) Printer {
	return func(w io.Writer, err error, frames, parent Frames) {
		DefaultPrinterWithIndent(w, err, frames, parent, indent)
	}
}

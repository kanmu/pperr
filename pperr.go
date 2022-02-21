package pperr

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func Print(err error) {
	Fprint(os.Stdout, err)
}

func PrintFunc(err error, puts Printer) {
	FprintFunc(os.Stdout, err, puts)
}

func Sprint(err error) string {
	var buf strings.Builder
	Fprint(&buf, err)
	return buf.String()
}

func SprintFunc(err error, puts Printer) string {
	var buf strings.Builder
	FprintFunc(&buf, err, puts)
	return buf.String()
}

func Fprint(w io.Writer, err error) {
	FprintFunc(w, err, DefaultPrinter)
}

func FprintFunc(w io.Writer, err error, puts Printer) {
	fprintFuncWithParent(w, err, puts, nil)
}

func fprintFuncWithParent(w io.Writer, err error, puts Printer, parent Frames) {
	if err == nil {
		return
	}

	realErr := err
	var frames Frames

	if withStack, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
		frames = ExtractFrames(withStack.StackTrace())

		if withCause, ok := withStack.(interface{ Unwrap() error }); ok {
			realErr = withCause.Unwrap()
		}
	}

	if withCause, ok := realErr.(interface{ Unwrap() error }); ok {
		var causeParent Frames

		if frames != nil {
			causeParent = frames
		} else {
			causeParent = parent
		}

		fprintFuncWithParent(w, withCause.Unwrap(), puts, causeParent)
	}

	puts(w, err, frames, parent)
}

func CauseType(err error) string {
	for {
		wrappedErr, ok := err.(interface{ Unwrap() error })

		if !ok {
			return fmt.Sprintf("%T", err)
		}

		cause := wrappedErr.Unwrap()

		if cause == nil {
			return fmt.Sprintf("%T", err)
		}

		err = cause
	}
}

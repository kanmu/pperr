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
	for _, e := range extractErrorSets(err, nil) {
		puts(w, e.Error, e.Frames, e.Parent)
	}
}

func ExtractErrorSets(err error) ErrorSets {
	return extractErrorSets(err, nil)
}

func extractErrorSets(err error, parent Frames) []ErrorSet {
	if err == nil {
		return nil
	}

	realErr := err
	var frames Frames

	if withStack, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
		frames = ExtractFrames(withStack.StackTrace())

		if withCause, ok := withStack.(interface{ Unwrap() error }); ok {
			realErr = withCause.Unwrap()
		}
	}

	var errs []ErrorSet

	if withCause, ok := realErr.(interface{ Unwrap() error }); ok {
		var causeParent Frames

		if frames != nil {
			causeParent = frames
		} else {
			causeParent = parent
		}

		if es := extractErrorSets(withCause.Unwrap(), causeParent); es != nil {
			errs = es
		}
	}

	errs = append(errs, ErrorSet{Error: err, Frames: frames, Parent: parent})

	return errs
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

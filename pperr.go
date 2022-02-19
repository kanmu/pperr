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

func CauseType(err error) string {
	if err == nil {
		return ""
	}

	if ws, ok := err.(interface {
		StackTrace() errors.StackTrace
		Cause() error
	}); ok {
		if cause, ok := ws.Cause().(interface{ Cause() error }); ok {
			return CauseType(cause.Cause())
		}
	}

	return fmt.Sprintf("%T", err)
}

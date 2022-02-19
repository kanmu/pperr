package pperr_test

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/pperr"
)

func f1() error {
	return errors.Wrap(f2(), "from f1()")
}

func f2() error {
	return errors.Wrap(f3(), "from f2()")
}

func f3() error {
	_, err := os.Open("not_found")
	return errors.Wrap(err, "from f3()")
}

func TestFprint(t *testing.T) {
	assert := assert.New(t)

	var buf strings.Builder
	err := f1()
	pperr.Fprint(&buf, err)

	actual := buf.String()
	actual = regexp.MustCompile(`(?m)[^\s>]+/go/.*:\d+$`).ReplaceAllString(actual, ".../go/...:NN")
	actual = regexp.MustCompile(`(?m)[^\s>]+/pperr_test.go:\d+$`).ReplaceAllString(actual, ".../pperr_test.go:NN")
	actual = regexp.MustCompile(`(?m):\d+$`).ReplaceAllString(actual, ":NN")

	expected := `*errors.withStack: from f1(): from f2(): from f3(): open not_found: no such file or directory
	github.com/winebarrel/pperr_test.f1
		.../pperr_test.go:NN
	github.com/winebarrel/pperr_test.TestFprint
		.../pperr_test.go:NN
	testing.tRunner
		.../go/...:NN
	runtime.goexit
		.../go/...:NN
*errors.withStack: from f2(): from f3(): open not_found: no such file or directory
	github.com/winebarrel/pperr_test.f2
		.../pperr_test.go:NN
*errors.withStack: from f3(): open not_found: no such file or directory
	github.com/winebarrel/pperr_test.f3
		.../pperr_test.go:NN
*fs.PathError: open not_found: no such file or directory
	(no stack trace available)
syscall.Errno: no such file or directory
	(no stack trace available)
`

	assert.Equal(expected, actual)
}

func TestFprint_StandardError(t *testing.T) {
	assert := assert.New(t)

	var buf strings.Builder
	pperr.Fprint(&buf, fmt.Errorf("standard error"))

	actual := buf.String()
	expected := `*errors.errorString: standard error
	(no stack trace available)
`

	assert.Equal(expected, actual)
}

func TestFprint_Nil(t *testing.T) {
	assert := assert.New(t)
	var buf strings.Builder
	pperr.Fprint(&buf, nil)
	assert.Equal("", buf.String())
}

func TestFprint_Indent(t *testing.T) {
	assert := assert.New(t)

	var buf strings.Builder
	err := f1()
	pperr.FprintFunc(&buf, err, pperr.NewPrinterWithIndent(">>"))

	actual := buf.String()
	actual = regexp.MustCompile(`(?m)[^\s>]+/go/.*:\d+$`).ReplaceAllString(actual, ".../go/...:NN")
	actual = regexp.MustCompile(`(?m)[^\s>]+/pperr_test.go:\d+$`).ReplaceAllString(actual, ".../pperr_test.go:NN")
	actual = regexp.MustCompile(`(?m):\d+$`).ReplaceAllString(actual, ":NN")

	expected := `*errors.withStack: from f1(): from f2(): from f3(): open not_found: no such file or directory
>>github.com/winebarrel/pperr_test.f1
>>>>.../pperr_test.go:NN
>>github.com/winebarrel/pperr_test.TestFprint_Indent
>>>>.../pperr_test.go:NN
>>testing.tRunner
>>>>.../go/...:NN
>>runtime.goexit
>>>>.../go/...:NN
*errors.withStack: from f2(): from f3(): open not_found: no such file or directory
>>github.com/winebarrel/pperr_test.f2
>>>>.../pperr_test.go:NN
*errors.withStack: from f3(): open not_found: no such file or directory
>>github.com/winebarrel/pperr_test.f3
>>>>.../pperr_test.go:NN
*fs.PathError: open not_found: no such file or directory
>>(no stack trace available)
syscall.Errno: no such file or directory
>>(no stack trace available)
`

	assert.Equal(expected, actual)
}

func TestSprint(t *testing.T) {
	assert := assert.New(t)

	err := f1()
	actual := pperr.Sprint(err)

	actual = regexp.MustCompile(`(?m)[^\s>]+/go/.*:\d+$`).ReplaceAllString(actual, ".../go/...:NN")
	actual = regexp.MustCompile(`(?m)[^\s>]+/pperr_test.go:\d+$`).ReplaceAllString(actual, ".../pperr_test.go:NN")
	actual = regexp.MustCompile(`(?m):\d+$`).ReplaceAllString(actual, ":NN")

	expected := `*errors.withStack: from f1(): from f2(): from f3(): open not_found: no such file or directory
	github.com/winebarrel/pperr_test.f1
		.../pperr_test.go:NN
	github.com/winebarrel/pperr_test.TestSprint
		.../pperr_test.go:NN
	testing.tRunner
		.../go/...:NN
	runtime.goexit
		.../go/...:NN
*errors.withStack: from f2(): from f3(): open not_found: no such file or directory
	github.com/winebarrel/pperr_test.f2
		.../pperr_test.go:NN
*errors.withStack: from f3(): open not_found: no such file or directory
	github.com/winebarrel/pperr_test.f3
		.../pperr_test.go:NN
*fs.PathError: open not_found: no such file or directory
	(no stack trace available)
syscall.Errno: no such file or directory
	(no stack trace available)
`

	assert.Equal(expected, actual)
}

func TestSprintFunc(t *testing.T) {
	assert := assert.New(t)

	err := f1()
	actual := pperr.SprintFunc(err, pperr.NewPrinterWithIndent(">>"))

	actual = regexp.MustCompile(`(?m)[^\s>]+/go/.*:\d+$`).ReplaceAllString(actual, ".../go/...:NN")
	actual = regexp.MustCompile(`(?m)[^\s>]+/pperr_test.go:\d+$`).ReplaceAllString(actual, ".../pperr_test.go:NN")
	actual = regexp.MustCompile(`(?m):\d+$`).ReplaceAllString(actual, ":NN")

	expected := `*errors.withStack: from f1(): from f2(): from f3(): open not_found: no such file or directory
>>github.com/winebarrel/pperr_test.f1
>>>>.../pperr_test.go:NN
>>github.com/winebarrel/pperr_test.TestSprintFunc
>>>>.../pperr_test.go:NN
>>testing.tRunner
>>>>.../go/...:NN
>>runtime.goexit
>>>>.../go/...:NN
*errors.withStack: from f2(): from f3(): open not_found: no such file or directory
>>github.com/winebarrel/pperr_test.f2
>>>>.../pperr_test.go:NN
*errors.withStack: from f3(): open not_found: no such file or directory
>>github.com/winebarrel/pperr_test.f3
>>>>.../pperr_test.go:NN
*fs.PathError: open not_found: no such file or directory
>>(no stack trace available)
syscall.Errno: no such file or directory
>>(no stack trace available)
`

	assert.Equal(expected, actual)
}

func TestCauseType(t *testing.T) {
	assert := assert.New(t)
	err := f1()
	assert.Equal("*fs.PathError", pperr.CauseType(err))
}

func TestCauseType_fundamental(t *testing.T) {
	assert := assert.New(t)
	err := errors.New("")
	assert.Equal("*errors.fundamental", pperr.CauseType(err))
}

func TestCauseType_errorString(t *testing.T) {
	assert := assert.New(t)
	err := fmt.Errorf("")
	assert.Equal("*errors.errorString", pperr.CauseType(err))
}

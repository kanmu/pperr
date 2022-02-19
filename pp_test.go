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
	actual = regexp.MustCompile(`/usr/local/Cellar/go/.*`).ReplaceAllString(actual, "/usr/local/Cellar/go/...")
	actual = regexp.MustCompile(`(?m):\d+$`).ReplaceAllString(actual, ":NN")

	expected := `*errors.withStack: from f1(): from f2(): from f3(): open not_found: no such file or directory
	github.com/winebarrel/pperr_test.f1
		/Users/sugawara/com/winebarrel/pperr/pp_test.go:NN
	github.com/winebarrel/pperr_test.TestFprint
		/Users/sugawara/com/winebarrel/pperr/pp_test.go:NN
	testing.tRunner
		/usr/local/Cellar/go/...
	runtime.goexit
		/usr/local/Cellar/go/...
*errors.withStack: from f2(): from f3(): open not_found: no such file or directory
	github.com/winebarrel/pperr_test.f2
		/Users/sugawara/com/winebarrel/pperr/pp_test.go:NN
*errors.withStack: from f3(): open not_found: no such file or directory
	github.com/winebarrel/pperr_test.f3
		/Users/sugawara/com/winebarrel/pperr/pp_test.go:NN
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
	actual = regexp.MustCompile(`/usr/local/Cellar/go/.*`).ReplaceAllString(actual, "/usr/local/Cellar/go/...")
	actual = regexp.MustCompile(`(?m):\d+$`).ReplaceAllString(actual, ":NN")

	expected := `*errors.withStack: from f1(): from f2(): from f3(): open not_found: no such file or directory
>>github.com/winebarrel/pperr_test.f1
>>>>/Users/sugawara/com/winebarrel/pperr/pp_test.go:NN
>>github.com/winebarrel/pperr_test.TestFprint_Indent
>>>>/Users/sugawara/com/winebarrel/pperr/pp_test.go:NN
>>testing.tRunner
>>>>/usr/local/Cellar/go/...
>>runtime.goexit
>>>>/usr/local/Cellar/go/...
*errors.withStack: from f2(): from f3(): open not_found: no such file or directory
>>github.com/winebarrel/pperr_test.f2
>>>>/Users/sugawara/com/winebarrel/pperr/pp_test.go:NN
*errors.withStack: from f3(): open not_found: no such file or directory
>>github.com/winebarrel/pperr_test.f3
>>>>/Users/sugawara/com/winebarrel/pperr/pp_test.go:NN
*fs.PathError: open not_found: no such file or directory
>>(no stack trace available)
syscall.Errno: no such file or directory
>>(no stack trace available)
`

	assert.Equal(expected, actual)
}

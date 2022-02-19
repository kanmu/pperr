# pperr

Pretty print for [pkg/errors](https://github.com/pkg/errors).

[![Build](https://github.com/winebarrel/pperr/actions/workflows/build.yml/badge.svg)](https://github.com/winebarrel/pperr/actions/workflows/build.yml)

# Usage

```go
package pp

import (
	"os"

	"github.com/pkg/errors"
	"github.com/winebarrel/pperr"
)

func main() {
	err := f1()
	pperr.Print(err)
}

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
```

```
*errors.withStack: from f1(): from f2(): from f3(): open not_found: no such file or directory
	main.f1
		/Users/.../main.go:16
	main.main
		/Users/.../main.go:11
	runtime.main
		/usr/local/Cellar/go/1.17.6/libexec/src/runtime/proc.go:255
	runtime.goexit
		/usr/local/Cellar/go/1.17.6/libexec/src/runtime/asm_amd64.s:1581
*errors.withStack: from f2(): from f3(): open not_found: no such file or directory
	main.f2
		/Users/.../main.go:20
*errors.withStack: from f3(): open not_found: no such file or directory
	main.f3
		/Users/.../main.go:25
*fs.PathError: open not_found: no such file or directory
	(no stack trace available)
syscall.Errno: no such file or directory
	(no stack trace available)
```

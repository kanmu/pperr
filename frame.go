package pp

import (
	"runtime"

	"github.com/pkg/errors"
)

type Frame struct {
	File string
	Line int
	Name string
}

func ExtractFrames(st errors.StackTrace) []*Frame {
	return ExtractFramesN(st, 0)
}

func ExtractFramesN(st errors.StackTrace, n int) []*Frame {
	var frames []*Frame

	if n > 0 {
		frames = make([]*Frame, 0, n)
	} else {
		frames = make([]*Frame, 0, len(st))
	}

	for i, v := range st {
		if n > 0 && i >= n {
			break
		}

		pc := uintptr(v) - 1
		fn := runtime.FuncForPC(pc)
		var frm *Frame

		if fn == nil {
			frm = &Frame{File: "unknown", Line: 0, Name: "unknown"}
		} else {
			file, line := fn.FileLine(pc)
			frm = &Frame{File: file, Line: line, Name: fn.Name()}
		}

		frames = append(frames, frm)
	}

	return frames
}

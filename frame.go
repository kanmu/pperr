package pperr

import (
	"runtime"

	"github.com/pkg/errors"
)

type ErrorSets []ErrorSet

func (es ErrorSets) TopFrame() *Frame {
	if len(es) == 0 {
		return nil
	} else if len(es[0].Frames) == 0 {
		return nil
	}

	top := es[0].Frames

	if p := es[0].Parent; p != nil {
		top = top.Exclude(p)
	}
	if len(top) == 0 {
		return nil
	}

	return top[0]
}

type ErrorSet struct {
	Error  error
	Frames Frames
	Parent Frames
}

type Frame struct {
	File string
	Line int
	Name string
}

type Frames []*Frame

func ExtractFrames(st errors.StackTrace) Frames {
	frames := make([]*Frame, 0, len(st))

	for _, v := range st {
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

func (frames Frames) Exclude(excludes Frames) Frames {
	newFrames := make(Frames, 0, len(frames))

L1:
	for _, f := range frames {
		for _, e := range excludes {
			if f.File == e.File && f.Line == e.Line && f.Name == e.Name {
				break L1
			}
		}

		newFrames = append(newFrames, f)
	}

	return newFrames
}

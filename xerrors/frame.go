package xerrors

import (
	"bytes"
	"runtime"
	"strconv"
)

// caller returns a Frame that describes a frame on the caller's stack.
// The argument skip is the number of frames to skip over.
// caller(0) returns the frame for the caller of caller.
func caller(skip uint8) [3]uintptr {
	var frames [3]uintptr
	runtime.Callers(int(skip+1), frames[:])
	return frames
}

// location reports the file, line, and function of a frame.
//
// The returned function may be "" even if file and line are not.
func location(framesPtrs [3]uintptr) (function, file string, line int) {
	frames := runtime.CallersFrames(framesPtrs[:])
	if _, ok := frames.Next(); !ok {
		return "", "", 0
	}
	fr, ok := frames.Next()
	if !ok {
		return "", "", 0
	}
	return fr.Function, fr.File, fr.Line
}

// formatFrames prints the stack as error detail.
// It should be called from an error's Format implementation,
// before printing any other error detail.
func formatFrames(function, file string, line int, buf *bytes.Buffer) {
	if function != "" {
		buf.WriteString(function)
		buf.WriteString(":")
	}
	if file != "" {
		buf.WriteString(file)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(line))
		buf.WriteString(" ")
	}
	if function != "" || file != "" {
		buf.Truncate(buf.Len() - 1)
	}
}

// FrameError is an error with part of a call stack.
type FrameError interface {
	Wrapper
	FrameLocation() (string, string, int)
}

// IsFrameError is a helper for type casting to FrameError
func IsFrameError(err error) bool {
	_, ok := err.(FrameError)
	return ok
}

// LastFrameError is a helper for Last with IsFrameError, returning a typed FrameError
func LastFrameError(err error) FrameError {
	err = Last(err, IsFrameError)
	if err == nil {
		return nil
	}
	return err.(FrameError)
}

type frameError struct {
	// Make room for three PCs: the one we were asked for, what it called,
	// and possibly a PC for skipPleaseUseCallersFrames. See:
	// https://go.googlesource.com/go/+/032678e0fb/src/runtime/extern.go#169
	frames [3]uintptr
	Wrapping
}

func (err *frameError) Error() string {
	function, file, line := err.FrameLocation()
	buf := bytes.Buffer{}
	formatFrames(function, file, line, &buf)
	return buf.String()
}

func (err *frameError) FrameLocation() (string, string, int) {
	return location(err.frames)
}

func newFrameError(skip uint8, err error) error {
	return &frameError{
		frames:   caller(skip + 1),
		Wrapping: Wrapping{err: err},
	}
}

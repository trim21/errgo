package errgo

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
)

// stack represents a stack of program counters.
type stack []uintptr

func (s stack) Format(st fmt.State, verb rune) {
	if verb == 'v' && st.Flag('+') {
		frames := runtime.CallersFrames(s)

		for {
			frame, more := frames.Next()
			_, _ = io.WriteString(st, "\n")
			_, _ = io.WriteString(st, frame.Function)
			_, _ = io.WriteString(st, "\n\t")
			_, _ = io.WriteString(st, frame.File)
			_, _ = io.WriteString(st, ":")
			_, _ = io.WriteString(st, strconv.Itoa(frame.Line))
			if !more {
				break
			}
		}
	}
}

func callers() stack {
	const depth = 16
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[:n]
}

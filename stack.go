package errgo

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

// stack represents a stack of program counters.
type stack []uintptr

func (s stack) String() string {
	st := &strings.Builder{}

	s.writeTo(st)

	return st.String()
}

func (s stack) writeTo(w io.Writer) {
	frames := runtime.CallersFrames(s)

	for {
		frame, more := frames.Next()
		_, _ = fmt.Fprintf(w, "  at %s (%s:%d)\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}

func (s stack) Format(st fmt.State, verb rune) {
	if verb == 'v' && st.Flag('+') {
		s.writeTo(st)
	}
}

func callers() stack {
	const depth = 16
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[:n]
}

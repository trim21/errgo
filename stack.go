package errgo

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
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
		_, _ = io.WriteString(w, frame.Function)
		_, _ = io.WriteString(w, "\n\t")
		_, _ = io.WriteString(w, frame.File)
		_, _ = io.WriteString(w, ":")
		_, _ = io.WriteString(w, strconv.Itoa(frame.Line))
		if !more {
			break
		}
		_, _ = io.WriteString(w, "\n")
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

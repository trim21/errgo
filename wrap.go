package errgo

import (
	"fmt"
	"io"
	"runtime"
	"strconv"

	"github.com/valyala/bytebufferpool"
)

type unwrap interface {
	Unwrap() error
}

type Stack interface {
	Stack() string
}

var _ unwrap = (*wrapError)(nil)
var _ unwrap = (*msgError)(nil)
var _ unwrap = (*withStackError)(nil)
var _ Stack = (*withStackError)(nil)

type wrapError struct {
	err error
	msg string
}

func (e *wrapError) Error() string {
	return e.msg + ": " + e.err.Error()
}

func (e *wrapError) Unwrap() error {
	return e.err
}

type msgError struct {
	err error
	msg string
}

func (e *msgError) Error() string {
	return e.msg
}

func (e *msgError) Unwrap() error {
	return e.err
}

type withStackError struct {
	Err   error
	stack stack
}

func (w *withStackError) Error() string {
	return w.Err.Error()
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withStackError) Unwrap() error { return w.Err }

// Format implement fmt.Formatter, add trace to zap.Error(err).
func (w *withStackError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'T':
		_, _ = io.WriteString(s, "*withStackError")
	case 'v':
		// _, _ = io.WriteString(s, w.Error())
		if s.Flag('#') {
			fmt.Fprintf(s, "&errgo.withStackError{Err: %#v, stack: ...}", w.Err)
			return
		}

		_, _ = io.WriteString(s, "error: ")
		_, _ = io.WriteString(s, w.Err.Error())
		if s.Flag('+') {
			s.Write([]byte("\nstack:\n"))
			w.stack.Format(s, verb)
			return
		}
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = io.WriteString(s, strconv.Quote(w.Error()))
	}
}

func (w *withStackError) Stack() string {
	return fmt.Sprintf("%+v", w.stack)
}

// MarshalJSON marshal error with stack to
//
//	{
//		"error": "context: real error",
//		"stack": [
//			"main.main ...main.go:54",
//			"..."
//		]
//	}
func (w *withStackError) MarshalJSON() ([]byte, error) {
	if w == nil {
		return []byte("null"), nil
	}

	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	b.WriteString(`{"error":`)
	b.B = strconv.AppendQuote(b.B, w.Error())
	b.WriteString(",")

	b.WriteString(`"stack":[`)

	frames := runtime.CallersFrames(w.stack)
	for {
		frame, more := frames.Next()
		b.WriteString(`"`)
		b.WriteString(frame.Function)
		b.WriteString(" ")
		b.WriteString(frame.File)
		b.WriteString(":")
		b.B = strconv.AppendInt(b.B, int64(frame.Line), 10)
		b.WriteString(`"`)
		if !more {
			break
		}
		b.WriteString(",")
	}

	b.WriteString("]}")

	return append([]byte{}, b.Bytes()...), nil
}

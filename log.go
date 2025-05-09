package errgo

import (
	"errors"
	"runtime"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/valyala/bytebufferpool"
)

func ZerologErrorMarshaler(err error) any {
	var j *withStackError
	if errors.As(err, &j) {
		return j
	}

	return err
}

func (w *withStackError) MarshalZerologObject(e *zerolog.Event) {
	e.Str("message", w.Error())
	a := zerolog.Arr()

	frames := runtime.CallersFrames(w.stack)
	for {
		frame, more := frames.Next()
		a.Object(zerologStack{
			Func: frame.Function,
			File: frame.File,
			Lino: frame.Line,
		})
		if !more {
			break
		}
	}

	e.Array("stack", a)
}

var _ zerolog.LogObjectMarshaler = stack{}
var _ zerolog.LogObjectMarshaler = (*withStackError)(nil)
var _ zerolog.LogObjectMarshaler = zerologStack{}

type zerologStack struct {
	Func string
	File string
	Lino int
}

func (z zerologStack) MarshalZerologObject(e *zerolog.Event) {
	e.Str("function", z.Func)
	e.Str("file", z.File)
	e.Int("lino", z.Lino)
}

func (s stack) MarshalZerologObject(e *zerolog.Event) {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	b.WriteString("[")

	frames := runtime.CallersFrames(s)
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

	b.WriteString("]")

	e.RawJSON("stack", b.Bytes())
}

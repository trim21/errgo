//nolint:goerr113
package errgo_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/trim21/errgo"
)

func TestWrap(t *testing.T) {
	t.Parallel()
	err := errors.New("raw")
	require.Equal(t, "wrap: raw", errgo.Wrap(err, "wrap").Error())
	require.Equal(t, "e: wrap: raw", errgo.Wrap(errgo.Wrap(err, "wrap"), "e").Error())
}

func TestStackTrace(t *testing.T) {
	t.Parallel()

	err := errgo.Wrap(errors.New("a error"), "m")
	s := fmt.Sprintf("%+v", err)
	require.Regexp(t,
		regexp.MustCompile("^m: a error\nerror stack:\n.*errgo_test\\.TestStackTrace\n.*wrap_test.go:\\d+\n.*"),
		s)
}

func TestErrorIs(t *testing.T) {
	t.Parallel()

	e := errors.New("expected")

	err := errgo.Wrap(e, "ctx")
	require.True(t, errors.Is(err, e))

	err = errgo.MsgNoTrace(e, "ctx")
	require.True(t, errors.Is(err, e))

	err = errgo.Msg(e, "ctx")
	require.True(t, errors.Is(err, e))
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	e := errors.New("expected")
	err := errgo.Wrap(e, "ctx")

	b, jerr := json.Marshal(err)
	require.NoError(t, jerr)

	var m struct {
		Error string   `json:"error"`
		Stack []string `json:"stack"`
	}

	require.NoError(t, json.Unmarshal(b, &m))

	require.Equal(t, "ctx: expected", m.Error)
	require.NotZero(t, len(m.Stack), "stack should not be zero")
}

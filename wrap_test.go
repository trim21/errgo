//nolint:goerr113
package errgo_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/require"

	"github.com/trim21/errgo"
)

// not depth enough error stack may contain code from stdlib testing and runtime
// and their lino will change with go version
func errWithTraceDepthEnough() error {
	return errDepth(0, errors.New("original error"))
}

func errDepth(depth int, err error) error {
	if depth == 17 {
		return errgo.Wrap(err, fmt.Sprintf("ctx %d", depth))
	}

	return errgo.Wrap(errDepth(depth+1, err), fmt.Sprintf("ctx %d", depth))
}

func TestWrap(t *testing.T) {
	t.Parallel()
	err := errors.New("raw")
	require.Equal(t, "wrap: raw", errgo.Wrap(err, "wrap").Error())
	require.Equal(t, "e: wrap: raw", errgo.Wrap(errgo.Wrap(err, "wrap"), "e").Error())
}

func TestFormat(t *testing.T) {
	t.Parallel()

	err := errWithTraceDepthEnough()
	t.Run("+v", func(t *testing.T) {
		s := fmt.Sprintf("%+v", err)
		cupaloy.SnapshotT(t, s)
	})

	t.Run("v", func(t *testing.T) {
		s := fmt.Sprintf("%v", err)
		cupaloy.SnapshotT(t, s)
	})
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

	err := errWithTraceDepthEnough()

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")

	require.NoError(t, enc.Encode(err))
	cupaloy.SnapshotT(t, buf.String())

	var m struct {
		Error string   `json:"error"`
		Stack []string `json:"stack"`
	}

	require.NoError(t, json.Unmarshal(buf.Bytes(), &m))

	require.True(t, strings.HasSuffix(err.Error(), "ctx 17: original error"))
	require.NotZero(t, len(m.Stack), "stack should not be zero")
}

func TestStack(t *testing.T) {
	t.Parallel()

	err := errWithTraceDepthEnough()

	cupaloy.SnapshotT(t, err.(errgo.Stack).Stack())
}

package errgo

// Wrap add context to error message.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withStackError); ok { //nolint:errorlint
		// keep Stack
		return &withStackError{
			Err:   &wrapError{msg: msg, err: e.Err},
			Stack: e.Stack,
		}
	}

	return &withStackError{Err: &wrapError{msg: msg, err: err}, Stack: callers()}
}

// Msg replace error message.
func Msg(err error, msg string) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withStackError); ok { //nolint:errorlint
		// keep traces
		return &withStackError{
			Err:   &msgError{msg: msg, err: e.Err},
			Stack: e.Stack,
		}
	}

	return &withStackError{Err: &msgError{msg: msg, err: err}, Stack: callers()}
}

// MsgNoTrace replace error message without adding trace.
// this is used to create global errors, which also avoid add trace.
func MsgNoTrace(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &msgError{msg: msg, err: err}
}

// Trace add trace to error, without change error message.
func Trace(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withStackError); ok { //nolint:errorlint
		// keep Stack
		return &withStackError{
			Err:   e.Err,
			Stack: e.Stack,
		}
	}

	return &withStackError{Err: err, Stack: callers()}
}

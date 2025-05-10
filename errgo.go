package errgo

// Wrap add context to error message.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withStackError); ok {
		// keep stack
		return &withStackError{
			Err:   &wrapError{msg: msg + ": " + err.Error(), err: e.Err},
			stack: e.stack,
		}
	}

	return &withStackError{Err: &wrapError{msg: msg + ": " + err.Error(), err: err}, stack: callers()}
}

// Msg replace error message.
func Msg(err error, msg string) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withStackError); ok {
		// keep traces
		return &withStackError{
			Err:   &msgError{msg: msg, err: e.Err},
			stack: e.stack,
		}
	}

	return &withStackError{Err: &msgError{msg: msg, err: err}, stack: callers()}
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

	if e, ok := err.(*withStackError); ok {
		// keep stack
		return &withStackError{
			Err:   e.Err,
			stack: e.stack,
		}
	}

	return &withStackError{Err: err, stack: callers()}
}

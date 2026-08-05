package server

import "fmt"

type InternalError struct {
	err  error
	msg  string
	args []any
}

func NewInternalError(err error, msg string, args ...any) *InternalError {
	return &InternalError{
		err:  err,
		msg:  msg,
		args: args,
	}
}

func (e *InternalError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e *InternalError) Unwrap() error { return e.err }
func (e *InternalError) Args() []any   { return e.args }
func (e *InternalError) Msg() string   { return e.msg }

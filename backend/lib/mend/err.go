package mend

import (
	"fmt"
	"runtime/debug"
)

type AppErr struct {
	msg         string
	debugValues map[string]string
	stackTrace  []string
	wrapped     error
}

func (e *AppErr) Error() string {
	if e.wrapped == nil || e.wrapped.Error() == e.msg {
		return e.msg
	}
	return e.msg + " caused by: " + e.wrapped.Error()
}

func (e *AppErr) Unwrap() error {
	return e.wrapped
}

func (e *AppErr) Str(key string, val any) *AppErr {
	e.debugValues[key] = fmt.Sprintf("%v", val)
	return e
}

func (e *AppErr) StackTrace() []string {
	return e.stackTrace
}

func Err(msg string, withStackTrace bool) *AppErr {
	err := &AppErr{
		msg:         msg,
		debugValues: map[string]string{},
	}
	if withStackTrace {
		err.stackTrace = []string{string(debug.Stack())}
	}
	return err
}

func Wrap(err error, withStackTrace bool) *AppErr {
	if appErr, ok := err.(*AppErr); ok {
		return appErr
	}

	appErr := Err(err.Error(), withStackTrace)
	appErr.wrapped = err
	return appErr
}

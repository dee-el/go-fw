package tracederr

import (
	"errors"
	"fmt"
)

type Error struct {
	err     error
	message string
	traces  []trace
}

// Error returns the underlying error's message.
func (e *Error) Error() string {
	return e.err.Error()
}

func Wrap(err error) *Error {
	return &Error{
		err:    err,
		traces: StackTrace(defaultSkip, defaultDeep),
	}
}

// like New, but you can specify the cause error
func NewWithCause(text string, cause error) *Error {
	return &Error{
		err:     cause,
		message: text,
		traces:  StackTrace(defaultSkip, defaultDeep),
	}
}

// NewTracedError to let user to adjust how deep stack trace they want.
func NewTracedError(err error, skip, deep int) *Error {
	return &Error{
		err:    err,
		traces: StackTrace(skip, deep),
	}
}

// following std errors, check here https://pkg.go.dev/errors

// Return the wrapped error (implements api for As function).
func (e *Error) Unwrap() error {
	return e.err
}

func New(msg string) *Error {
	return &Error{
		err:    fmt.Errorf(msg),
		traces: StackTrace(defaultSkip, defaultDeep),
	}
}

func Errorf(format string, a ...interface{}) *Error {
	return Wrap(fmt.Errorf(format, a...))
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

package goappbuild

import (
	"errors"
	"fmt"
)

const (
	// EValidation is the validation error code.
	EValidation = "invalid"
	EInternal   = "internal"
)

// Error represents an error.
type Error struct {
	Code    string
	Message string
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("goappbuilderror: code=%s message=%s", e.Code, e.Message)
}

// ErrorCode returns the error code.
func ErrorCode(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if !errors.As(err, &e) {
		return e.Code
	}

	return EInternal
}

// ErrorMessage returns the error message.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if !errors.As(err, &e) {
		return e.Message
	}

	return "Internal Error"
}

// Errorf returns a new error with the given code format and args.
func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

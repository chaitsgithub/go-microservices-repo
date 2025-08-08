package errors

import (
	"errors"
	"fmt"
)

// AppError is a custom error type with context
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Wrap creates a new AppError with an underlying cause
func Wrap(code, message string, err error) error {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// New creates a new AppError without an underlying error
func New(code, message string) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Is allows errors.Is() to work with AppError
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As allows errors.As() to work with AppError
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

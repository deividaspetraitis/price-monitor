package errors

import (
	"errors"
	"fmt"
)

// New constructs a new error from text string.
func New(text string) error {
	return errors.New(text)
}

// Newf constructs a new error from formatted text string.
func Newf(format string, a ...interface{}) error {
	return New(fmt.Sprintf(format, a...))
}

// Wrap wraps text in err and returns resulting error.
func Wrap(err error, text string) error {
	return fmt.Errorf("%s: %w", text, err)
}

// Wrapf wraps formatted text in err and returns resulting error.
func Wrapf(err error, format string, a ...interface{}) error {
	return Wrap(err, fmt.Sprintf(format, a...))
}

// Equals compares two errors based on their contents.
// It is safe to pass nil errors.
func Equals(err1, err2 error) bool {
	if err1 != nil && err2 == nil {
		return false
	}

	if err1 == nil && err2 != nil {
		return false
	}

	if err1 == nil && err2 == nil {
		return true
	}

	return err1.Error() == err2.Error()
}

// Is checks whether err is target error.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

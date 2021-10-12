package errors

import (
	"errors"
	"fmt"
)

func Wrap(err error, format string, v ...interface{}) error {
	msg := fmt.Sprintf(format, v...)
	return fmt.Errorf("%s %w", msg, err)
}

func Is(err, target error) bool { return errors.Is(err, target) }

func New(format string, v ...interface{}) error { return fmt.Errorf(format, v...) }

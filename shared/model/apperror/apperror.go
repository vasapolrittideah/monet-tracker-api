package apperror

import (
	"fmt"

	"github.com/vasapolrittideah/money-tracker-api/shared/constants/errorcode"
)

type Error struct {
	Code errorcode.Code
	Err  error
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %v", e.Code, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code errorcode.Code, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

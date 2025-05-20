package apperror

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

type Error struct {
	Code codes.Code
	Err  error
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %v", e.Code.String(), e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code codes.Code, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

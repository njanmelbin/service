package errs

import (
	"errors"
	"fmt"
)

// Error represents an error in the system.
type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

// New constructs an error based on an app error.
func New(code ErrCode, err error) Error {
	return Error{
		Code:    code,
		Message: err.Error(),
	}
}

// Newf constructs an error based on a error message.
func Newf(code ErrCode, format string, v ...any) Error {
	return Error{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
	}
}

func (err Error) Error() string {
	return err.Message
}

func IsError(err error) bool {
	var er Error
	return errors.As(err, &er)
}

func GetError(err error) Error {
	var er Error
	if !errors.As(err, &er) {
		return Error{}
	}
	return er
}

// HTTPStatus implements the web package httpStatus interface so the
// web framework can use the correct http status.
func (e Error) HTTPStatus() int {
	return httpStatus[e.Code]
}

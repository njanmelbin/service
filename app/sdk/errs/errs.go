package errs

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Error represents an error in the system.
type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

// New constructs an error based on an app error.
func New(code ErrCode, err error) *Error {
	return &Error{
		Code:    code,
		Message: err.Error(),
	}
}

// Newf constructs an error based on a error message.
func Newf(code ErrCode, format string, v ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
	}
}

func (err *Error) Error() string {
	return err.Message
}

func IsError(err error) bool {
	var er *Error
	return errors.As(err, &er)
}

func GetError(err error) *Error {
	var er *Error
	if !errors.As(err, &er) {
		return &Error{}
	}
	return er
}

// HTTPStatus implements the web package httpStatus interface so the
// web framework can use the correct http status.
func (e Error) HTTPStatus() int {
	return httpStatus[e.Code]
}

// =============================================================================

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

type FieldErrors []FieldError

func NewFieldErrors(field string, err error) *Error {
	fe := FieldErrors{
		{
			Field: field,
			Err:   err.Error(),
		},
	}

	return fe.ToError()
}

// Add adds a field error to the collection.
func (fe *FieldErrors) Add(field string, err error) {
	*fe = append(*fe, FieldError{
		Field: field,
		Err:   err.Error(),
	})
}

// ToError converts the field errors to an Error.
func (fe FieldErrors) ToError() *Error {
	return New(InvalidArgument, fe)
}

// Error implements the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

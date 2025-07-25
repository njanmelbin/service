package web

import (
	"fmt"
	"net/http"

	"github.com/go-json-experiment/json"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	return r.PathValue(key)
}

// Decoder represents data that can be decoded.
type Decoder interface {
	Decode(data []byte) error
}

type validator interface {
	Validate() error
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
// If the provided value is a struct then it is checked for validation tags.
// If the value implements a validate function, it is executed.
func Decode(r *http.Request, v Decoder) error {
	if err := json.UnmarshalRead(r.Body, v, json.RejectUnknownMembers(false)); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if v, ok := v.(validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

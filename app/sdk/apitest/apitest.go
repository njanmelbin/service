package apitest

import (
	"net/http"
	"service/app/sdk/auth"
	"service/business/sdk/dbtest"
)

type testOption struct {
	skip    bool
	skipMsg string
}

type OptionFunc func(*testOption)

// WithSkip can be used to skip running a test.
func WithSkip(skip bool, msg string) OptionFunc {
	return func(to *testOption) {
		to.skip = skip
		to.skipMsg = msg
	}
}

// Test contains functions for executing an api test.
type Test struct {
	DB   *dbtest.Database
	Auth *auth.Auth
	mux  http.Handler
}

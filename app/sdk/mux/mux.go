package mux

import (
	"net/http"
	"service/app/domain/checkapi"
)

func WebAPI() *http.ServeMux {
	mux := http.NewServeMux()
	checkapi.Routes(mux)
	return mux
}

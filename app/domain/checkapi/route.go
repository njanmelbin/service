package checkapi

import "net/http"

func Routes(mux *http.ServeMux) {

	mux.HandleFunc("GET /v1/liveness", liveness)
	mux.HandleFunc("GET /v1/readiness", readiness)
}

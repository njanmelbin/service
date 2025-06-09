package checkapi

import (
	"net/http"
	"service/foundation/web"
)

func Routes(app *web.App) {

	const version = "v1"

	app.HandleFunc(http.MethodGet, version, "/liveness", liveness)
	app.HandleFunc(http.MethodGet, version, "/readiness", readiness)
	app.HandleFunc(http.MethodGet, version, "/panic", panics)
	app.HandleFunc(http.MethodGet, version, "/errors", errorsHandler)
}

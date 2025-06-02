package checkapi

import "service/foundation/web"

func Routes(app *web.App) {

	app.HandleFunc("GET /v1/liveness", liveness)
	app.HandleFunc("GET /v1/readiness", readiness)
	app.HandleFunc("GET /v1/panic", panics)
}

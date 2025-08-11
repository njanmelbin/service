package checkapp

import (
	"net/http"
	"service/foundation/logger"
	"service/foundation/web"

	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

func Routes(app *web.App, cfg Config) {

	const version = "v1"

	api := newApp(cfg.Build, cfg.Log, cfg.DB)

	app.HandlerFuncNoMid(http.MethodGet, version, "/liveness", api.liveness)
	app.HandlerFuncNoMid(http.MethodGet, version, "/readiness", api.readiness)
	// app.HandleFunc(http.MethodGet, version, "/panic", panics)
	// app.HandleFunc(http.MethodGet, version, "/errors", errorsHandler)
}

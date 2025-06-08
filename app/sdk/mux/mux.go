package mux

import (
	"os"
	"service/app/domain/checkapi"
	"service/app/sdk/mid"
	"service/foundation/logger"
	"service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build    string
	Log      *logger.Logger
	Shutdown chan os.Signal
}

func WebAPI(cfg Config) *web.App {
	mux := web.New(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	checkapi.Routes(mux)

	return mux
}

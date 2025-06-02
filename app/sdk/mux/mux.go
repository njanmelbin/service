package mux

import (
	"os"
	"service/app/domain/checkapi"
	"service/app/sdk/mid"
	"service/foundation/logger"
	"service/foundation/web"
)

func WebAPI(shutdown chan os.Signal, log *logger.Logger) *web.App {
	mux := web.New(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics())

	checkapi.Routes(mux)

	return mux
}

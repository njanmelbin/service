package mux

import (
	"os"
	"service/app/domain/checkapi"
	"service/app/domain/mid"
	"service/foundation/logger"
	"service/foundation/web"
)

func WebAPI(shutdown chan os.Signal, log *logger.Logger) *web.App {
	mux := web.New(shutdown, mid.Logger(log))

	checkapi.Routes(mux)

	return mux
}

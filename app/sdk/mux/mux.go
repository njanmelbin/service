package mux

import (
	"os"
	"service/app/domain/checkapi"
	"service/foundation/web"
)

func WebAPI(shutdown chan os.Signal) *web.App {
	mux := web.New(shutdown)
	checkapi.Routes(mux)
	return mux
}

package all

import (
	"service/app/domain/checkapp"
	"service/app/domain/userapp"
	"service/app/sdk/mux"
	"service/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

func (add) Add(app *web.App, cfg mux.Config) {

	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})
}

package all

import (
	"service/app/domain/authapp"
	"service/app/domain/checkapi"
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
	checkapi.Routes(app)

	authapp.Routes(app, authapp.Config{
		Auth: cfg.AuthConfig.Auth,
	})
}

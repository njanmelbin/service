package all

import (
	"service/app/domain/checkapi"
	"service/app/domain/userapp"
	"service/app/sdk/mux"
	"service/business/domain/userbus"
	"service/business/domain/userbus/stores/userdb"
	"service/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

func (add) Add(app *web.App, cfg mux.Config) {
	userBus := userbus.NewBusiness(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))
	checkapi.Routes(app)
	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    userBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})
}

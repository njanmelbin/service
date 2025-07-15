package all

import (
	"service/app/domain/checkapp"
	"service/app/domain/userapp"
	"service/app/sdk/mux"
	"service/business/domain/userbus"
	"service/business/domain/userbus/stores/userdb"
	"service/business/sdk/delegate"
	"service/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

func (add) Add(app *web.App, cfg mux.Config) {
	delegate := delegate.New(cfg.Log)

	userBus := userbus.NewBusiness(cfg.Log, delegate, userdb.NewStore(cfg.Log, cfg.DB))

	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    userBus,
		AuthClient: cfg.SalesConfig.AuthClient,
	})
}

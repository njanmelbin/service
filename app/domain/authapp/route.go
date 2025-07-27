package authapp

import (
	"net/http"
	"service/app/sdk/auth"
	"service/app/sdk/mid"
	"service/business/domain/userbus"
	"service/foundation/web"
)

type Config struct {
	UserBus userbus.ExtBusiness
	Auth    *auth.Auth
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.Auth)
	basic := mid.Basic(cfg.Auth, cfg.UserBus)
	bearer := mid.Bearer(cfg.Auth)

	app.HandleFunc(http.MethodGet, version, "/auth/token/{kid}", api.token, basic)
	app.HandleFunc(http.MethodGet, version, "/auth/authenticate", api.authenticate, bearer)
	app.HandleFunc(http.MethodPost, version, "/auth/authorize", api.authorize)
}

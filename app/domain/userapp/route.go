package userapp

import (
	"net/http"
	"service/app/sdk/auth"
	"service/app/sdk/authclient"
	"service/app/sdk/mid"
	"service/business/domain/userbus"
	"service/foundation/logger"
	"service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	UserBus    userbus.ExtBusiness
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"
	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := NewApp(cfg.UserBus)
	app.HandleFunc(http.MethodPost, version, "/users", api.create, authen, ruleAdmin)

}

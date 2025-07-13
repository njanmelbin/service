package mux

import (
	"os"
	"service/app/sdk/auth"
	"service/app/sdk/authclient"
	"service/app/sdk/mid"
	"service/foundation/logger"
	"service/foundation/web"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

// AuthConfig contains auth service specific config.
type AuthConfig struct {
	Auth *auth.Auth
}

// SalesConfig contains sales service specific config.
type SalesConfig struct {
	AuthClient *authclient.Client
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build       string
	Log         *logger.Logger
	Shutdown    chan os.Signal
	AuthConfig  AuthConfig
	SalesConfig SalesConfig
	DB          *sqlx.DB
	Tracer      trace.Tracer
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

func WebAPI(cfg Config, routeAdder RouteAdder) *web.App {
	mux := web.New(cfg.Shutdown, cfg.Tracer,
		mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	routeAdder.Add(mux, cfg)

	return mux
}

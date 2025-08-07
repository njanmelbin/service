package mux

import (
	"embed"
	"os"
	"service/app/sdk/auth"
	"service/app/sdk/authclient"
	"service/app/sdk/mid"
	"service/business/domain/userbus"
	"service/foundation/logger"
	"service/foundation/web"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

// StaticSite represents a static site to run.
type StaticSite struct {
	react      bool
	static     embed.FS
	staticDir  string
	staticPath string
}

// Options represent optional parameters.
type Options struct {
	corsOrigin []string
	sites      []StaticSite
}

// AuthConfig contains auth service specific config.
type AuthConfig struct {
	Auth *auth.Auth
}

// SalesConfig contains sales service specific config.
type SalesConfig struct {
	AuthClient *authclient.Client
}

type BusConfig struct {
	UserBus userbus.ExtBusiness
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build       string
	Log         *logger.Logger
	Shutdown    chan os.Signal
	BusConfig   BusConfig
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

func WebAPI(cfg Config, routeAdder RouteAdder, options ...func(opts *Options)) *web.App {
	mux := web.New(cfg.Shutdown, cfg.Tracer, mid.Otel(cfg.Tracer),
		mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	routeAdder.Add(mux, cfg)

	var opts Options
	for _, option := range options {
		option(&opts)
	}

	return mux
}

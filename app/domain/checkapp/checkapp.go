package checkapp

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"service/app/sdk/errs"
	"service/business/sdk/sqldb"
	"service/foundation/logger"
	"service/foundation/web"
	"time"

	"github.com/jmoiron/sqlx"
)

type app struct {
	build string
	log   *logger.Logger
	db    *sqlx.DB
}

func newApp(build string, log *logger.Logger, db *sqlx.DB) *app {
	return &app{
		build: build,
		log:   log,
		db:    db,
	}
}

func (a *app) liveness(ctx context.Context, r *http.Request) web.Encoder {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := Info{
		Status:     "up",
		Build:      a.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
	}

	// This handler provides a free timer loop.

	return info
}

func (a *app) readiness(ctx context.Context, r *http.Request) web.Encoder {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := sqldb.StatusCheck(ctx, a.db); err != nil {
		a.log.Info(ctx, "readiness failure", "ERROR", err)
		return errs.New(errs.Internal, err)
	}

	return nil
}

// func panics(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 	if n := rand.Intn(100); n%2 == 0 {
// 		panic("we are panicking")
// 	}

// 	status := struct {
// 		Status string
// 	}{
// 		Status: "OK",
// 	}

// 	return web.Respond(ctx, w, status, http.StatusOK)

// }

// func errorsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 	if n := rand.Intn(100); n%2 == 0 {
// 		return errs.Newf(errs.Canceled, "error cancelled")
// 	}

// 	status := struct {
// 		Status string
// 	}{
// 		Status: "OK",
// 	}

// 	return web.Respond(ctx, w, status, http.StatusOK)

// }

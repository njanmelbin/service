package mid

import (
	"context"
	"net/http"
	"service/foundation/logger"
	"service/foundation/web"
)

func Logger(log *logger.Logger) web.MidFunc {
	m := func(handler web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			log.Info(ctx, "request started", "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Info(ctx, "request completed", "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr)

			return err
		}
		return h
	}
	return m
}

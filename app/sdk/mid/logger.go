package mid

import (
	"context"
	"net/http"
	"service/foundation/logger"
	"service/foundation/web"
	"time"
)

func Logger(log *logger.Logger) web.MidFunc {
	m := func(handler web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			v := web.GetValues(ctx)

			log.Info(ctx, "request started", "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr)

			err := handler(ctx, r)

			log.Info(ctx, "request completed", "method", r.Method, "path", r.URL.Path, "remoteAddr", r.RemoteAddr, "statuscode", v.StatusCode, "since", time.Since(v.Now).String())

			return err
		}
		return h
	}
	return m
}

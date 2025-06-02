package mid

import (
	"context"
	"net/http"
	"service/app/sdk/metrics"
	"service/foundation/web"
)

func Metrics() web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)

			err := next(ctx, w, r)

			n := metrics.AddRequests(ctx)

			if n%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}

			if err != nil {
				metrics.AddErrors(ctx)
			}

			return err
		}
		return h
	}
	return m
}

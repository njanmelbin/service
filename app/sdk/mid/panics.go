package mid

import (
	"context"
	"net/http"
	"runtime/debug"
	"service/app/sdk/errs"
	"service/app/sdk/metrics"
	"service/foundation/web"
)

func Panics() web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) (resp web.Encoder) {

			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					resp = errs.Newf(errs.InternalOnlyLog, "PANIC [%v] TRACE[%s]", rec, string(trace))
					metrics.AddPanics(ctx)
				}
			}()

			return next(ctx, r)
		}
		return h

	}
	return m
}

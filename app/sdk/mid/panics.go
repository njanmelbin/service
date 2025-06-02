package mid

import (
	"context"
	"net/http"
	"runtime/debug"
	"service/app/sdk/errs"
	"service/foundation/web"
)

func Panics() web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = errs.Newf(errs.InternalOnlyLog, "PANIC [%v] TRACE[%s]", rec, string(trace))

				}
			}()

			return next(ctx, w, r)
		}
		return h

	}
	return m
}

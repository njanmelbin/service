package mid

import (
	"context"
	"net/http"
	"service/app/sdk/errs"
	"service/foundation/logger"
	"service/foundation/web"
)

func Errors(log *logger.Logger) web.MidFunc {

	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := next(ctx, w, r); err != nil {
				log.Error(ctx, "message", "ERROR", err)

				var er *errs.Error

				switch {
				case errs.IsError(err):
					er = errs.GetError(err)
				default:
					er = errs.Newf(errs.Unknown, errs.Unknown.String())
				}

				if err := web.Respond(ctx, w, er, er.HTTPStatus()); err != nil {
					return err
				}
			}
			return nil
		}
		return h
	}

	return m
}

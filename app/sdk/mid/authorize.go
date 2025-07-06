package mid

import (
	"context"
	"net/http"
	"service/app/sdk/authclient"
	"service/app/sdk/errs"
	"service/foundation/web"
	"time"
)

func Authorize(client *authclient.Client, rule string) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			userID, err := GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}
			auth := authclient.Authorize{
				Claims: GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
			}
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := client.Authorize(ctx, auth); err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			return next(ctx, r)
		}
		return h
	}
	return m
}

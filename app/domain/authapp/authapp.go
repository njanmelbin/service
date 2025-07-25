package authapp

import (
	"context"
	"net/http"
	"service/app/sdk/auth"
	"service/app/sdk/authclient"
	"service/app/sdk/errs"
	"service/app/sdk/mid"
	"service/foundation/web"
)

type app struct {
	auth *auth.Auth
}

func newApp(ath *auth.Auth) *app {
	return &app{
		auth: ath,
	}
}

func (a *app) token(ctx context.Context, r *http.Request) web.Encoder {
	kid := web.Param(r, "kid")

	if kid == "" {
		return errs.NewFieldErrors("kid", errs.Newf(errs.FailedPrecondition, "missing kid"))
	}

	// The BearerBasic middleware function generates the claims.
	claims := mid.GetClaims(ctx)

	tkn, err := a.auth.GenerateToken(kid, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}
	return token{Token: tkn}

}

func (a *app) authenticate(ctx context.Context, r *http.Request) web.Encoder {

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	resp := authclient.AuthenticateResp{
		UserID: userID,
		Claims: mid.GetClaims(ctx),
	}
	return resp
}

func (a *app) authorize(ctx context.Context, r *http.Request) web.Encoder {

	var auth authclient.Authorize
	if err := web.Decode(r, &auth); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := a.auth.Authorize(ctx, auth.Claims, auth.UserID, auth.Rule); err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]", auth.Claims.Roles, auth.Rule)
	}

	return nil
}

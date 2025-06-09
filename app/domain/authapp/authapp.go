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

func (a *app) token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
	token := token{Token: tkn}
	return web.Respond(ctx, w, token, http.StatusOK)

}

func (a *app) authenticate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	resp := authclient.AuthenticateResp{
		UserID: userID,
		Claims: mid.GetClaims(ctx),
	}
	return web.Respond(ctx, w, resp, http.StatusOK)
}

func (a *app) authorize(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var auth authclient.Authorize
	if err := web.Decode(r, &auth); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := a.auth.Authorize(ctx, auth.Claims, auth.UserID, auth.Rule); err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]", auth.Claims.Roles, auth.Rule)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

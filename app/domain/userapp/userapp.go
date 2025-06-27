// Package userapp maintains the app layer api for the user domain.
package userapp

import (
	"context"
	"errors"
	"net/http"
	"service/app/sdk/auth"
	"service/app/sdk/errs"
	"service/app/sdk/mid"
	"service/business/domain/userbus"
	"service/foundation/web"
)

// App manages the set of app layer api functions for the user domain.
type App struct {
	userBus userbus.ExtBusiness
	auth    *auth.Auth
}

// NewApp constructs a user app API for use.
func NewApp(userBus userbus.ExtBusiness) *App {
	return &App{
		userBus: userBus,
	}
}

// NewAppWithAuth constructs a user app API for use with auth support.
func NewAppWithAuth(userBus userbus.ExtBusiness, ath *auth.Auth) *App {
	return &App{
		auth:    ath,
		userBus: userBus,
	}
}

// Create adds a new user to the system.
func (a *App) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewUser
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nc, err := toBusNewUser(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, mid.GetSubjectID(ctx), nc)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return errs.Newf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}

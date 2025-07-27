package mid

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/mail"
	"service/app/sdk/auth"
	"service/app/sdk/authclient"
	"service/app/sdk/errs"
	"service/business/domain/userbus"
	"service/business/types/role"
	"service/foundation/web"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func Authenticate(client *authclient.Client) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {

			resp, err := client.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			ctx = setUserID(ctx, resp.UserID)
			ctx = setClaims(ctx, resp.Claims)

			return next(ctx, r)

		}
		return h
	}
	return m
}

// Bearer processes JWT authentication logic.
func Bearer(ath *auth.Auth) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			claims, err := ath.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			if claims.Subject == "" {
				return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claim")
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
			}

			ctx = setUserID(ctx, subjectID)
			ctx = setClaims(ctx, claims)

			return next(ctx, r)
		}
		return h
	}
	return m
}

// Basic processes basic authentication logic.
func Basic(ath *auth.Auth, userBus userbus.ExtBusiness) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			email, pass, ok := parseBasicAuth(r.Header.Get("authorization"))
			if !ok {
				return errs.Newf(errs.Unauthenticated, "invalid basic auth")
			}

			addr, err := mail.ParseAddress(email)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			usr, err := userBus.Authenticate(ctx, *addr, pass)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   usr.ID.String(),
					Issuer:    ath.Issuer(),
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
				Roles: role.ParseToString(usr.Roles),
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
			}

			ctx = setUserID(ctx, subjectID)
			ctx = setClaims(ctx, claims)

			return next(ctx, r)
		}
		return h
	}
	return m
}

func parseBasicAuth(auth string) (string, string, bool) {
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}

	c, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", false
	}

	return username, password, true
}

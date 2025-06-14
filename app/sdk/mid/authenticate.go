package mid

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/mail"
	"service/app/sdk/auth"
	"service/app/sdk/authclient"
	"service/app/sdk/errs"
	"service/foundation/web"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func Authenticate(client *authclient.Client) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			resp, err := client.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			ctx = setUserID(ctx, resp.UserID)
			ctx = setClaims(ctx, resp.Claims)

			return next(ctx, w, r)

		}
		return h
	}
	return m
}

// Bearer processes JWT authentication logic.
func Bearer(ath *auth.Auth) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

			return next(ctx, w, r)
		}
		return h
	}
	return m
}

// Basic processes basic authentication logic.
func Basic(ath *auth.Auth) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			email, _, ok := parseBasicAuth(r.Header.Get("authorization"))
			if !ok {
				return errs.Newf(errs.Unauthenticated, "invalid basic auth")
			}

			_, err := mail.ParseAddress(email)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "5cf37266-3473-4006-984f-9325122678b7",
					Issuer:    ath.Issuer(),
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
				Roles: []string{"USER"},
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
			}

			ctx = setUserID(ctx, subjectID)
			ctx = setClaims(ctx, claims)

			return next(ctx, w, r)
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

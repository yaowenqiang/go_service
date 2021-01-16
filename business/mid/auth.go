package mid

import (

	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/yaowenqiang/service/business/auth"
	"github.com/yaowenqiang/service/foundation/web"
	"github.com/pkg/errors"
)

var ErrForbidden = web.NewRequestError(
	errors.New("you are not authorized for that action"),
	http.StatusForbidden,
)


func Authenticate(a *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			authStr := r.Header.Get("authorization")
			parts := strings.Split(authStr, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return web.NewRequestError(err ,http.StatusUnauthorized)
			}

			ctx = context.WithValue(ctx, auth.Key, claims)

			return handler(ctx, w, r)
		}
			return h
	}
	return m
}


func Authorize(roles ...string) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				claims, ok := ctx.Value(auth.Key).(auth.Claims)
				if !ok {
					return errors.New("claims missing from context")
				}

				if !claims.Authorized(roles...) {
					return web.NewRequestError(
							fmt.Errorf("you are not authorized for that action : claims: %v, exp: %v", claims.Roles, roles),
							http.StatusForbidden)
				}
				return handler(ctx, w, r)
			}
		return h
	}
	return m
}

package middleware

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/util"
)

func Authentication(app *util.Application) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		BeforeFunc: func(ctx echo.Context) {
			ctx.Response().Header().Add(
				echo.HeaderVary,
				echo.HeaderAuthorization,
			)
		},
		SigningKey: []byte(app.Config.JWT.Access),
		NewClaimsFunc: func(ctx echo.Context) jwt.Claims {
			return &util.CustomAccessJWTClaims{}
		},
		ContextKey: util.ContextKeyUser,
	})
}

func PublicAuthentication(app *util.Application) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		BeforeFunc: func(ctx echo.Context) {
			ctx.Response().Header().Add(
				echo.HeaderVary,
				echo.HeaderAuthorization,
			)
		},
		SigningKey: []byte(app.Config.JWT.Access),
		NewClaimsFunc: func(ctx echo.Context) jwt.Claims {
			return &util.CustomAccessJWTClaims{}
		},
		ContextKey:             util.ContextKeyUser,
		ContinueOnIgnoredError: true,
		ErrorHandler: func(ctx echo.Context, err error) error {
			token := ctx.Request().Header.Get(echo.HeaderAuthorization)
			if strings.TrimSpace(token) == "" {
				ctx.Set(util.ContextKeyUser, new(jwt.Token))
				return nil
			}

			return &echo.HTTPError{
				Code:    echojwt.ErrJWTInvalid.Code,
				Message: echojwt.ErrJWTInvalid.Message,
			}
		},
	})
}

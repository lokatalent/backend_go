package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/util"
)

// RequireVerification ensures that user attempting to access a resource
// has been verified.
func RequireVerification(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := util.ContextGetUser(ctx)
		if !user.IsVerified {
			return &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "User is unverified!",
			}
		}
		return nil
	}
}

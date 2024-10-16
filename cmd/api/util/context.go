package util

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/internal/models"
)

type ContextKey string

// ContextGetUser returns the current user data stored in the request.
func ContextGetUser(ctx echo.Context) models.User {
	claims, ok := ctx.Get(ContextKeyUser).(*jwt.Token).Claims.(*CustomAccessJWTClaims)
	if !ok {
		return models.User{}
	}
	return models.User{
		ID:         claims.ID,
		Email:      claims.Email,
		IsVerified: claims.Verified,
	}
}

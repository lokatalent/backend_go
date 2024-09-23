package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/server/handlers"
	"github.com/lokatalent/backend_go/internal/server/middleware"
)

func setAuthRoutes(app *util.Application, engine *echo.Echo) {
	handler := handlers.NewAuthHandler(app)

	auth := engine.Group("auth")
	auth.POST("/refresh-token", handler.RefreshToken)
	auth.PATCH("/verify-user", handler.VerifyUser, middleware.Authentication(app))
}

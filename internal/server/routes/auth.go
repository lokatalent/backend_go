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
	auth.GET("/:provider", handler.ProviderAuthentication)
	auth.GET("/:provider/callback", handler.ProviderAuthCallback)
	auth.POST("/signup", handler.SignUp)
	auth.GET("/signin", handler.SignIn)
	auth.PATCH("/verify-user", handler.VerifyUser, middleware.Authentication(app))
	auth.GET(
		"/verify-user/contact", handler.VerifyContact,
		middleware.Authentication(app))
	auth.PATCH(
		"/verify-user/contact/callback", handler.VerifyContactCallback,
		middleware.Authentication(app))
}

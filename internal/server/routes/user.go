package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/server/handlers"
	"github.com/lokatalent/backend_go/internal/server/middleware"
)

func setUserRoutes(app *util.Application, engine *echo.Echo) {
	handler := handlers.NewUserHandler(app)

	user := engine.Group("users")
	user.GET("", handler.ListUsers, middleware.PublicAuthentication(app))
	user.GET("/search", handler.Search, middleware.PublicAuthentication(app))
	user.PATCH(
		"/:id/set-role", handler.ChangeRole, middleware.Authentication(app),
		middleware.RequireVerification)
	user.PATCH(
		"/set-service-role", handler.ChangeServiceRole,
		middleware.Authentication(app), middleware.RequireVerification)

	// profile
	user.GET("/profile", handler.GetOwnProfile, middleware.Authentication(app))
	user.GET(
		"/profile/:id", handler.GetProfile, middleware.PublicAuthentication(app))
	user.PATCH(
		"/profile/update", handler.UpdateProfile, middleware.Authentication(app))
	user.PATCH(
		"/profile/picture-update", handler.UpdateProfileImage,
		middleware.Authentication(app))
	user.DELETE(
		"/profile/picture-delete", handler.DeleteProfileImage,
		middleware.Authentication(app), middleware.RequireVerification)
}

package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/server/handlers"
	"github.com/lokatalent/backend_go/internal/server/middleware"
)

func setNotificationRoutes(app *util.Application, engine *echo.Echo) {
	handler := handlers.NewNotificationHandler(app)

	notification := engine.Group("notification")
	notification.GET(
		"",
		handler.Get,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	notification.PATCH(
		"/read",
		handler.MarkAllAsRead,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	notification.PATCH(
		"/:id/read",
		handler.MarkAsRead,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	notification.GET(
		"/count",
		handler.Count,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
}

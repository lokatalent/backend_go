package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/server/handlers"
	"github.com/lokatalent/backend_go/internal/server/middleware"
)

func setBookingRoutes(app *util.Application, engine *echo.Echo) {
	handler := handlers.NewBookingHandler(app)

	booking := engine.Group("booking")
	booking.POST(
		"",
		handler.CreateBooking,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	booking.GET(
		"/:id",
		handler.GetBooking,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	booking.GET(
		"/all",
		handler.GetAllBookings,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	booking.PATCH(
		"/:id/status",
		handler.UpdateBookingStatus,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	booking.PATCH(
		"/:id/accept",
		handler.AcceptBooking,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	booking.PATCH(
		"/:id/reject",
		handler.RejectBooking,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	booking.GET(
		"/:id/find-providers",
		handler.FindProviders,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	booking.PUT(
		"/select-provider",
		handler.SelectProvider,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
}

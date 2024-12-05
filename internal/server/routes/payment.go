package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/server/handlers"
	"github.com/lokatalent/backend_go/internal/server/middleware"
)

func setPaymentRoutes(app *util.Application, engine *echo.Echo) {
	handler := handlers.NewPaymentHandler(app)

	payment := engine.Group("payment")
	payment.POST(
		"/initialize-transaction",
		handler.InitializeTransaction,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	payment.GET(
		"/verify-transaction",
		handler.VerifyTransaction,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	payment.GET(
		"/track-transactions",
		handler.TrackPayments,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
}

package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/server/handlers"
	"github.com/lokatalent/backend_go/internal/server/middleware"
)

func setServicePricingRoutes(app *util.Application, engine *echo.Echo) {
	handler := handlers.NewServicePricingHandler(app)

	servicePricing := engine.Group("service-pricing")
	servicePricing.POST(
		"",
		handler.CreateServicePricing,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	servicePricing.GET(
		"",
		handler.GetServicePricing,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	servicePricing.GET(
		"/all",
		handler.GetAllServicesPricing,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	servicePricing.PATCH(
		"",
		handler.UpdateServicePricing,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	servicePricing.DELETE(
		"",
		handler.DeleteServicePricing,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)

	// commission
	servicePricing.POST(
		"/commission",
		handler.CreateServiceCommission,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	servicePricing.GET("/commission", handler.GetServiceCommission)
	servicePricing.PATCH(
		"/commission/:id",
		handler.UpdateServiceCommission,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
}

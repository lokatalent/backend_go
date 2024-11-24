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
		"/profile/education",
		handler.GetOwnEducationProfile, middleware.Authentication(app))
	user.GET(
		"/profile/bank",
		handler.GetOwnBankProfile, middleware.Authentication(app))
	user.GET(
		"/:id/profile", handler.GetProfile, middleware.PublicAuthentication(app))
	user.GET(
		"/:id/profile/education",
		handler.GetEducationProfile, middleware.PublicAuthentication(app))
	user.PATCH(
		"/profile",
		handler.UpdateProfilePersonal,
		middleware.Authentication(app),
	)
	user.PATCH(
		"/profile/education",
		handler.UpdateProfileEducation,
		middleware.Authentication(app),
	)
	user.PATCH(
		"/profile/bank",
		handler.UpdateProfileBank,
		middleware.Authentication(app),
	)
	user.PATCH(
		"/profile/picture-update", handler.UpdateProfileImage,
		middleware.Authentication(app))
	user.DELETE(
		"/profile/picture-delete", handler.DeleteProfileImage,
		middleware.Authentication(app), middleware.RequireVerification)

	user.GET(
		"/:id/profile/certifications",
		handler.GetCertifications,
		middleware.Authentication(app),
	)
	user.POST(
		"/profile/certifications",
		handler.UploadCertifications,
		middleware.Authentication(app),
	)
	user.DELETE(
		"/profile/certifications/:id",
		handler.DeleteCertification,
		middleware.Authentication(app),
	)

	// services
	user.GET("/:id/service", handler.GetService, middleware.Authentication(app))
	user.GET("/:id/service/list", handler.ListServices, middleware.Authentication(app))
	user.POST(
		"/service", handler.CreateService, middleware.Authentication(app))
	user.PATCH(
		"/service", handler.UpdateService, middleware.Authentication(app))
	user.DELETE(
		"/service", handler.DeleteService, middleware.Authentication(app))

	user.GET(
		"/:id/service/images",
		handler.GetServiceImages,
		middleware.Authentication(app),
	)
	user.POST(
		"/service/images",
		handler.UploadServiceImages,
		middleware.Authentication(app),
	)
	user.DELETE(
		"/service/images/:id",
		handler.DeleteServiceImage,
		middleware.Authentication(app),
	)

	// wallet
	user.GET(
		"/wallet",
		handler.GetWallet,
		middleware.Authentication(app),
		middleware.RequireVerification,
	)
	user.GET(
		"/:id/wallet/debits",
		handler.GetDebits,
		middleware.Authentication(app),
	)

	user.POST(
		"/waitlist",
		handler.JoinWaitlist,
	)
}

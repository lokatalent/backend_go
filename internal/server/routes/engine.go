package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/lokatalent/backend_go/cmd/api/util"
)

// Engine configures server engine
func Engine(app *util.Application) *echo.Echo {
	engine := echo.New()
	engine.HideBanner = true
	engine.Validator = util.NewCustomValidator()

	var CORSConfig middleware.CORSConfig

	switch app.Config.Env {
	case util.ENVIRONMENT_PRODUCTION:
		CORSConfig = middleware.CORSConfig{
			Skipper: middleware.DefaultCORSConfig.Skipper,
			AllowOrigins: []string{
				"https://lokatalent.io",
				"https://staging.lokatalent.io",
			},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPut,
				http.MethodPatch,
				http.MethodPost,
				http.MethodDelete,
				http.MethodOptions,
			},
			AllowHeaders: []string{"*"},
		}
	default:
		CORSConfig = middleware.CORSConfig{
			Skipper:      middleware.DefaultCORSConfig.Skipper,
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPut,
				http.MethodPatch,
				http.MethodPost,
				http.MethodDelete,
				http.MethodOptions,
			},
			AllowHeaders: []string{"*"},
		}
	}

	// set middlewares
	engine.Use(middleware.CORSWithConfig(CORSConfig))
	engine.Pre(middleware.RemoveTrailingSlash())
	engine.Use(middleware.Recover())
	engine.Use(middleware.Logger())

	// set routes
	setAuthRoutes(app, engine)
	setUserRoutes(app, engine)

	return engine
}

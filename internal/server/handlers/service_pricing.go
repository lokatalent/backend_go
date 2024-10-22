package handlers

import (
	// "errors"
	"net/http"

	// "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	// "golang.org/x/sync/errgroup"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/models"
	// "github.com/lokatalent/backend_go/internal/repository"
)

type ServicePricingHandler struct {
	app *util.Application
}

func NewServicePricingHandler(app *util.Application) ServicePricingHandler {
	return ServicePricingHandler{app: app}
}

func (s ServicePricingHandler) CreateServicePricing(ctx echo.Context) error {
	reqData := struct {
		ServiceType string  `json:"service_type" validate:"required"`
		RatePerHour float64 `json:"rate_per_hour" validate:"required,gte=1"`
	}{}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := s.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil || (authUser.Role != models.USER_ADMIN_SUPER) {
		if err == nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"only super admin can create service price")
		}
		return util.ErrInternalServer(ctx, err)
	}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if !util.IsValidServiceType(reqData.ServiceType) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrInvalidServiceType,
		)
	}

	servicePricing := models.ServicePricing{
		ServiceType: reqData.ServiceType,
		RatePerHour: reqData.RatePerHour,
	}
	err = s.app.Repositories.ServicePricing.CreateServicePricing(&servicePricing)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, servicePricing)
}

func (s ServicePricingHandler) GetServicePricing(ctx echo.Context) error {
	serviceType := ctx.QueryParam("service_type")
	if !util.IsValidServiceType(serviceType) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrInvalidServiceType,
		)
	}

	servicePricing, err := s.app.Repositories.ServicePricing.GetServicePricing(serviceType)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, servicePricing)
}

func (s ServicePricingHandler) GetAllServicesPricing(ctx echo.Context) error {
	servicesPricing, err := s.app.Repositories.ServicePricing.GetAllServicesPricing()
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, servicesPricing)
}

func (s ServicePricingHandler) UpdateServicePricing(ctx echo.Context) error {
	reqData := struct {
		ServiceType string  `json:"service_type" validate:"required"`
		RatePerHour float64 `json:"rate_per_hour" validate:"required,gte=1"`
	}{}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := s.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil || (authUser.Role != models.USER_ADMIN_SUPER) {
		if err == nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"only super admin can update service price")
		}
		return util.ErrInternalServer(ctx, err)
	}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if !util.IsValidServiceType(reqData.ServiceType) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrInvalidServiceType,
		)
	}

	servicePricing := models.ServicePricing{
		ServiceType: reqData.ServiceType,
		RatePerHour: reqData.RatePerHour,
	}
	err = s.app.Repositories.ServicePricing.UpdateServicePricing(&servicePricing)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, servicePricing)
}

func (s ServicePricingHandler) DeleteServicePricing(ctx echo.Context) error {
	serviceType := ctx.QueryParam("service_type")
	if !util.IsValidServiceType(serviceType) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrInvalidServiceType,
		)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := s.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil || (authUser.Role != models.USER_ADMIN_SUPER) {
		if err == nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"only super admin can delete service price")
		}
		return util.ErrInternalServer(ctx, err)
	}

	err = s.app.Repositories.ServicePricing.DeleteServicePricing(serviceType)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, "successfully deleted.")
}

// commission

func (s ServicePricingHandler) CreateServiceCommission(ctx echo.Context) error {
	reqData := struct {
		Percentage int `json:"percentage" validate:"required,gte=1"`
	}{}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := s.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil || (authUser.Role != models.USER_ADMIN_SUPER) {
		if err == nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"only super admin can create commission")
		}
		return util.ErrInternalServer(ctx, err)
	}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.ErrBadRequest
	}

	serviceCommission := models.ServiceCommission{
		Percentage: reqData.Percentage,
	}
	err = s.app.Repositories.Commission.CreateServiceCommission(&serviceCommission)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, serviceCommission)
}

func (s ServicePricingHandler) GetServiceCommission(ctx echo.Context) error {
	serviceCommission, err := s.app.Repositories.Commission.GetServiceCommission()
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, serviceCommission)
}

func (s ServicePricingHandler) UpdateServiceCommission(ctx echo.Context) error {
	id := ctx.Param("id")
	reqData := struct {
		Percentage int `json:"percentage" validate:"required,gte=1"`
	}{}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := s.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil || (authUser.Role != models.USER_ADMIN_SUPER) {
		if err == nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"only super admin can update commission")
		}
		return util.ErrInternalServer(ctx, err)
	}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.ErrBadRequest
	}

	serviceCommission := models.ServiceCommission{
		ID:         id,
		Percentage: reqData.Percentage,
	}
	err = s.app.Repositories.Commission.UpdateServiceCommission(&serviceCommission)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, serviceCommission)
}

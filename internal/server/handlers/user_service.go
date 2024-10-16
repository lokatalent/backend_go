package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

func (u UserHandler) GetService(ctx echo.Context) error {
	id := ctx.Param("id")
	serviceType := ctx.QueryParam("service_type")
	// authenticatedUser := util.ContextGetUser(ctx)

	if !util.IsValidServiceType(serviceType) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrInvalidServiceType,
		)
	}

	service, err := u.app.Repositories.User.GetService(id, serviceType)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, service)
}

func (u UserHandler) ListServices(ctx echo.Context) error {
	id := ctx.Param("id")
	// authenticatedUser := util.ContextGetUser(ctx)

	services, err := u.app.Repositories.User.GetAllServices(id)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, services)
}

func (u UserHandler) CreateService(ctx echo.Context) error {
	reqData := struct {
		ServiceType     string              `json:"service_type" validate:"required"`
		ServiceDesc     string              `json:"service_desc" validate:"required"`
		RatePerHour     float64             `json:"rate_per_hour" validate:"required"`
		ExperienceYears int                 `json:"experience_years" validate:"required"`
		Address         string              `json:"address" validate:"required"`
		Availability    models.Availability `json:"availability" validate:"required"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if !util.IsValidServiceType(reqData.ServiceType) {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidServiceType)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := u.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}
	if authUser.ServiceRole == models.SERVICE_REQUESTER {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"service requester can not create service!",
		)
	}

	newService := models.UserService{
		UserID:          authUser.ID,
		ServiceType:     reqData.ServiceType,
		ServiceDesc:     reqData.ServiceDesc,
		RatePerHour:     reqData.RatePerHour,
		ExperienceYears: reqData.ExperienceYears,
		Address:         reqData.Address,
		Availability:    reqData.Availability,
	}

	err = u.app.Repositories.User.CreateService(&newService)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateService) {
			return echo.NewHTTPError(
				http.StatusForbidden,
				repository.ErrDuplicateService,
			)
		}
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, newService)
}

func (u UserHandler) UpdateService(ctx echo.Context) error {
	reqData := struct {
		ServiceType     string              `json:"service_type" validate:"required"`
		ServiceDesc     string              `json:"service_desc" validate:"required"`
		RatePerHour     float64             `json:"rate_per_hour" validate:"required"`
		ExperienceYears int                 `json:"experience_years" validate:"required"`
		Address         string              `json:"address" validate:"required"`
		Availability    models.Availability `json:"availability" validate:"required"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if !util.IsValidServiceType(reqData.ServiceType) {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidServiceType)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := u.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	newService := models.UserService{
		UserID:          authUser.ID,
		ServiceType:     reqData.ServiceType,
		ServiceDesc:     reqData.ServiceDesc,
		RatePerHour:     reqData.RatePerHour,
		ExperienceYears: reqData.ExperienceYears,
		Address:         reqData.Address,
		Availability:    reqData.Availability,
	}

	err = u.app.Repositories.User.UpdateService(&newService)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, newService)
}

func (u UserHandler) DeleteService(ctx echo.Context) error {
	serviceType := ctx.QueryParam("service_type")

	if !util.IsValidServiceType(serviceType) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrInvalidServiceType,
		)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := u.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	err = u.app.Repositories.User.DeleteService(authUser.ID, serviceType)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, "service successfully deleted")
}

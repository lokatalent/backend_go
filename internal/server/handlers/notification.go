package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/models/response"
	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/models"
	// "github.com/lokatalent/backend_go/internal/repository"
)

type NotificationHandler struct {
	app *util.Application
}

func NewNotificationHandler(app *util.Application) NotificationHandler {
	return NotificationHandler{app: app}
}

func (n NotificationHandler) Get(ctx echo.Context) error {
	reqPage, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid page value")
	}
	reqSize, err := strconv.Atoi(ctx.QueryParam("size"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid size value")
	}
	reqData := struct {
		Type      string `json:"type"`
		BookingID string `json:"booking_id"`
		// Page      int    `json:"page" validate:"required,gte=1"`
		// Size      int    `json:"size" validate:"required,gte=1"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	/*
		if err := ctx.Validate(&reqData); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	*/

	authenticatedUser := util.ContextGetUser(ctx)

	filter := models.NotificationFilter{
		Type:      reqData.Type,
		UserID:    authenticatedUser.ID,
		BookingID: reqData.BookingID,
		Page:      reqPage,
		Limit:     reqSize,
	}

	notifications, err := n.app.Repositories.Notification.GetForUser(filter)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	resp := []response.NotificationResponse{}
	for _, notification := range notifications {
		resp = append(resp, response.NotificationResponseFromModel(notification))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (n NotificationHandler) MarkAsRead(ctx echo.Context) error {
	notificationID := ctx.Param("id")
	authenticatedUser := util.ContextGetUser(ctx)
	err := n.app.Repositories.Notification.MarkAsRead(authenticatedUser.ID, notificationID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, map[string]string{})
}

func (n NotificationHandler) MarkAllAsRead(ctx echo.Context) error {
	authenticatedUser := util.ContextGetUser(ctx)
	err := n.app.Repositories.Notification.MarkAllAsRead(authenticatedUser.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, map[string]string{})
}

func (n NotificationHandler) Count(ctx echo.Context) error {
	authenticatedUser := util.ContextGetUser(ctx)

	count, err := n.app.Repositories.Notification.CountUnseen(authenticatedUser.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, count)
}

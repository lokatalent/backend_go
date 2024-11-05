package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/models/response"
	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

var (
	BOOKING_PROVIDER_SELECTION_MESSAGE = "You have been selected for booking %s."
)

type BookingHandler struct {
	app *util.Application
}

func NewBookingHandler(app *util.Application) BookingHandler {
	return BookingHandler{app: app}
}

func (b BookingHandler) CreateBooking(ctx echo.Context) error {
	reqData := struct {
		RequesterID   string  `json:"requester_id" validate:"required"`
		RequesterAddr string  `json:"requester_addr" validate:"required"`
		ServiceType   string  `json:"service_type" validate:"required"`
		BookingType   string  `json:"booking_type" validate:"required"`
		ServiceDesc   string  `json:"service_desc" validate:"required"`
		StartTime     string  `json:"start_time" validate:"required"`
		EndTime       string  `json:"end_time" validate:"required"`
		StartDate     string  `json:"start_date" validate:"required"`
		EndDate       string  `json:"end_date" validate:"required"`
		TotalPrice    float64 `json:"total_price" validate:"required"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if !util.IsValidServiceType(reqData.ServiceType) {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidServiceType)
	}
	if !util.IsValidBookingType(reqData.BookingType) {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidBookingType)
	}
	if !util.ValidPlaceAddress(reqData.RequesterAddr) {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidPlaceAddress)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := b.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}
	if authUser.ServiceRole == models.SERVICE_PROVIDER {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"service provider can not place booking",
		)
	}
	if reqData.RequesterID != authUser.ID {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"requester should be authenticated user.",
		)
	}

	var parsedStartDate, parsedEndDate, parsedStartTime, parsedEndTime time.Time
	// var err error

	parsedStartDate, err = util.ParseDate(reqData.StartDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	parsedEndDate, err = util.ParseDate(reqData.EndDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	parsedStartTime, err = util.ParseTime(
		fmt.Sprintf(
			"%sT%s",
			reqData.StartDate,
			reqData.StartTime,
		),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	parsedEndTime, err = util.ParseTime(
		fmt.Sprintf(
			"%sT%s",
			reqData.EndDate,
			reqData.EndTime,
		),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	calculatedTotalPrice, err := calculateBookingPrice(
		b.app,
		reqData.ServiceType,
		reqData.BookingType,
		parsedStartDate,
		parsedEndDate,
		parsedStartTime,
		parsedEndTime,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newBooking := models.Booking{
		ID:            uuid.NewString(),
		RequesterID:   authUser.ID,
		RequesterAddr: reqData.RequesterAddr,
		ServiceType:   reqData.ServiceType,
		BookingType:   reqData.BookingType,
		ServiceDesc:   reqData.ServiceDesc,
		StartTime:     parsedStartTime,
		EndTime:       parsedEndTime,
		StartDate:     parsedStartDate,
		EndDate:       parsedEndDate,
		TotalPrice:    calculatedTotalPrice,
		Status:        models.BOOKING_OPEN,
	}

	err = b.app.Repositories.Booking.Create(&newBooking)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response.BookingResponseFromModel(newBooking))
}

func (b BookingHandler) GetBooking(ctx echo.Context) error {
	bookingID := ctx.Param("id")

	booking, err := b.app.Repositories.Booking.GetByID(bookingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := b.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	if (booking.RequesterID != authUser.ID &&
		booking.ProviderID.String != authUser.ID) &&
		authUser.Role == models.USER_REGULAR {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"restricted from view booking.",
		)
	}

	return ctx.JSON(http.StatusOK, response.BookingResponseFromModel(booking))
}

func (b BookingHandler) GetAllBookings(ctx echo.Context) error {
	reqPage, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid page value")
	}
	reqSize, err := strconv.Atoi(ctx.QueryParam("size"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid size value")
	}
	reqData := struct {
		RequesterID string `json:"requester_id"`
		ProviderID  string `json:"provider_id"`
		ServiceType string `json:"service_type"`
		BookingType string `json:"booking_type"`
		Status      string `json:"status"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	/*
		if err := ctx.Validate(&reqData); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	*/
	if reqData.Status != "" && !util.IsValidBookingStatus(reqData.Status) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid booking status.",
		)
	}
	if reqData.BookingType != "" && !util.IsValidBookingType(reqData.BookingType) {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidBookingType)
	}
	if reqData.ServiceType != "" && !util.IsValidServiceType(reqData.ServiceType) {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidServiceType)
	}

	// ensure that there is at lease an ID, whether requester or provider,
	// else assume request from admin.
	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := b.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	if (reqData.RequesterID != authUser.ID && reqData.ProviderID != authUser.ID) && authUser.Role == models.USER_REGULAR {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"restricted from view bookings.",
		)
	}

	filter := models.BookingFilter{
		RequesterID: reqData.RequesterID,
		ProviderID:  reqData.ProviderID,
		ServiceType: reqData.ServiceType,
		BookingType: reqData.BookingType,
		Status:      reqData.Status,
		Page:        reqPage,
		Limit:       reqSize,
	}

	// ensure proper date and time format
	if reqData.StartDate != "" {
		_, err := util.ParseDate(reqData.StartDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}
	if reqData.EndDate != "" {
		_, err := util.ParseDate(reqData.EndDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}
	if reqData.StartTime != "" {
		_, err := util.ParseTime(
			fmt.Sprintf(
				"%sT%s",
				reqData.EndDate,
				reqData.EndTime,
			),
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}
	if reqData.EndTime != "" {
		_, err := util.ParseTime(
			fmt.Sprintf(
				"%sT%s",
				reqData.EndDate,
				reqData.EndTime,
			),
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	filter.StartDate = reqData.StartDate
	filter.EndDate = reqData.EndDate
	filter.StartTime = reqData.StartTime

	bookings, err := b.app.Repositories.Booking.GetAll(filter)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	resp := []response.BookingResponse{}
	for _, booking := range bookings {
		resp = append(resp, response.BookingResponseFromModel(booking))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (b BookingHandler) UpdateBookingStatus(ctx echo.Context) error {
	bookingID := ctx.Param("id")
	status := ctx.QueryParam("status")

	if !util.IsValidBookingStatus(status) {
		return echo.ErrBadRequest
	}

	booking, err := b.app.Repositories.Booking.GetByID(bookingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				repository.ErrRecordNotFound)
		}
		return util.ErrInternalServer(ctx, err)
	}
	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := b.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	if booking.RequesterID != authUser.ID {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"only requester can update booking status.",
		)
	}

	var notification models.Notification
	switch status {
	case models.BOOKING_COMPLETED:
		if (time.Now().UTC().Hour() - booking.EndTime.UTC().Hour()) < 0 {
			return echo.NewHTTPError(
				http.StatusForbidden,
				"booking still in progress.",
			)
		}
		if booking.Status != models.BOOKING_IN_PROGRESS {
			return echo.NewHTTPError(
				http.StatusForbidden,
				"can only change 'in_progress' booking.",
			)
		}
		notification = models.Notification{
			Type:    models.NOTIFICATION_TYPE_BOOKING,
			UserID:  booking.ProviderID.String,
			Message: "booking completed.",
		}
		notification.BookingID.String = booking.ID
		notification.BookingID.Valid = true
	case models.BOOKING_CANCELED:
		notification = models.Notification{
			Type:    models.NOTIFICATION_TYPE_BOOKING,
			UserID:  booking.ProviderID.String,
			Message: "booking completed.",
		}
		notification.BookingID.String = booking.ID
		notification.BookingID.Valid = true
	default:
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"can only change to 'canceled' or 'completed'",
		)
	}

	_, err = b.app.Repositories.Booking.UpdateStatus(booking.ID, status)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	err = b.app.Repositories.Notification.Create(&notification)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(
		http.StatusOK,
		response.NotificationResponseFromModel(notification))
}

func (b BookingHandler) AcceptBooking(ctx echo.Context) error {
	bookingID := ctx.Param("id")
	booking, err := b.app.Repositories.Booking.GetByID(bookingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				repository.ErrRecordNotFound)
		}
		return util.ErrInternalServer(ctx, err)
	}

	if booking.Status != models.BOOKING_OPEN {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"can only accept open booking.",
		)
	}
	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := b.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}
	if authUser.ServiceRole == models.SERVICE_REQUESTER {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"only providers can accept booking.",
		)
	}

	rejected, err := b.app.Repositories.Booking.CheckRejected(booking.ID, authUser.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	if rejected {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"booking has been previously rejected.",
		)
	}

	// check if user has notification about the booking
	notifications, err := b.app.Repositories.Notification.GetForUser(
		models.NotificationFilter{
			Type:      models.NOTIFICATION_TYPE_BOOKING,
			UserID:    authUser.ID,
			BookingID: bookingID,
			Page:      1,
			Limit:     1,
		},
	)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	if len(notifications) < 1 {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"user does not have access to booking.",
		)
	}

	_, err = b.app.Repositories.Booking.AssignProvider(booking.ID, authUser.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	_, err = b.app.Repositories.Booking.UpdateStatus(
		booking.ID, models.BOOKING_IN_PROGRESS)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	notification := models.Notification{
		Type:    models.NOTIFICATION_TYPE_BOOKING,
		UserID:  booking.RequesterID,
		Message: "booking accepted.",
	}
	notification.BookingID.String = booking.ID
	notification.BookingID.Valid = true
	err = b.app.Repositories.Notification.Create(&notification)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(
		http.StatusOK,
		response.NotificationResponseFromModel(notification))
}

func (b BookingHandler) RejectBooking(ctx echo.Context) error {
	bookingID := ctx.Param("id")
	booking, err := b.app.Repositories.Booking.GetByID(bookingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				repository.ErrRecordNotFound)
		}
		return util.ErrInternalServer(ctx, err)
	}
	if booking.Status != models.BOOKING_OPEN {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"can only reject open booking.",
		)
	}
	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := b.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	if authUser.ServiceRole == models.SERVICE_REQUESTER {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"only providers can reject booking.",
		)
	}

	rejected, err := b.app.Repositories.Booking.CheckRejected(booking.ID, authUser.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	if rejected {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"booking has been previously rejected.",
		)
	}

	// check if user has notification about the booking
	notifications, err := b.app.Repositories.Notification.GetForUser(
		models.NotificationFilter{
			Type:      models.NOTIFICATION_TYPE_BOOKING,
			UserID:    authUser.ID,
			BookingID: bookingID,
			Page:      1,
			Limit:     1,
		},
	)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	if len(notifications) < 1 {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"user does not have access to booking.",
		)
	}

	// add entry to rejected booking
	err = b.app.Repositories.Booking.RejectBooking(booking.ID, authUser.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	notification := models.Notification{
		Type:    models.NOTIFICATION_TYPE_BOOKING,
		UserID:  booking.RequesterID,
		Message: "booking rejected.",
	}
	notification.BookingID.String = booking.ID
	notification.BookingID.Valid = true
	err = b.app.Repositories.Notification.Create(&notification)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(
		http.StatusOK,
		response.NotificationResponseFromModel(notification))
}

func (b BookingHandler) FindProviders(ctx echo.Context) error {
	bookingID := ctx.Param("id")
	reqPage, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid page value")
	}
	reqSize, err := strconv.Atoi(ctx.QueryParam("size"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid size value")
	}

	booking, err := b.app.Repositories.Booking.GetByID(bookingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				repository.ErrRecordNotFound)
		}
		return util.ErrInternalServer(ctx, err)
	}
	if booking.Status != models.BOOKING_OPEN {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"can only find providers for open booking.",
		)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := b.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	if authUser.ID != booking.RequesterID {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"only booking requester can search for providers.",
		)
	}

	services, err := b.app.Repositories.Booking.MatchServices(
		&booking,
		models.ServiceFilter{
			Page:  reqPage,
			Limit: reqSize,
		},
	)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	distEstimations := []response.DistanceEstimation{}
	for _, service := range services {
		distance, err := fetchDistanceInfo(
			service.Address,
			booking.RequesterAddr,
			b.app.Config.Google.MapSecret,
		)
		if err != nil || distance.Status != "OK" {
			distEstimations = append(
				distEstimations,
				response.DistanceEstimation{
					Distance: "unknown",
					Duration: "unknown",
				},
			)
		} else {
			distEstimations = append(
				distEstimations,
				response.DistanceEstimation{
					Distance: distance.Rows[0].Elements[0].Distance.Text,
					Duration: distance.Rows[0].Elements[0].Duration.Text,
				},
			)
		}
	}

	resp := []response.ServiceDistanceResponse{}
	for idx, service := range services {
		resp = append(
			resp,
			response.BookingServiceProviderResponse(service, distEstimations[idx]),
		)
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (b BookingHandler) SelectProvider(ctx echo.Context) error {
	reqData := struct {
		ProviderID string `json:"provider_id" validate:"required"`
		BookingID  string `json:"booking_id" validate:"required"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	booking, err := b.app.Repositories.Booking.GetByID(reqData.BookingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				repository.ErrRecordNotFound)
		}
		return util.ErrInternalServer(ctx, err)
	}
	if booking.Status != models.BOOKING_OPEN {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"can only select provider for open booking.",
		)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := b.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}
	if authUser.ID == reqData.ProviderID {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"attempting to select self as provider.",
		)
	}
	serviceProvider, err := b.app.Repositories.User.GetByID(reqData.ProviderID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}
	if serviceProvider.ServiceRole == models.SERVICE_REQUESTER {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"attempting to select requester as provider.",
		)
	}

	notification := models.Notification{
		Type:   models.NOTIFICATION_TYPE_BOOKING,
		UserID: reqData.ProviderID,
		Message: fmt.Sprintf(
			BOOKING_PROVIDER_SELECTION_MESSAGE,
			booking.ID,
		),
	}
	notification.BookingID.String = booking.ID
	notification.BookingID.Valid = true
	err = b.app.Repositories.Notification.Create(&notification)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}
	return ctx.JSON(
		http.StatusOK,
		response.NotificationResponseFromModel(notification))
}

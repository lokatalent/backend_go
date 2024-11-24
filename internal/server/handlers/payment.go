package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	// "github.com/lokatalent/backend_go/cmd/api/models/response"
	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

type PaymentHandler struct {
	app *util.Application
}

func NewPaymentHandler(app *util.Application) PaymentHandler {
	return PaymentHandler{app: app}
}

func (p PaymentHandler) InitializeTransaction(ctx echo.Context) error {
	reqData := struct {
		BookingID   string `json:"booking_id" validate:"required"`
		CallBackURL string `json:"callback_url" validate:"required"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := p.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	var booking models.Booking
	// var err error
	if reqData.BookingID != "" {
		booking, err = p.app.Repositories.Booking.GetByID(reqData.BookingID)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				return echo.ErrNotFound
			}
			return util.ErrInternalServer(ctx, err)
		}
	}
	if booking.RequesterID != authUser.ID {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"only requester can pay for open booking.",
		)
	}

	var accessCode string
	payment, err := p.app.Repositories.Payment.GetPayment(models.PaymentFilter{
		Type:      models.PAYMENT_TYPE_CREDIT,
		BookingID: reqData.BookingID,
		Status:    models.PAYMENT_STATUS_PENDING,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			paymentRef := uuid.NewString()
			accessCode, err = initTransaction(
				authUser.Email,
				paymentRef,
				reqData.CallBackURL,
				booking.TotalPrice,
				p.app.Config.Paystack.APIKey,
			)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}

			newPayment := models.Payment{
				ID:         uuid.NewString(),
				Type:       models.PAYMENT_TYPE_CREDIT,
				PaymentRef: paymentRef,
				Amount:     booking.TotalPrice,
				Status:     models.PAYMENT_STATUS_PENDING,
			}
			newPayment.BookingID.String = booking.ID
			newPayment.BookingID.Valid = true
			err = p.app.Repositories.Payment.CreatePayment(&newPayment)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}

			err = p.app.Repositories.Payment.CreateAccessCode(
				uuid.NewString(),
				newPayment.ID,
				accessCode,
			)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}
		} else {
			return util.ErrInternalServer(ctx, err)
		}
	} else {
		accessCode, err = p.app.Repositories.Payment.GetAccessCode(payment.ID)
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
		return ctx.JSON(http.StatusProcessing, accessCode)
	}
	return ctx.JSON(http.StatusOK, accessCode)
}

func (p PaymentHandler) VerifyTransaction(ctx echo.Context) error {
	authenticatedUser := util.ContextGetUser(ctx)
	reqData := struct {
		BookingID string `json:"booking_id" validate:"required"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	booking, err := p.app.Repositories.Booking.GetByID(reqData.BookingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	var payment models.Payment
	if authenticatedUser.ID == booking.RequesterID {
		payment, err = p.app.Repositories.Payment.GetPayment(models.PaymentFilter{
			Type:      models.PAYMENT_TYPE_CREDIT,
			BookingID: reqData.BookingID,
		})
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				return echo.NewHTTPError(
					http.StatusNotFound,
					"no payment record found for booking.",
				)
			}
			return util.ErrInternalServer(ctx, err)
		}
	} else if authenticatedUser.ID == booking.ProviderID.String {
		payment, err = p.app.Repositories.Payment.GetPayment(models.PaymentFilter{
			Type:      models.PAYMENT_TYPE_DEBIT,
			BookingID: reqData.BookingID,
		})
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				return echo.NewHTTPError(
					http.StatusNotFound,
					"no payment record found for booking.",
				)
			}
			return util.ErrInternalServer(ctx, err)
		}
	}

	fmt.Println(payment)
	if payment.Status == models.PAYMENT_STATUS_VERIFIED {
		return ctx.JSON(http.StatusOK, models.PAYMENT_STATUS_VERIFIED)
	}

	var status string
	if authenticatedUser.ID == booking.RequesterID {
		status, err = verifyTransaction(payment.PaymentRef, p.app.Config.Paystack.APIKey)
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
	}

	if authenticatedUser.ID == booking.ProviderID.String {
		status, err = verifyTransfer(payment.PaymentRef, p.app.Config.Paystack.APIKey)
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
	}

	switch status {
	case "success":
		_, err = p.app.Repositories.Payment.UpdatePaymentStatus(payment.ID, models.PAYMENT_STATUS_VERIFIED)
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
		// update wallet
		if authenticatedUser.ID == booking.RequesterID {
			err = p.app.Repositories.Payment.UpdateWallet(
				booking.RequesterID,
				models.PAYMENT_TYPE_CREDIT,
				booking.TotalPrice,
			)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}
		}
		return ctx.JSON(http.StatusOK, models.PAYMENT_STATUS_VERIFIED)
	case "abandoned", "failed", "reversed":
		err = p.app.Repositories.Payment.DeleteAccessCode(payment.ID)
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
		_, err = p.app.Repositories.Payment.UpdatePaymentStatus(payment.ID, models.PAYMENT_STATUS_CANCELED)
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, status)
}

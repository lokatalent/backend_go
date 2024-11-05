package handlers

import (
	"errors"
	// "math"
	"time"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/models"
)

const (
	DAY_IN_HOURS              = 24 // numbers of hours in a day.
	MAX_BOOKING_DAYS          = 30 // maximum numbers of booking days.
	MAX_INSTANT_BOOKING_HOURS = 10 // maximum 10 hours a day.
	// MAX_SCHEDULED_BOOKING_HOURS = 300   // maximum of 30 days and 10 hours a day.
)

func calculateBookingPrice(
	app *util.Application,
	serviceType, bookingType string,
	startDate, endDate, startTime, endTime time.Time,
) (float64, error) {
	servicePrice, err := app.Repositories.ServicePricing.GetServicePricing(serviceType)
	if err != nil {
		return float64(0.0), err
	}

	var hours float64
	switch bookingType {
	case models.BOOKING_INSTANT:
		hours = endTime.Sub(startTime).Hours()
		if hours < float64(0.0) || hours > float64(MAX_INSTANT_BOOKING_HOURS) {
			return float64(0.0), errors.New("invalid start and end time.")
		}
		return servicePrice.RatePerHour * hours, nil
	case models.BOOKING_SCHEDULED:
		days := endDate.Sub(startDate).Hours() / DAY_IN_HOURS
		if days < float64(0.0) || days > MAX_BOOKING_DAYS {
			return float64(0.0), errors.New("invalid start and end days.")
		}
		hours = endTime.Sub(startTime).Hours()
		if hours < float64(0.0) || hours > float64(MAX_INSTANT_BOOKING_HOURS) {
			return float64(0.0), errors.New("invalid start and end time.")
		}
		return servicePrice.RatePerHour * hours * days, nil

	default:
		return 0.0, errors.New("invalid booking type!")
	}
}

package util

import (
	"errors"
	"fmt"
	"math/rand"
	"net/mail"
	"net/textproto"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/internal/models"
)

// ValidateContentType verifies that request header content-type is supported
func ValidateContentType(header textproto.MIMEHeader, validTypes ...string) (string, error) {
	errMsg := fmt.Sprintf("%s Expected one of:", ErrInvalidContentType)
	contentType := header.Get(echo.HeaderContentType)

	for _, validType := range validTypes {
		if strings.EqualFold(contentType, validType) {
			return validType, nil
		}
		errMsg = fmt.Sprintf("%s %s", errMsg, validType)
	}

	return "", errors.New(errMsg)
}

// IsValidUUID validates given UUID
func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// IsValidEmail validates the structure of a given email address
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsValidPhoneNumber validates the structure of a given phone number
func IsValidPhoneNumber(num string) bool {
	matched, err := regexp.MatchString(PhoneNumPattern, num)
	if err != nil {
		return false
	}

	return matched
}

// IsAdmin checks if a user has admin priviledge
func IsAdmin(userRole string) bool {
	switch userRole {
	case models.USER_ADMIN, models.USER_ADMIN_SUPER:
		return true
	default:
		return false
	}
}

// IsValidServiceRole checks if the specified role is a valid
// service role.
func IsValidServiceRole(role string) bool {
	switch role {
	case models.SERVICE_PROVIDER, models.SERVICE_REQUESTER, models.SERVICE_BOTH:
		return true
	default:
		return false
	}
}

func IsValidServiceType(serviceType string) bool {
	switch serviceType {
	case models.SERVICE_CLEANING, models.SERVICE_PLUMBING, models.SERVICE_COOKING:
		return true
	default:
		return false
	}
}

func IsValidBookingType(bookingType string) bool {
	switch bookingType {
	case models.BOOKING_INSTANT, models.BOOKING_SCHEDULED:
		return true
	default:
		return false
	}
}

func IsValidBookingStatus(status string) bool {
	switch status {
	case models.BOOKING_OPEN, models.BOOKING_CANCELED, models.BOOKING_IN_PROGRESS,
		models.BOOKING_COMPLETED, models.BOOKING_ACCEPT, models.BOOKING_REJECT:
		return true
	default:
		return false
	}
}

func IsValidNotificationType(notificationType string) bool {
	switch notificationType {
	case models.NOTIFICATION_TYPE_BOOKING:
		return true
	default:
		return false
	}
}

// ParseBool converts boolean values in string
func ParseBool(value string) (bool, error) {
	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, errors.New("Invalid boolean string")
	}
}

// RandomInt generates a radom integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func ValidVerificationType(value string) bool {
	switch value {
	case models.EMAIL_VERIFICATION, models.PHONE_VERIFICATION:
		return true
	default:
		return false
	}
}

// checks if a place address satisfy the expected format.
func ValidPlaceAddress(address string) bool {
	return len(strings.Split(address, ",")) >= 4
}

func ParseDate(date string) (time.Time, error) {
	parsedTime, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return time.Time{}, errors.New("invalid date format. Expected YYYY-MM-DD")
	}

	return parsedTime, nil
}

func ParseTime(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, errors.New("invalid time format. Expected HH:MM:SS+ZZ:ZZ")
	}

	return parsedTime, nil
}

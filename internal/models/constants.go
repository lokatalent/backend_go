package models

// user roles
const (
	USER_REGULAR     = "regular"
	USER_ADMIN       = "admin"
	USER_ADMIN_SUPER = "admin_super"
)

const (
	EMAIL_VERIFICATION = "email"
	PHONE_VERIFICATION = "phone"
)

// user service roles
const (
	SERVICE_PROVIDER  = "service_provider"
	SERVICE_REQUESTER = "service_requester"
	SERVICE_BOTH      = "service_both"
)

// user service types
const (
	SERVICE_CLEANING = "cleaning"
	SERVICE_PLUMBING = "plumbing"
	SERVICE_COOKING  = "cooking"
)

// bookings
const (
	// booking status
	BOOKING_OPEN        = "open"
	BOOKING_COMPLETED   = "completed"
	BOOKING_IN_PROGRESS = "in_progress"
	BOOKING_CANCELED    = "canceled"
	BOOKING_ACCEPT      = "accept"
	BOOKING_REJECT      = "reject"

	// booking type
	BOOKING_INSTANT   = "instant"
	BOOKING_SCHEDULED = "scheduled"
)

// notifications
const (
	NOTIFICATION_TYPE_BOOKING = "booking"
)

const (
	DefaultPage      = 1
	DefaultPageLimit = 10
)

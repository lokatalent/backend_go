package util

import "time"

// application environment
const (
	ENVIRONMENT_DEVELOPMENT = "DEVELOPMENT"
	ENVIRONMENT_STAGING     = "STAGING"
	ENVIRONMENT_PRODUCTION  = "PRODUCTION"
)

// Database connection string format
const DB_CONN_FMT = "postgres://%s:%s@%s:%s/%s?sslmode=disable"

const ContextKeyUser = "user"

// tokens configuration
const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 24 * time.Hour
)

// Content-Type configurations
const (
	ContentTypeJPEG = "image/jpeg"
	ContentTypePNG  = "image/png"

	// Errors
	ErrInvalidContentType = "Unexpected content type."
)

const PhoneNumPattern = `^\+234[789]\d{9}$`

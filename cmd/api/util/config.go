package util

import (
	// "errors"
	"fmt"
	"os"
	"strconv"
	// "strings"
)

type SendGridSecret struct {
	APIKey string
	Sender string
}

type AWSSecret struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	S3Bucket        string
}

type GoogleSecret struct {
	ClientID     string
	ClientSecret string
	MapSecret    string
}

// JSON Web Token
type JWTSecret struct {
	Access  string
	Refresh string
}

type TwilioSecret struct {
	APIKey     string
	APISecret  string
	AccountSID string
	Sender     string
}

// Config holds configuration data loaded from .env file.
type Config struct {
	Env string
	// CertPath    string
	// CertKeyPath string
	Port   int
	Origin string

	JWT JWTSecret

	DB struct {
		DSN string
	}

	AWS      AWSSecret
	Google   GoogleSecret
	SendGrid SendGridSecret
	Twilio   TwilioSecret
}

// Load reads in all required environment variable to start the
// application.
func (c *Config) Load() error {
	var err error
	if c.Env, err = loadAppEnv(); err != nil {
		return err
	}
	/*
		if c.CertPath, err = loadAppCertPath(); err != nil {
			return err
		}
		if c.CertKeyPath, err = loadAppCertKeyPath(); err != nil {
			return err
		}
	*/

	if c.Origin, err = loadAppOrigin(); err != nil {
		return err
	}

	if c.Port, err = loadAppPort(); err != nil {
		return err
	}

	if err := loadJWTSecrets(&c.JWT); err != nil {
		return err
	}

	if err := loadGoogleSecrets(&c.Google); err != nil {
		return err
	}

	if c.DB.DSN, err = loadDB(); err != nil {
		return err
	}

	if err := loadAWSSecrets(&c.AWS); err != nil {
		return err
	}

	if err := loadSendGridSecrets(&c.SendGrid); err != nil {
		return err
	}

	if err := loadTwilioSecrets(&c.Twilio); err != nil {
		return err
	}

	return nil
}

// loadAppEnv loads application environment.
func loadAppEnv() (string, error) {
	env, ok := os.LookupEnv("APP_ENV")
	if !ok {
		return "", missingEnvVar("APP_ENV")
	}

	switch env {
	case ENVIRONMENT_PRODUCTION, ENVIRONMENT_DEVELOPMENT, ENVIRONMENT_STAGING:
		return env, nil
	default:
		return "", invalidEnvVar("APP_ENV", "PRODUCTION|DEVELOPMENT", env)
	}
}

/*
// loadAppEnv loads application SSL certificate path.
func loadAppCertPath() (string, error) {
	env, ok := os.LookupEnv("APP_CERT_PATH")
	if !ok {
		return "", missingEnvVar("APP_CERT_PATH")
	}

	return env, nil
}

// loadAppEnv loads application SSL certificate key path.
func loadAppCertKeyPath() (string, error) {
	env, ok := os.LookupEnv("APP_CERT_KEY_PATH")
	if !ok {
		return "", missingEnvVar("APP_CERT_KEY_PATH")
	}

	return env, nil
}
*/

// loadAppOigin loads application origin.
func loadAppOrigin() (string, error) {
	origin, ok := os.LookupEnv("ORIGIN")
	if !ok {
		return "", missingEnvVar("ORIGIN")
	}
	return origin, nil
}

// loadAppPort loads application listening port.
func loadAppPort() (int, error) {
	portEnv, ok := os.LookupEnv("APP_PORT")
	if !ok {
		return 0, missingEnvVar("APP_PORT")
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		return 0, invalidEnvVar("APP_PORT", "integer", portEnv)
	}

	return port, nil
}

// loadDB loads database connection parameters.
func loadDB() (string, error) {
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return "", missingEnvVar("DB_NAME")
	}
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return "", missingEnvVar("DB_HOST")
	}
	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return "", missingEnvVar("DB_PORT")
	}
	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		return "", missingEnvVar("DB_USER")
	}
	dbPass, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return "", missingEnvVar("DB_PASSWORD")
	}

	appEnv, err := loadAppEnv()
	if err != nil {
		return "", err
	}

	switch appEnv {
	case ENVIRONMENT_DEVELOPMENT:
		return fmt.Sprintf(DB_CONN_FMT_TEST, dbUser, dbPass, dbHost, dbPort, dbName), nil
	default:
		return fmt.Sprintf(DB_CONN_FMT, dbUser, dbPass, dbHost, dbPort, dbName), nil
	}
}

// loadJWTSecrets loads secrets for generation JSON web tokens.
func loadJWTSecrets(jwt *JWTSecret) error {
	jwtAccess, ok := os.LookupEnv("JWT_ACCESS")
	if !ok {
		return missingEnvVar("JWT_ACCESS")
	}
	jwtRefresh, ok := os.LookupEnv("JWT_REFRESH")
	if !ok {
		return missingEnvVar("JWT_REFRESH")
	}

	// validate secrets
	if len(jwtAccess) != 64 {
		return invalidEnvVar("JWT_ACCESS", "string of length 64", jwtAccess)
	}
	if len(jwtRefresh) != 64 {
		return invalidEnvVar("JWT_REFRESH", "string of length 64", jwtRefresh)
	}

	jwt.Access = jwtAccess
	jwt.Refresh = jwtRefresh

	return nil
}

// loadGoogleSecrets loads secrets for Google API
func loadGoogleSecrets(google *GoogleSecret) error {
	googleClientID, ok := os.LookupEnv("GOOGLE_CLIENT_ID")
	if !ok {
		return missingEnvVar("GOOGLE_CLIENT_ID")
	}
	googleClientSecret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET")
	if !ok {
		return missingEnvVar("GOOGLE_CLIENT_SECRET")
	}
	googleMapSecret, ok := os.LookupEnv("GOOGLE_MAP_API_KEY")
	if !ok {
		return missingEnvVar("GOOGLE_MAP_API_KEY")
	}

	// validate secrets
	if len(googleClientID) < 1 {
		return invalidEnvVar(
			"GOOGLE_CLIENT_ID", "string of length > 1", googleClientID)
	}
	if len(googleClientSecret) < 1 {
		return invalidEnvVar(
			"GOOGLE_CLIENT_SECRET", "string of length > 1", googleClientSecret)
	}
	if len(googleMapSecret) < 1 {
		return invalidEnvVar(
			"GOOGLE_MAP_API_KEY", "string of length > 1", googleMapSecret)
	}

	google.ClientID = googleClientID
	google.ClientSecret = googleClientSecret
	google.MapSecret = googleMapSecret

	return nil
}

// loadAWSSecrets loads secrets for AWS API.
func loadAWSSecrets(aws *AWSSecret) error {
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return missingEnvVar("AWS_REGION")
	}
	awsAccessKeyId, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok {
		return missingEnvVar("AWS_ACCESS_KEY_ID")
	}
	awsSecretAccessKey, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok {
		return missingEnvVar("AWS_SECRET_ACCESS_KEY")
	}
	awsS3Bucket, ok := os.LookupEnv("AWS_S3_BUCKET")
	if !ok {
		return missingEnvVar("AWS_S3_BUCKET")
	}

	// validate secrets
	if len(awsRegion) < 1 {
		return invalidEnvVar("AWS_REGION", "string of length > 1", awsRegion)
	}
	if len(awsAccessKeyId) < 1 {
		return invalidEnvVar(
			"AWS_ACCESS_KEY_ID", "string of length > 1", awsAccessKeyId)
	}
	if len(awsSecretAccessKey) < 1 {
		return invalidEnvVar(
			"AWS_SECRET_ACCESS_KEY", "string of length > 1", awsSecretAccessKey)
	}
	if len(awsS3Bucket) < 1 {
		return invalidEnvVar(
			"AWS_S3_BUCKET", "string of length > 1", awsS3Bucket)
	}

	aws.Region = awsRegion
	aws.AccessKeyId = awsAccessKeyId
	aws.SecretAccessKey = awsSecretAccessKey
	aws.S3Bucket = awsS3Bucket

	return nil
}

// loadSendGridSecrets loads secrets for SendGrid APIs.
func loadSendGridSecrets(sendgrid *SendGridSecret) error {
	sendGridAPIKey, ok := os.LookupEnv("SENDGRID_API_KEY")
	if !ok {
		return missingEnvVar("SENDGRID_API_KEY")
	}
	sendGridSender, ok := os.LookupEnv("SENDGRID_SENDER")
	if !ok {
		return missingEnvVar("SENDGRID_SENDER")
	}

	// validate secrets
	if len(sendGridAPIKey) < 1 {
		return invalidEnvVar(
			"SENDGRID_API_KEY", "string of length > 1", sendGridAPIKey,
		)
	}
	if len(sendGridSender) < 1 {
		return invalidEnvVar(
			"SENDGRID_SENDER", "string of length > 1", sendGridSender,
		)
	}

	sendgrid.APIKey = sendGridAPIKey
	sendgrid.Sender = sendGridSender

	return nil
}

// loadTwilioSecrets loads secrets for Twilio APIs.
func loadTwilioSecrets(twilio *TwilioSecret) error {
	twilioAPIKey, ok := os.LookupEnv("TWILIO_API_KEY")
	if !ok {
		return missingEnvVar("TWILIO_API_KEY")
	}
	twilioSender, ok := os.LookupEnv("TWILIO_SENDER")
	if !ok {
		return missingEnvVar("TWILIO_SENDER")
	}
	twilioAPISecret, ok := os.LookupEnv("TWILIO_API_SECRET")
	if !ok {
		return missingEnvVar("TWILIO_API_SECRET")
	}
	twilioAccountSID, ok := os.LookupEnv("TWILIO_ACCOUNT_SID")
	if !ok {
		return missingEnvVar("TWILIO_ACCOUNT_SID")
	}

	// validate secrets
	if len(twilioAPIKey) < 1 {
		return invalidEnvVar(
			"TWILIO_API_KEY", "string of length > 1", twilioAPIKey,
		)
	}
	if len(twilioSender) < 1 {
		return invalidEnvVar(
			"TWILIO_SENDER", "string of length > 1", twilioSender,
		)
	}
	if len(twilioAPISecret) < 1 {
		return invalidEnvVar(
			"TWILIO_API_SECRET", "string of length > 1", twilioAPISecret,
		)
	}
	if len(twilioAccountSID) < 1 {
		return invalidEnvVar(
			"TWILIO_ACCOUNT_SID", "string of length > 1", twilioAccountSID,
		)
	}

	twilio.APIKey = twilioAPIKey
	twilio.Sender = twilioSender
	twilio.AccountSID = twilioAccountSID
	twilio.APISecret = twilioAPISecret

	return nil
}

// missingEnvVar reports missing environment variable.
func missingEnvVar(envVar string) error {
	return fmt.Errorf("missing environment var: %s\n", envVar)
}

func invalidEnvVar(envVar, expected, got string) error {
	return fmt.Errorf(
		"invalid value for %s, expected: '%s' got: '%s'.\n",
		envVar, expected, got)
}

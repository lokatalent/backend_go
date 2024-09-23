package util

import (
	// "errors"
	"fmt"
	"os"
	"strconv"
	// "strings"
)

type AWSSecret struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	S3Bucket        string
}

// JSON Web Token
type JWTSecret struct {
	Access  string
	Refresh string
}

// Config holds configuration data loaded from .env file.
type Config struct {
	Env  string
	Port int

	JWT JWTSecret

	DB struct {
		DSN string
	}

	AWS AWSSecret
}

// Load reads in all required environment variable to start the
// application.
func (c *Config) Load() error {
	var err error
	if c.Env, err = loadAppEnv(); err != nil {
		return err
	}

	if c.Port, err = loadAppPort(); err != nil {
		return err
	}

	if err := loadJWTSecrets(&c.JWT); err != nil {
		return err
	}

	if c.DB.DSN, err = loadDB(); err != nil {
		return err
	}

	if err := loadAWSSecrets(&c.AWS); err != nil {
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
	case "PRODUCTION", "DEVELOPMENT":
		return env, nil
	default:
		return "", invalidEnvVar("APP_ENV", "PRODUCTION|DEVELOPMENT", env)
	}
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

	return fmt.Sprintf(DB_CONN_FMT, dbUser, dbPass, dbHost, dbPort, dbName), nil
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

// loadAWSSecrets loads secrets for generation JSON web tokens.
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

// missingEnvVar reports missing environment variable.
func missingEnvVar(envVar string) error {
	return fmt.Errorf("missing environment var: %s\n", envVar)
}

func invalidEnvVar(envVar, expected, got string) error {
	return fmt.Errorf(
		"invalid value for %s, expected: '%s' got: '%s'.\n",
		envVar, expected, got)
}

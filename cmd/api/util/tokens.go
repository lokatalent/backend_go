package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/lokatalent/backend_go/internal/models"
)

type CustomAccessJWTClaims struct {
	Email    string `json:"email,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	jwt.RegisteredClaims
}

type CustomRefreshJWTClaims struct {
	Email string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

// GenerateTokens generates the signed access and refresh tokens, with
// expiration time.
func GenerateTokens(app *Application, user *models.User) (string, string, int64, error) {
	accessTokenExpiration := time.Now().Add(AccessTokenDuration)
	accessClaims := CustomAccessJWTClaims{
		Email:    user.Email,
		Verified: user.IsVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        user.ID,
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		},
	}

	// generate refresh token with lifetime of 1 day
	refreshClaims := CustomAccessJWTClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration)),
		},
	}

	accessToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessClaims,
	).SignedString(
		[]byte(app.Config.JWT.Access))

	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	).SignedString(
		[]byte(app.Config.JWT.Refresh))

	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, accessTokenExpiration.Unix(), nil
}

// ValidateRefreshToken validates the provided JWT.
func ValidateRefreshToken(secret, signedToken string) (*CustomRefreshJWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&CustomRefreshJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*CustomRefreshJWTClaims)
	if !ok {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}

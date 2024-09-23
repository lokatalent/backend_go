package handlers

import (
	"errors"
	"net/http"
	// "strings"

	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/models/response"
	"github.com/lokatalent/backend_go/cmd/api/util"
	// "github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

type AuthHandler struct {
	app *util.Application
}

func NewAuthHandler(app *util.Application) AuthHandler {
	return AuthHandler{app: app}
}

// RefreshToken refreshes JWT access token attached to a user session.
func (a AuthHandler) RefreshToken(ctx echo.Context) error {
	reqBody := struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}{}

	if err := ctx.Bind(&reqBody); err != nil {
		return echo.ErrBadRequest
	}

	// validate refresh token
	claims, err := util.ValidateRefreshToken(a.app.Config.JWT.Refresh, reqBody.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid or expired refresh token")
	}

	// retrieve user details to generate new access token
	user, err := a.app.Repositories.User.GetByEmail(claims.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	accessToken, _, expiration, err := util.GenerateTokens(a.app, &user)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	resp := response.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: reqBody.RefreshToken,
		ExpiresAt:    expiration,
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (a AuthHandler) VerifyUser(ctx echo.Context) error {
	reqBody := struct {
		Email    string `json:"email" validate:"required,email"`
		PhoneNum string `json:"phone_num" validate:"required,len=14"`
		Status   bool   `json:"verification_status" validate:"required"`
	}{}

	if err := ctx.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	fetchedUser, err := a.app.Repositories.User.GetByEmail(authenticatedUser.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	// prevent multiple verification
	if fetchedUser.IsVerified {
		return echo.NewHTTPError(
			http.StatusConflict,
			"user is already verified!")
	}

	if !util.IsValidEmail(reqBody.Email) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid email address!")
	}
	if !util.IsValidPhoneNumber(reqBody.PhoneNum) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid phone number!")
	}

	err = a.app.Repositories.User.Verify(authenticatedUser.ID, reqBody.Status)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	} else {
		fetchedUser.IsVerified = reqBody.Status
	}

	// generate new access and refresh tokens
	accessToken, refreshToken, expiration, err := util.GenerateTokens(a.app, &fetchedUser)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	resp := response.AuthResponse{
		TokensResponse: response.TokensResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiration,
		},
		UserResponse: response.UserResponseFromModel(&fetchedUser),
	}

	return ctx.JSON(http.StatusOK, resp)
}

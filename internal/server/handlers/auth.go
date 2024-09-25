package handlers

import (
	"errors"
	"fmt"
	"net/http"
	// "strings"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"

	"github.com/lokatalent/backend_go/cmd/api/models/response"
	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
	"github.com/lokatalent/backend_go/internal/server/egothic"
)

type AuthHandler struct {
	app *util.Application
}

func NewAuthHandler(app *util.Application) AuthHandler {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options.HttpOnly = true
	store.Options.Secure = (app.Config.Env == util.ENVIRONMENT_PRODUCTION)
	egothic.SetStore(store)

	// Setup third-parties authentication handler
	goth.UseProviders(
		google.New(
			app.Config.Google.ClientID,
			app.Config.Google.ClientSecret,
			fmt.Sprintf("http://localhost:%d/auth/google/callback", app.Config.Port),
			"email",
			"profile",
			"openid",
		),
	)
	return AuthHandler{app: app}
}

// SignUp creates a new user, using the details from the form data
func (a AuthHandler) SignUp(ctx echo.Context) error {
	reqData := struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=8"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	hashedPassword, err := util.HashPassword(reqData.Password)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	user, err := a.app.Repositories.User.Create(models.User{
		FirstName: reqData.FirstName,
		LastName:  reqData.LastName,
		Email:     reqData.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateDetails) {
			return echo.NewHTTPError(
				http.StatusForbidden,
				"user with email already exists!")
		}
		return util.ErrInternalServer(ctx, err)
	}

	/*
		// generate access and refresh tokens
		accessToken, refreshToken, expiration, err := util.GenerateTokens(a.app, &user)
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}

		resp := response.AuthResponse{
			TokensResponse: response.TokensResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				ExpiresAt:    expiration,
			},
			UserResponse: response.UserResponseFromModel(&user),
		}
	*/
	return ctx.JSON(http.StatusOK, response.PublicUserResponseFromModel(&user))
}

// SignIn authenticate an already existing user, using the details
// from the formdata
func (a AuthHandler) SignIn(ctx echo.Context) error {
	reqData := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	user, err := a.app.Repositories.User.GetByEmail(reqData.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"user does not exist!")
		}
		return util.ErrInternalServer(ctx, err)
	}

	// validate password
	err = util.ValidatePassword(reqData.Password, user.Password)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"invalid password!")
	}

	// generate access and refresh tokens
	accessToken, refreshToken, expiration, err := util.GenerateTokens(a.app, &user)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	resp := response.AuthResponse{
		TokensResponse: response.TokensResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiration,
		},
		UserResponse: response.UserResponseFromModel(&user),
	}
	return ctx.JSON(http.StatusOK, resp)
}

// ProviderAuthentication handles third-party authentication.
func (a AuthHandler) ProviderAuthentication(ctx echo.Context) error {
    err := egothic.BeginAuthHandler(ctx)
    if err != nil {
        return util.ErrInternalServer(ctx, err)
    }
	return nil
}

func (a AuthHandler) ProviderAuthCallback(ctx echo.Context) error {
	user, err := egothic.CompleteUserAuth(ctx)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	// retrieve an existing user, or create an account for a new user
	fetchedUser, err := a.app.Repositories.User.GetByEmail(user.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			// create a new user
			newUser := models.User{
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Avatar:    user.AvatarURL,
			}
			fetchedUser, err = a.app.Repositories.User.Create(newUser)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}
		} else {
			return util.ErrInternalServer(ctx, err)
		}
	}

	// generate access and refresh tokens
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

// RefreshToken refreshes JWT access token attached to a user session.
func (a AuthHandler) RefreshToken(ctx echo.Context) error {
	reqData := struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// validate refresh token
	claims, err := util.ValidateRefreshToken(a.app.Config.JWT.Refresh, reqData.RefreshToken)
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
		RefreshToken: reqData.RefreshToken,
		ExpiresAt:    expiration,
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (a AuthHandler) VerifyUser(ctx echo.Context) error {
	reqData := struct {
		Email    string `json:"email" validate:"required,email"`
		PhoneNum string `json:"phone_num" validate:"required,len=14"`
		Status   bool   `json:"status" validate:"required"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(&reqData); err != nil {
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

	if !util.IsValidEmail(reqData.Email) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid email address!")
	}
	if !util.IsValidPhoneNumber(reqData.PhoneNum) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid phone number!")
	}

	err = a.app.Repositories.User.Verify(authenticatedUser.ID, reqData.Status)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	} else {
		fetchedUser.IsVerified = reqData.Status
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

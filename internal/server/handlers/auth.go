package handlers

import (
	"errors"
	"fmt"
	"net/http"
	// "strings"
	"os"
	"time"

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

const (
	verificationCodeMin   = 100000
	verificationCodeMax   = 999999
	emailVerificationTmpl = "email_verification.gotmpl"
	phoneVerificationTmpl = "phone_verification.gotmpl"

	verificationTypeEmail = "email"
	verificationTypePhone = "phone"
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
			fmt.Sprintf(
				"http://%s/auth/google/callback",
				app.Config.Origin,
			),
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
				repository.ErrDuplicateDetails)
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
				repository.ErrRecordNotFound)
		}
		return util.ErrInternalServer(ctx, err)
	}

	// validate password
	err = util.ValidatePassword(reqData.Password, user.Password)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			ErrInvalidPassword)
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
			ErrInvalidToken)
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
			http.StatusForbidden,
			ErrAlreadyVerified)
	}

	if !fetchedUser.EmailVerified {
		return echo.NewHTTPError(
			http.StatusFailedDependency,
			ErrVerificationDependency)
	}
	if !fetchedUser.PhoneVerified {
		return echo.NewHTTPError(
			http.StatusFailedDependency,
			ErrVerificationDependency)
	}

	err = a.app.Repositories.User.Verify(authenticatedUser.ID, true)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	} else {
		fetchedUser.IsVerified = true
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

func (a AuthHandler) VerifyContact(ctx echo.Context) error {
	reqVerificationType := ctx.QueryParam("verification_type")

	if !util.ValidVerificationType(reqVerificationType) {
		return echo.ErrBadRequest
	}

	authenticatedUser := util.ContextGetUser(ctx)
	fetchedUser, err := a.app.Repositories.User.GetByEmail(authenticatedUser.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	// prevent multiple email verification
	if (reqVerificationType == verificationTypeEmail) && fetchedUser.EmailVerified {
		return echo.NewHTTPError(
			http.StatusForbidden,
			ErrAlreadyVerified)
	}
	if (reqVerificationType == verificationTypePhone) && fetchedUser.PhoneVerified {
		return echo.NewHTTPError(
			http.StatusForbidden,
			ErrAlreadyVerified)
	}

	// ensure no multiple active verification code
	activeCode, err := a.app.Repositories.User.GetVerificationCode(
		fetchedUser.ID, reqVerificationType)
	if err != nil {
		// do not flag not found error
		if !errors.Is(err, repository.ErrRecordNotFound) {
			return util.ErrInternalServer(ctx, err)
		}
	} else {
		switch time.Now().Compare(activeCode.ExpiresAt) {
		case -1:
			// code is still active
			return echo.NewHTTPError(
				http.StatusForbidden,
				"verification code has already been sent.",
			)
		case 0, 1:
			// remove expired code and proceed to send new one.
			err = a.app.Repositories.User.DeleteVerificationCode(
				fetchedUser.ID, reqVerificationType)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}
		}
	}

	verificationCode := util.RandomInt(verificationCodeMin, verificationCodeMax)

	// send verification code
	err = a.app.Mailer.Send(
		fetchedUser.Email,
		emailVerificationTmpl,
		struct {
			FirstName string
			Code      int
		}{
			FirstName: fetchedUser.FirstName,
			Code:      int(verificationCode),
		},
	)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	// create new verification code
	err = a.app.Repositories.User.CreateVerificationCode(
		fetchedUser.ID, reqVerificationType, int(verificationCode))
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, "verification code sent.")
}

func (a AuthHandler) VerifyContactCallback(ctx echo.Context) error {
	reqData := struct {
		Code int `json:"verification_code" validate:"required"`
	}{}
	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if reqData.Code < verificationCodeMin || reqData.Code > verificationCodeMax {
		return echo.ErrBadRequest
	}

	reqVerificationType := ctx.QueryParam("verification_type")
	if !util.ValidVerificationType(reqVerificationType) {
		return echo.ErrBadRequest
	}

	authenticatedUser := util.ContextGetUser(ctx)
	fetchedUser, err := a.app.Repositories.User.GetByEmail(authenticatedUser.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	// prevent multiple email verification
	if (reqVerificationType == verificationTypeEmail) && fetchedUser.EmailVerified {
		return echo.NewHTTPError(
			http.StatusForbidden,
			ErrAlreadyVerified)
	}
	if (reqVerificationType == verificationTypePhone) && fetchedUser.PhoneVerified {
		return echo.NewHTTPError(
			http.StatusForbidden,
			ErrAlreadyVerified)
	}

	activeCode, err := a.app.Repositories.User.GetVerificationCode(
		fetchedUser.ID, reqVerificationType)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(
				http.StatusForbidden,
				ErrInvalidVerificationCode)
		}
		return util.ErrInternalServer(ctx, err)
	} else {
		switch time.Now().Compare(activeCode.ExpiresAt) {
		case -1:
			// code is still active
			if activeCode.Code != reqData.Code {
				return echo.NewHTTPError(
					http.StatusForbidden,
					ErrInvalidVerificationCode)
			}
		case 0, 1:
			// expired code.
			err = a.app.Repositories.User.DeleteVerificationCode(
				fetchedUser.ID, reqVerificationType)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}
			return echo.NewHTTPError(
				http.StatusForbidden,
				ErrExpiredVerificationCode)
		}
	}

	err = a.app.Repositories.User.VerifyContact(
		fetchedUser.ID, reqVerificationType, true)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	} else {
		if reqVerificationType == verificationTypeEmail {
			fetchedUser.EmailVerified = true
		} else {
			fetchedUser.PhoneVerified = true
		}
	}

	// generate new access and refresh tokens
	accessToken, refreshToken, expiration, err := util.GenerateTokens(a.app, &fetchedUser)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	err = a.app.Repositories.User.DeleteVerificationCode(
		fetchedUser.ID, reqVerificationType)
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

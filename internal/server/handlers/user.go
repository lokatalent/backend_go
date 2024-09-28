package handlers

import (
	"errors"
	// "fmt"
	"net/http"
	// "strings"

	"github.com/labstack/echo/v4"

	"github.com/lokatalent/backend_go/cmd/api/models/response"
	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

type UserHandler struct {
	app *util.Application
}

func NewUserHandler(app *util.Application) UserHandler {
	return UserHandler{app: app}
}

func (u UserHandler) ListUsers(ctx echo.Context) error {
	reqData := struct {
		Page     int `query:"page" validate:"gte=0"`
		PageSize int `query:"size" validate:"required,gte=1"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if reqData.Page < 1 {
		reqData.Page = models.DefaultPage
	}
	if reqData.PageSize < 1 {
		reqData.PageSize = models.DefaultPageLimit
	}

	filter := models.Filter{
		Page:  reqData.Page,
		Limit: reqData.PageSize,
	}

	users, err := u.app.Repositories.User.GetAllUsers(filter)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	// restrict full user detail to admin
	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := u.app.Repositories.User.GetByEmail(authenticatedUser.Email)
	if err == nil && util.IsAdmin(authUser.Role) {
		resp := []response.UserResponse{}
		for _, user := range users {
			resp = append(resp, response.UserResponseFromModel(&user))
		}
		return ctx.JSON(http.StatusOK, resp)
	}
	resp := []response.PublicUserResponse{}
	for _, user := range users {
		resp = append(resp, response.PublicUserResponseFromModel(&user))
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (u UserHandler) Search(ctx echo.Context) error {
	reqData := struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNum    string `json:"phone_num"`
		Role        string `json:"role"`
		ServiceRole string `json:"service_role"`
		Page        int    `query:"page" validate:"gte=0"`
		PageSize    int    `query:"size" validate:"required,gte=1"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if reqData.Page < 1 {
		reqData.Page = models.DefaultPage
	}
	if reqData.PageSize < 1 {
		reqData.PageSize = models.DefaultPageLimit
	}

	filter := models.Filter{
		FirstName:   reqData.FirstName,
		LastName:    reqData.LastName,
		PhoneNum:    reqData.PhoneNum,
		Role:        reqData.Role,
		ServiceRole: reqData.ServiceRole,
		Page:        reqData.Page,
		Limit:       reqData.PageSize,
	}

	users, err := u.app.Repositories.User.Search(filter)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	// restrict full user detail to admin
	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := u.app.Repositories.User.GetByEmail(authenticatedUser.Email)
	if err == nil && util.IsAdmin(authUser.Role) {
		resp := []response.UserResponse{}
		for _, user := range users {
			resp = append(resp, response.UserResponseFromModel(&user))
		}
		return ctx.JSON(http.StatusOK, resp)
	}
	resp := []response.PublicUserResponse{}
	for _, user := range users {
		resp = append(resp, response.PublicUserResponseFromModel(&user))
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (u UserHandler) GetOwnProfile(ctx echo.Context) error {
	authenticatedUser := util.ContextGetUser(ctx)

	user, err := u.app.Repositories.User.GetByEmail(authenticatedUser.Email)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	resp := response.UserResponseFromModel(&user)
	return ctx.JSON(http.StatusOK, resp)
}

func (u UserHandler) GetProfile(ctx echo.Context) error {
	id := ctx.Param("id")

	if !util.IsValidUUID(id) {
		return echo.ErrBadRequest
	}

	user, err := u.app.Repositories.User.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	// restrict full detail to admin
	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := u.app.Repositories.User.GetByEmail(authenticatedUser.Email)
	if err == nil && util.IsAdmin(authUser.Role) {
		resp := response.UserResponseFromModel(&user)
		return ctx.JSON(http.StatusOK, resp)
	}
	resp := response.PublicUserResponseFromModel(&user)
	return ctx.JSON(http.StatusOK, resp)
}

func (u UserHandler) ChangeRole(ctx echo.Context) error {
	id := ctx.Param("id")
	newRole := ctx.QueryParam("role")

	if !util.IsValidUUID(id) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid user id!")
	}
	if (newRole != models.USER_REGULAR) && !util.IsAdmin(newRole) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid role. expected one of regular|admin|admin_super")
	}

	authenticatedUser := util.ContextGetUser(ctx)
	authUser, err := u.app.Repositories.User.GetByEmail(authenticatedUser.Email)
	if err != nil || (authUser.Role != models.USER_ADMIN_SUPER) {
		if err == nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"only super admin can change role")
		}
		return util.ErrInternalServer(ctx, err)
	}

	err = u.app.Repositories.User.ChangeRole(id, newRole)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, echo.Map{})
}

func (u UserHandler) ChangeServiceRole(ctx echo.Context) error {
	// id := ctx.Param("id")
	newRole := ctx.QueryParam("role")

	/*
		if !util.IsValidUUID(id) {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				"invalid user id!")
		}
	*/
	if !util.IsValidServiceRole(newRole) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid role. expected one of service_provider|service_requester|service_both")
	}

	authenticatedUser := util.ContextGetUser(ctx)
	/*
		authUser, err := u.app.Repositories.User.GetByEmail(authenticatedUser.Email)
		if err != nil || !util.IsAdmin(authUser.Role) {
			if err == nil {
				return echo.NewHTTPError(
					http.StatusUnauthorized,
					"only admin can change role")
			}
			return util.ErrInternalServer(ctx, err)
		}
	*/

	err := u.app.Repositories.User.ChangeServiceRole(authenticatedUser.ID, newRole)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, echo.Map{})
}

func (u UserHandler) UpdateProfile(ctx echo.Context) error {
	reqData := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		PhoneNum  string `json:"phone_num"`
		Bio       string `json:"bio"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}

	user := util.ContextGetUser(ctx)
	existingUser, err := u.app.Repositories.User.GetByEmail(user.Email)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	user.IsVerified = existingUser.IsVerified
	user.PhoneNum = existingUser.PhoneNum
	user.FirstName = existingUser.FirstName
	user.LastName = existingUser.LastName
	user.Bio = existingUser.Bio

	if len(existingUser.Bio) == 0 || len(reqData.Bio) > 1 {
		user.Bio = reqData.Bio
	}

	if len(reqData.FirstName) > 1 {
		user.FirstName = reqData.FirstName
	}

	if len(reqData.LastName) > 1 {
		user.LastName = reqData.LastName
	}

	if len(reqData.PhoneNum) > 1 {
		if !util.IsValidPhoneNumber(reqData.PhoneNum) {
			return echo.ErrBadRequest
		}
		// remove verification if user has been verified before, this is done
		// to ensure that new phone number is verified again.
		if existingUser.IsVerified {
			user.IsVerified = false
		}
		user.PhoneNum = reqData.PhoneNum
	}

	updatedUser, err := u.app.Repositories.User.Update(user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateDetails) {
			return echo.NewHTTPError(
				http.StatusConflict, repository.ErrDuplicateDetails)
		}
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response.UserResponseFromModel(&updatedUser))
}

func (u UserHandler) UpdateProfileImage(ctx echo.Context) error {
	authenticatedUser := util.ContextGetUser(ctx)

	// validate the existence of a file
	file, fileHeader, err := ctx.Request().FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return echo.NewHTTPError(
				http.StatusBadRequest, "missing image file")
		}
		return util.ErrInternalServer(ctx, err)
	}
	defer file.Close()

	// verify the file is an accepted image type
	contentType, err := util.ValidateContentType(
		fileHeader.Header,
		util.ContentTypeJPEG,
		util.ContentTypePNG,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// upload image file to storage
	imageUrl, err := u.app.Repositories.Storage.UploadFile(file, authenticatedUser.AvatarPath(), contentType)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	// update image url in database
	err = u.app.Repositories.User.UpdateImage(authenticatedUser.ID, imageUrl)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, echo.Map{"url": imageUrl})
}

func (u UserHandler) DeleteProfileImage(ctx echo.Context) error {
	authenticatedUser := util.ContextGetUser(ctx)

	// delete image in storage
	err := u.app.Repositories.Storage.DeleteFile(authenticatedUser.AvatarPath())
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	// reset the corresponding image URL in the database
	err = u.app.Repositories.User.UpdateImage(authenticatedUser.ID, "")
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, "image successfully deleted")
}

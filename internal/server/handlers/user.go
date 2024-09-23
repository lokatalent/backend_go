package handlers

import (
	"errors"
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
	req := struct {
		Page     int `query:"page" validate:"required,gte=0"`
		PageSize int `query:"size" validate:"required,gte=0"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 {
		req.PageSize = models.DefaultPageLimit
	}

	filter := models.Filter{
		Page:  req.Page,
		Limit: req.PageSize,
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
	req := struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNum    string `json:"phone_num"`
		Role        string `json:"role"`
		ServiceRole string `json:"service_role"`
		Page        int    `query:"page" validate:"required,gte=0"`
		PageSize    int    `query:"size" validate:"required,gte=0"`
	}{}

	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if req.Page < 1 {
		req.Page = models.DefaultPage
	}
	if req.PageSize < 1 {
		req.PageSize = models.DefaultPageLimit
	}

	filter := models.Filter{
		Page:  req.Page,
		Limit: req.PageSize,
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
	if err != nil || authUser.Role != models.USER_ADMIN_SUPER {
		if err == nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"only admin can change role")
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
	if err != nil || !util.IsAdmin(authUser.Role) {
		if err == nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"only admin can change role")
		}
		return util.ErrInternalServer(ctx, err)
	}

	err = u.app.Repositories.User.ChangeRole(id, newRole)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, echo.Map{})
}

func (u UserHandler) UpdateProfile(ctx echo.Context) error {
	reqBody := struct {
		Bio      string `json:"bio"`
		PhoneNum string `json:"phone_num" validate:"len=14"`
	}{}

	if err := ctx.Bind(&reqBody); err != nil {
		return echo.ErrBadRequest
	}
	if !util.IsValidPhoneNumber(reqBody.PhoneNum) {
		return echo.ErrBadRequest
	}

	user := util.ContextGetUser(ctx)
	user.Bio = reqBody.Bio
	user.PhoneNum = reqBody.PhoneNum

	updatedUser, err := u.app.Repositories.User.Update(user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateDetails) {
			return echo.NewHTTPError(
				http.StatusConflict, "conflicting details.")
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

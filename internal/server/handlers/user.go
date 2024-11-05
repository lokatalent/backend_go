package handlers

import (
	"context"
	"errors"
	// "fmt"
	"net/http"
	// "strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"

	"github.com/lokatalent/backend_go/cmd/api/models/response"
	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

const (
	MAX_CONCURRENT_UPLOAD = 4
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
	authUser, err := u.app.Repositories.User.GetByID(authenticatedUser.ID)
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
		Email       string `json:"email"`
		Gender      string `json:"gender"`
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
		Email:       reqData.Email,
		Gender:      reqData.Gender,
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
	authUser, err := u.app.Repositories.User.GetByID(authenticatedUser.ID)
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

	user, err := u.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	resp := response.UserResponseFromModel(&user)
	return ctx.JSON(http.StatusOK, resp)
}

func (u UserHandler) GetOwnEducationProfile(ctx echo.Context) error {
	authenticatedUser := util.ContextGetUser(ctx)

	eduInfo, err := u.app.Repositories.User.GetEducationInfo(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, eduInfo)
}

func (u UserHandler) GetOwnBankProfile(ctx echo.Context) error {
	authenticatedUser := util.ContextGetUser(ctx)

	bankInfo, err := u.app.Repositories.User.GetBankInfo(authenticatedUser.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, bankInfo)
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
	authUser, err := u.app.Repositories.User.GetByID(authenticatedUser.ID)
	if err == nil && util.IsAdmin(authUser.Role) {
		resp := response.UserResponseFromModel(&user)
		return ctx.JSON(http.StatusOK, resp)
	}
	resp := response.PublicUserResponseFromModel(&user)
	return ctx.JSON(http.StatusOK, resp)
}

func (u UserHandler) GetEducationProfile(ctx echo.Context) error {
	id := ctx.Param("id")

	if !util.IsValidUUID(id) {
		return echo.ErrBadRequest
	}

	eduInfo, err := u.app.Repositories.User.GetEducationInfo(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, eduInfo)
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
	authUser, err := u.app.Repositories.User.GetByID(authenticatedUser.ID)
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

func (u UserHandler) UpdateProfilePersonal(ctx echo.Context) error {
	reqData := struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		PhoneNum    string `json:"phone_num"`
		Bio         string `json:"bio"`
		DateOfBirth string `json:"date_of_birth"`
		Address     string `json:"address"`
		Gender      string `json:"gender"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}

	if !util.ValidPlaceAddress(reqData.Address) {
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidPlaceAddress)
	}

	var parsedDOB time.Time
	var err error
	if reqData.DateOfBirth != "" {
		parsedDOB, err = util.ParseDate(reqData.DateOfBirth)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	user := util.ContextGetUser(ctx)
	existingUser, err := u.app.Repositories.User.GetByID(user.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	user.IsVerified = existingUser.IsVerified
	user.PhoneNum = existingUser.PhoneNum
	user.FirstName = existingUser.FirstName
	user.LastName = existingUser.LastName
	user.Bio = existingUser.Bio
	user.Gender = existingUser.Gender
	user.DateOfBirth = existingUser.DateOfBirth

	if len(existingUser.Bio) == 0 || len(reqData.Bio) > 1 {
		user.Bio = reqData.Bio
	}

	if len(existingUser.Gender) == 0 || len(reqData.Gender) > 0 {
		user.Gender = reqData.Gender
	}

	if existingUser.DateOfBirth.IsZero() || !parsedDOB.IsZero() {
		user.DateOfBirth = parsedDOB
	}

	if len(existingUser.Address) == 0 || len(reqData.Address) > 0 {
		user.Address = reqData.Address
	}

	if len(reqData.FirstName) > 1 {
		user.FirstName = reqData.FirstName
	}

	if len(reqData.LastName) > 1 {
		user.LastName = reqData.LastName
	}

	if len(reqData.PhoneNum) > 1 {
		if !util.IsValidPhoneNumber(reqData.PhoneNum) {
			return ErrInvalidPhone
		}
		// remove verification if user has been verified before, this is done
		// to ensure that new phone number is verified again.
		if existingUser.IsVerified {
			user.IsVerified = false
		}
		user.PhoneNum = reqData.PhoneNum
	}

	err = u.app.Repositories.User.Update(&user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateDetails) {
			return echo.NewHTTPError(
				http.StatusConflict, repository.ErrDuplicateDetails)
		}
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response.UserResponseFromModel(&user))
}

func (u UserHandler) UpdateProfileEducation(ctx echo.Context) error {
	isNewInfo := false

	reqData := struct {
		Institute  string `json:"institute"`
		Degree     string `json:"degree"`
		Discipline string `json:"discipline"`
		Start      string `json:"start"`
		Finish     string `json:"finish"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}

	var parsedStartDate, parsedEndDate time.Time
	var err error

	if reqData.Start != "" {
		parsedStartDate, err = util.ParseDate(reqData.Start)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	if reqData.Finish != "" {
		parsedEndDate, err = util.ParseDate(reqData.Finish)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	if parsedStartDate.IsZero() || !parsedStartDate.Before(parsedEndDate) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid start and finish date.",
		)
	}

	user := util.ContextGetUser(ctx)
	existingUser, err := u.app.Repositories.User.GetByID(user.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	existingEduInfo, err := u.app.Repositories.User.GetEducationInfo(existingUser.ID)
	if err != nil {
		if !errors.Is(err, repository.ErrRecordNotFound) {
			return util.ErrInternalServer(ctx, err)
		}
		isNewInfo = true
	}

	if len(existingEduInfo.Institute) == 0 || len(reqData.Institute) > 1 {
		existingEduInfo.Institute = reqData.Institute
	}

	if len(existingEduInfo.Degree) == 0 || len(reqData.Degree) > 1 {
		existingEduInfo.Degree = reqData.Degree
	}

	if len(existingEduInfo.Discipline) == 0 || len(reqData.Discipline) > 1 {
		existingEduInfo.Discipline = reqData.Discipline
	}

	if existingEduInfo.Start.IsZero() || !parsedStartDate.IsZero() {
		existingEduInfo.Start = parsedStartDate
	}

	if existingEduInfo.Finish.IsZero() || !parsedEndDate.IsZero() {
		existingEduInfo.Finish = parsedEndDate
	}

	if isNewInfo {
		existingEduInfo.UserID = existingUser.ID
		err = u.app.Repositories.User.CreateEducationInfo(&existingEduInfo)
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
	} else {
		err = u.app.Repositories.User.UpdateEducationInfo(&existingEduInfo)
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, existingEduInfo)
}

func (u UserHandler) UpdateProfileBank(ctx echo.Context) error {
	isNewInfo := false

	reqData := struct {
		BankName    string `json:"bank_name"`
		AccountName string `json:"account_name"`
		AccountNum  string `json:"account_num" validate:"max=10"`
	}{}

	if err := ctx.Bind(&reqData); err != nil {
		return echo.ErrBadRequest
	}
	if err := ctx.Validate(&reqData); err != nil {
		return echo.ErrBadRequest
	}

	user := util.ContextGetUser(ctx)
	existingUser, err := u.app.Repositories.User.GetByID(user.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	existingBankInfo, err := u.app.Repositories.User.GetBankInfo(existingUser.ID)
	if err != nil {
		if !errors.Is(err, repository.ErrRecordNotFound) {
			return util.ErrInternalServer(ctx, err)
		}
		isNewInfo = true
	}

	if len(existingBankInfo.BankName) == 0 || len(reqData.BankName) > 1 {
		existingBankInfo.BankName = reqData.BankName
	}

	if len(existingBankInfo.AccountName) == 0 || len(reqData.AccountName) > 1 {
		existingBankInfo.AccountName = reqData.AccountName
	}

	if len(existingBankInfo.AccountNum) == 0 || len(reqData.AccountNum) > 1 {
		existingBankInfo.AccountNum = reqData.AccountNum
	}

	if isNewInfo {
		existingBankInfo.UserID = existingUser.ID
		err = u.app.Repositories.User.CreateBankInfo(&existingBankInfo)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicateBankDetails) {
				return echo.NewHTTPError(
					http.StatusForbidden,
					repository.ErrDuplicateBankDetails,
				)
			}
			return util.ErrInternalServer(ctx, err)
		}
	} else {
		err = u.app.Repositories.User.UpdateBankInfo(&existingBankInfo)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicateBankDetails) {
				return echo.NewHTTPError(
					http.StatusForbidden,
					repository.ErrDuplicateBankDetails,
				)
			}
			return util.ErrInternalServer(ctx, err)
		}
	}

	return ctx.JSON(http.StatusOK, existingBankInfo)
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

func (u UserHandler) GetCertifications(ctx echo.Context) error {
	id := ctx.Param("id")

	certs, err := u.app.Repositories.User.GetCertifications(id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"no uploaded certifications",
			)
		}
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, certs)
}

func (u UserHandler) UploadCertifications(ctx echo.Context) error {
	var (
		sem = util.NewSemaphore(MAX_CONCURRENT_UPLOAD)
		mu  sync.Mutex
		g   = new(errgroup.Group)
	)

	uploadCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	certs := []models.UserCertification{}

	authenticatedUser := util.ContextGetUser(ctx)

	form, err := ctx.MultipartForm()
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return echo.NewHTTPError(http.StatusBadRequest, "certificate files are required")
		}
		return util.ErrInternalServer(ctx, err)
	}

	for _, image := range form.File["images"] {
		contentType, err := util.ValidateContentType(
			image.Header,
			util.ContentTypeJPEG,
			util.ContentTypePNG,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		src, err := image.Open()
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
		defer src.Close()

		g.Go(func() error {
			select {
			case <-uploadCtx.Done():
				return nil // context canceled: stop further upload.
			default:
				// proceed to upload file
			}

			sem.Acquire()
			defer sem.Release()

			newCert := models.UserCertification{
				ID:     uuid.NewString(),
				UserID: authenticatedUser.ID,
			}
			imgURL, err := u.app.Repositories.Storage.UploadFile(
				src,
				newCert.CertificationPath(),
				contentType,
			)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}

			newCert.URL = imgURL
			err = u.app.Repositories.User.CreateCertification(&newCert)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}

			mu.Lock()
			defer mu.Unlock()
			certs = append(certs, newCert)

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, certs)
}

func (u UserHandler) DeleteCertification(ctx echo.Context) error {
	certID := ctx.Param("id")
	authenticatedUser := util.ContextGetUser(ctx)

	cert := models.UserCertification{ID: certID, UserID: authenticatedUser.ID}

	err := u.app.Repositories.Storage.DeleteFile(cert.CertificationPath())
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	err = u.app.Repositories.User.DeleteCertification(certID, authenticatedUser.ID)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, "image successfully deleted.")
}

func (u UserHandler) GetServiceImages(ctx echo.Context) error {
	id := ctx.Param("id")
	serviceType := ctx.QueryParam("service_type")
	if !util.IsValidServiceType(serviceType) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrInvalidServiceType,
		)
	}

	imgs, err := u.app.Repositories.User.GetServiceImages(id, serviceType)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"no uploaded images",
			)
		}
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, imgs)
}

func (u UserHandler) UploadServiceImages(ctx echo.Context) error {
	serviceType := ctx.QueryParam("service_type")
	if !util.IsValidServiceType(serviceType) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrInvalidServiceType,
		)
	}

	var (
		sem = util.NewSemaphore(MAX_CONCURRENT_UPLOAD)
		mu  sync.Mutex
		g   = new(errgroup.Group)
	)

	uploadCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	imgs := []models.ServiceImage{}

	authenticatedUser := util.ContextGetUser(ctx)

	form, err := ctx.MultipartForm()
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return echo.NewHTTPError(
				http.StatusBadRequest, "certificate files are required")
		}
		return util.ErrInternalServer(ctx, err)
	}

	for _, image := range form.File["images"] {
		contentType, err := util.ValidateContentType(
			image.Header,
			util.ContentTypeJPEG,
			util.ContentTypePNG,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		src, err := image.Open()
		if err != nil {
			return util.ErrInternalServer(ctx, err)
		}
		defer src.Close()

		g.Go(func() error {
			select {
			case <-uploadCtx.Done():
				return nil // context canceled: stop further upload.
			default:
				// proceed to upload file
			}

			sem.Acquire()
			defer sem.Release()

			newImg := models.ServiceImage{
				ID:          uuid.NewString(),
				UserID:      authenticatedUser.ID,
				ServiceType: serviceType,
			}
			imgURL, err := u.app.Repositories.Storage.UploadFile(
				src,
				newImg.ServiceImagePath(),
				contentType,
			)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}

			newImg.URL = imgURL
			err = u.app.Repositories.User.CreateServiceImage(&newImg)
			if err != nil {
				return util.ErrInternalServer(ctx, err)
			}

			mu.Lock()
			defer mu.Unlock()
			imgs = append(imgs, newImg)

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, imgs)
}

func (u UserHandler) DeleteServiceImage(ctx echo.Context) error {
	imgID := ctx.Param("id")
	serviceType := ctx.QueryParam("service_type")
	if !util.IsValidServiceType(serviceType) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			ErrInvalidServiceType,
		)
	}

	authenticatedUser := util.ContextGetUser(ctx)
	img := models.ServiceImage{
		ID:          imgID,
		UserID:      authenticatedUser.ID,
		ServiceType: serviceType,
	}

	err := u.app.Repositories.Storage.DeleteFile(img.ServiceImagePath())
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	err = u.app.Repositories.User.DeleteServiceImage(imgID, authenticatedUser.ID, serviceType)
	if err != nil {
		return util.ErrInternalServer(ctx, err)
	}

	return ctx.JSON(http.StatusOK, "image successfully deleted.")
}

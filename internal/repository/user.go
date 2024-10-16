package repository

import "github.com/lokatalent/backend_go/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetByPhone(phone string) (models.User, error)
	GetAllUsers(filter models.Filter) ([]models.User, error)
	Update(user *models.User) error
	UpdateImage(id string, imageURL string) error
	ChangeRole(id, role string) error
	ChangeServiceRole(id, role string) error
	Search(filter models.Filter) ([]models.User, error)

	Verify(id string, status bool) error
	VerifyContact(id, verificationType string, status bool) error
	CreateVerificationCode(verCode *models.UserVerificationCode) error
	DeleteVerificationCode(id, verificationType string) error
	GetVerificationCode(id, verificationType string) (models.UserVerificationCode, error)

	CreateBankInfo(bankInfo *models.UserBankInfo) error
	GetBankInfo(userID string) (models.UserBankInfo, error)
	UpdateBankInfo(bankInfo *models.UserBankInfo) error

	CreateEducationInfo(eduInfo *models.UserEducationInfo) error
	GetEducationInfo(userID string) (models.UserEducationInfo, error)
	UpdateEducationInfo(eduInfo *models.UserEducationInfo) error

	CreateCertification(cert *models.UserCertification) error
	GetCertifications(userID string) ([]models.UserCertification, error)
	DeleteCertification(id, userID string) error

	CreateService(service *models.UserService) error
	GetService(userID, serviceType string) (models.UserService, error)
	GetAllServices(userID string) ([]models.UserService, error)
	UpdateService(service *models.UserService) error
	DeleteService(userID, serviceType string) error

	CreateServiceImage(seviceImg *models.ServiceImage) error
	GetServiceImages(userID, serviceType string) ([]models.ServiceImage, error)
	DeleteServiceImage(id, userID, serviceType string) error
}

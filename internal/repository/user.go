package repository

import "github.com/lokatalent/backend_go/internal/models"

type UserRepository interface {
	Create(user models.User) (models.User, error)
	GetByID(id string) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetAllUsers(filter models.Filter) ([]models.User, error)
	Update(user models.User) (models.User, error)
	UpdateImage(id string, imageURL string) error
	ChangeRole(id, role string) error
	ChangeServiceRole(id, role string) error
	Search(filter models.Filter) ([]models.User, error)

	Verify(id string, status bool) error
	VerifyContact(id, verificationType string, status bool) error
	CreateVerificationCode(id, verificationType string, code int) error
	DeleteVerificationCode(id, verificationType string) error
	GetVerificationCode(id, verificationType string) (models.UserVerificationCode, error)
}

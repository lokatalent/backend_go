package response

import (
	"time"

	"github.com/lokatalent/backend_go/internal/models"
)

type UserResponse struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	PhoneNum      string    `json:"phone_num"`
	Gender        string    `json:"gender"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	Bio           string    `json:"bio"`
	Address       string    `json:"address"`
	Avatar        string    `json:"avatar"`
	Role          string    `json:"role"`
	ServiceRole   string    `json:"service_role"`
	IsVerified    bool      `json:"is_verified"`
	EmailVerified bool      `json:"email_verified"`
	PhoneVerified bool      `json:"phone_verified"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PublicUserResponse struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Gender      string    `json:"gender"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Bio         string    `json:"bio"`
	Avatar      string    `json:"avatar"`
	Role        string    `json:"role"`
	ServiceRole string    `json:"service_role"`
	IsVerified  bool      `json:"is_verified"`

	CreatedAt time.Time `json:"created_at"`
}

func UserResponseFromModel(user *models.User) UserResponse {
	response := UserResponse{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		PhoneNum:      user.PhoneNum,
		Gender:        user.Gender,
		DateOfBirth:   user.DateOfBirth,
		Bio:           user.Bio,
		Address:       user.Address,
		Avatar:        user.Avatar,
		Role:          user.Role,
		ServiceRole:   user.ServiceRole,
		IsVerified:    user.IsVerified,
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		CreatedAt:     user.CreatedAt.UTC(),
		UpdatedAt:     user.UpdatedAt.UTC(),
	}

	return response
}

func PublicUserResponseFromModel(user *models.User) PublicUserResponse {
	response := PublicUserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Gender:      user.Gender,
		DateOfBirth: user.DateOfBirth,
		Bio:         user.Bio,
		Avatar:      user.Avatar,
		Role:        user.Role,
		ServiceRole: user.ServiceRole,
		IsVerified:  user.IsVerified,
		CreatedAt:   user.CreatedAt.UTC(),
	}

	return response
}

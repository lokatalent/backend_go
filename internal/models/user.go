package models

import (
	"path"
	"time"
)

type User struct {
	ID            string
	FirstName     string
	LastName      string
	Email         string
	PhoneNum      string
	Password      string
	Gender        string
	DateOfBirth   time.Time
	Bio           string
	Address       string
	Avatar        string
	Role          string
	ServiceRole   string
	IsVerified    bool
	EmailVerified bool
	PhoneVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserVerificationCode struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Code        int       `json:"code"`
	ContactType string    `json:"contact_type"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type UserBankInfo struct {
	UserID      string    `json:"user_id"`
	BankName    string    `json:"bank_name"`
	AccountName string    `json:"account_name"`
	AccountNum  string    `json:"account_num"`
	BankCode    string    `json:"bank_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserEducationInfo struct {
	UserID     string    `json:"user_id"`
	Institute  string    `json:"institute"`
	Degree     string    `json:"degree"`
	Discipline string    `json:"discipline"`
	Start      time.Time `json:"start"`
	Finish     time.Time `json:"finish"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserCertification struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	URL    string `json:"url"`
}

// AvatarPath returns the relative path to user's avatar.
func (u *User) AvatarPath() string {
	return path.Join("profiles", "avatars", u.ID)
}

func (u *UserCertification) CertificationPath() string {
	return path.Join("profiles", "certifications", u.UserID, u.ID)
}

package models

import (
	"path"
	"time"
)

type User struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	PhoneNum    string
	Password    string
	Bio         string
	Avatar      string
	Role        string
	ServiceRole string
	IsVerified  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// AvatarPath returns the relative path to user's avatar.
func (u *User) AvatarPath() string {
	return path.Join("profiles/avatar", u.ID)
}

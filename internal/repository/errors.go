package repository

import "errors"

var (
	ErrRecordNotFound   = errors.New("No matching record found")
	ErrDuplicateDetails = errors.New("User with email or phone number already exist")
)

package repository

import "errors"

var (
	ErrRecordNotFound       = errors.New("No matching record found")
	ErrDuplicateDetails     = errors.New("User with email or phone number already exist")
	ErrDuplicateService     = errors.New("Service already exist.")
	ErrDuplicateBankDetails = errors.New("Bank account number already exists.")
	ErrInvalidWalletBalance = errors.New("Invalid user wallet balance.")
)

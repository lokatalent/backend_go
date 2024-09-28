package handlers

import "errors"

var (
	ErrInvalidPassword         = errors.New("Invalid password!")
	ErrInvalidToken            = errors.New("Invalid or expired Token")
	ErrInvalidEmail            = errors.New("Invalid email address")
	ErrInvalidPhone            = errors.New("Invalid phone number")
	ErrAlreadyVerified         = errors.New("Already verified")
	ErrVerificationDependency  = errors.New("email and phone number needs to be verified!")
	ErrExpiredVerificationCode = errors.New("verification code has expired.")
	ErrInvalidVerificationCode = errors.New("invalid verification code.")
)
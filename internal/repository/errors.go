package repository

import "errors"

var (
	ErrRecordNotFound   = errors.New("No matching record found")
	ErrDuplicateDetails = errors.New("duplicate details found")
)

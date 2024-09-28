package repository

import (
	"io"
)

type StorageRepository interface {
	UploadFile(file io.Reader, key string, contentType string) (string, error)
	DeleteFile(key string) error
}

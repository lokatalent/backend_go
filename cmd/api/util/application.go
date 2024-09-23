package util

import (
	"github.com/lokatalent/backend_go/internal/repository"
)

type Application struct {
	Config       *Config
	Repositories *repository.Repositories
}

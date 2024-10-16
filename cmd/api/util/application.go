package util

import (
	"github.com/lokatalent/backend_go/internal/mailer"
	"github.com/lokatalent/backend_go/internal/repository"
	"github.com/lokatalent/backend_go/internal/sms"
)

type Application struct {
	Config       *Config
	Repositories *repository.Repositories
	Mailer       *mailer.Mailer
	SMSSender    *sms.SMSSender
}

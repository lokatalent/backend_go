package main

import (
	"database/sql"
	"fmt"
	"log"
	// "os"

	"github.com/lokatalent/backend_go/cmd/api/util"
	"github.com/lokatalent/backend_go/internal/database/postgres"
	"github.com/lokatalent/backend_go/internal/mailer"
	"github.com/lokatalent/backend_go/internal/repository"
	"github.com/lokatalent/backend_go/internal/s3"
	"github.com/lokatalent/backend_go/internal/server/routes"
	"github.com/lokatalent/backend_go/internal/sms"
)

func serveApp(config *util.Config, db *sql.DB) error {
	repos := &repository.Repositories{
		User:    postgres.NewUserImplementation(db),
		Storage: s3.NewStorageInfrastructure(config.AWS.S3Bucket),
	}

	app := util.Application{
		Config:       config,
		Repositories: repos,
		Mailer: mailer.New(
			mailer.Credentials{
				APIKey: config.SendGrid.APIKey,
				Sender: config.SendGrid.Sender,
			},
		),
		SMSSender: sms.New(
			sms.Credentials{
				APIKey:     config.Twilio.APIKey,
				APISecret:  config.Twilio.APISecret,
				AccountSID: config.Twilio.AccountSID,
				Sender:     config.Twilio.Sender,
			},
		),
	}

	engine := routes.Engine(&app)

	if err := engine.Start(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
		return err
	}

	log.Println("server stopped")

	return nil
}

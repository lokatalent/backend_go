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
		User:           postgres.NewUserImplementation(db),
		Storage:        s3.NewStorageInfrastructure(config.AWS.S3Bucket),
		Commission:     postgres.NewCommissionImplementation(db),
		ServicePricing: postgres.NewServicePricingImplementation(db),
		Booking:        postgres.NewBookingImplementation(db),
		Notification:   postgres.NewNotificationImplementation(db),
		Payment:        postgres.NewPaymentImplementation(db),
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

	/*
		switch app.Config.Env {
		case util.ENVIRONMENT_PRODUCTION:
			if err := engine.Start(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
				return err
			}
		default:
			if err := engine.StartTLS(
				fmt.Sprintf(":%d", app.Config.Port),
				app.Config.CertPath,
				app.Config.CertKeyPath,
			); err != nil {
				return err
			}
		}
	*/
	if err := engine.Start(fmt.Sprintf(":%d", app.Config.Port)); err != nil {
		return err
	}

	log.Println("server stopped")

	return nil
}

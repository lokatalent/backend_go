package sms

import (
	"bytes"
	"embed"
	"text/template"
	"time"

	"github.com/twilio/twilio-go"
	twilioAPI "github.com/twilio/twilio-go/rest/api/v2010"
)

//go:embed "templates"
var templateFS embed.FS

const smsDelay = 5 * time.Second

// SMSSender contains a client instance and the sender information for sending sms.
type SMSSender struct {
	client *twilio.RestClient
	sender string
}

type Credentials struct {
	APIKey     string
	APISecret  string
	AccountSID string
	Sender     string
}

// New returns a new SMSSender.
func New(details Credentials) *SMSSender {
	client := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username:   details.APIKey,
			Password:   details.APISecret,
			AccountSid: details.AccountSID,
		},
	)
	return &SMSSender{
		client: client,
		sender: details.Sender,
	}
}

// Send takes the recipient phone number, the sms template file and initiates SMS transfer.
func (s *SMSSender) Send(recipient string, templateFile string, data any) error {
	tmpl, err := template.New("sms").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	shortMsg := &bytes.Buffer{}
	err = tmpl.Execute(shortMsg, data)
	if err != nil {
		return err
	}

	msgParams := &twilioAPI.CreateMessageParams{}
	msgParams.SetTo(recipient)
	msgParams.SetFrom(s.sender)
	msgParams.SetBody(shortMsg.String())

	// attempt sending SMS 3 times
	for range 3 {
		_, err = s.client.Api.CreateMessage(msgParams)
		if err == nil {
			return nil
		}

		time.Sleep(smsDelay)
	}

	return err
}

package mailer

import (
	"bytes"
	"embed"
	// "fmt"
	"html/template"
	// "log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//go:embed "templates"
var templateFS embed.FS

const mailDelay = 5 * time.Second // delay between mail sending trials

// The Mailer struct contains a client instance and the sender information for sending emails.
type Mailer struct {
	client *sendgrid.Client
	sender string
}

type Credentials struct {
	APIKey string
	Sender string
}

// New returns a Mailer pointer with a sendgrid client
func New(details Credentials) *Mailer {
	client := sendgrid.NewSendClient(details.APIKey)
	return &Mailer{
		client: client,
		sender: details.Sender,
	}
}

// Send takes the recipient email address, the name of the email template file, and any dynamic data for the template.
func (m *Mailer) Send(recipient string, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	msg := mail.NewSingleEmail(
		mail.NewEmail(
			m.sender,
			m.sender),
		subject.String(),
		mail.NewEmail(recipient, recipient),
		plainBody.String(),
		htmlBody.String(),
	)

	// attempt sending mail 3 times
	for range 3 {
		_, err = m.client.Send(msg)
		// return if no error occurs
		if err == nil {
			return nil
		}
		// fmt.Println(err)

		// make next attempt on failure with a delay of 5 seconds
		time.Sleep(mailDelay)
	}

	// log.Println(err)
	return err
}

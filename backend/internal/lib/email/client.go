package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/anuragShingare30/go-boilerplate/internal/config"
	"github.com/pkg/errors"
	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog"
)

type Client struct {
	client *resend.Client
	logger *zerolog.Logger
}

// @dev Initialize the email client with Resend API key and logger
func NewClient(cfg *config.Config, logger *zerolog.Logger) *Client {
	return &Client{
		client: resend.NewClient(cfg.Integration.ResendAPIKey),
		logger: logger,
	}
}

// @dev Method which will actually use the Resend Client to actually send emails
func (c *Client) SendEmail(to string, subject string, templateName Template, data map[string]string) error {
	// templates/emails/Welcome.html
	tmplPath := fmt.Sprintf("%s/%s.html", "templates/emails", templateName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse email template %s", templateName)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return errors.Wrapf(err, "failed to execute email template %s", templateName)
	}

	// @dev SendEmailRequest is the request object for the Send call
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", "Boilerplate", "onboarding@resend.dev"),
		To:      []string{to},
		Subject: subject,
		Html:    body.String(),
	}

	_, err = c.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return err
}

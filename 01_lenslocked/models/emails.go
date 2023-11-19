package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	From, To, Subject, Html, Plaintext string
}

type SmtpConfig struct {
	Host, Username, Password string
	Port                     int
}

type EmailService struct {
	DefaultSender string
	dialer        *mail.Dialer
}

func (es *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	es.setFrom(msg, email)
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)
	switch {
	case email.Plaintext != "" && email.Html != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.Html)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.Html != "":
		msg.SetBody("text/html", email.Html)
	}
	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (es *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}

func NewEmailService(config *SmtpConfig) *EmailService {
	return &EmailService{
		dialer: mail.NewDialer(
			config.Host,
			config.Port,
			config.Username,
			config.Password,
		),
	}
}

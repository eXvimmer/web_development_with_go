package models

import (
	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type SmtpConfig struct {
	Host, Username, Password string
	Port                     int
}

type EmailService struct {
	DefaultSender string
	dialer        *mail.Dialer
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

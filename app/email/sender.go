// Package mail provides utilities for constructing and sending email messages.
package mail

import (
	"fmt"
	"log"
	"net/smtp"
)

// SMTPEmailSender implements an email sender using SMTP protocol.
type SMTPEmailSender struct {
	smtpHost string // SMTP server hostname
	smtpPort string // SMTP server port
}

// NewSMTPEmailSender creates a new SMTPEmailSender with given host and port.
//
// Parameters:
//   - smtpHost: SMTP server hostname (e.g., "smtp.example.com")
//   - smtpPort: SMTP server port (e.g., "587")
//
// Returns:
//   - pointer to a configured SMTPEmailSender instance
func NewSMTPEmailSender(smtpHost, smtpPort string) *SMTPEmailSender {
	return &SMTPEmailSender{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

// Send sends an EmailMessage using the configured SMTP server.
//
// Parameters:
//   - email: the EmailMessage to be sent
//
// Returns:
//   - error if sending the email fails; otherwise, nil
//
// Note:
//
//	This implementation uses smtp.SendMail without authentication.
//	Adjustments may be required to support authentication depending on your SMTP server.
func (s *SMTPEmailSender) Send(email EmailMessage) error {
	to := []string{email.To}

	msg := []byte(fmt.Sprintf(
		"Subject: %s\r\nFrom: %s\r\nTo: %s\r\nContent-Type: %s\r\n\r\n%s",
		email.Subject, email.From, email.To, email.ContentType, email.Body,
	))

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	if err := smtp.SendMail(addr, nil, email.From, to, msg); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}

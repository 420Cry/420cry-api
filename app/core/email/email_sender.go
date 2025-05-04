package core

import (
	EmailDomain "cry-api/app/domain/email"
	"log"
)

// SMTPEmailSender represents the structure for sending emails via SMTP
type SMTPEmailSender struct {
	smtpHost string
	smtpPort string
}

// NewSMTPEmailSender creates a new SMTPEmailSender instance
func NewSMTPEmailSender(smtpHost, smtpPort string) *SMTPEmailSender {
	return &SMTPEmailSender{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

// Send sends the email
func (s *SMTPEmailSender) Send(email EmailDomain.EmailMessage) error {
	// Set logging to include timestamp and source file for better debugging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Log email details
	log.Printf("Sending email to: %s\n", email.To)
	log.Printf("Sending email from: %s\n", email.From)
	log.Printf("Subject: %s\n", email.Subject)
	log.Printf("Body: %s\n", email.Body)

	// Simulate email sending (actual logic to send the email can go here)
	log.Printf("Email sent to: %s successfully", email.To)

	// Return success
	return nil
}

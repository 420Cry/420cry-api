package emailcore

import (
	config "cry-api/app/config"
	EmailDomain "cry-api/app/domain/email"
	"log"
	"net/smtp"
	"strings"
)

// SmtpEmailSender implements EmailSender interface for sending emails via SMTP
type SMTPEmailSender struct {
	smtpHost string
	smtpPort string
	from     string
}

// NewSMTPEmailSender creates a new instance of SmtpEmailSender
func NewSMTPEmailSender() *SMTPEmailSender {
	cfg := config.Get()
	return &SMTPEmailSender{
		smtpHost: cfg.SMTPConfig.Host,
		smtpPort: cfg.SMTPConfig.Port,
		from:     cfg.NoReplyEmail,
	}
}

// Send sends the email via SMTP
func (sender *SMTPEmailSender) Send(email EmailDomain.Email) error {
	// Prepare SMTP authentication
	auth := smtp.PlainAuth("", sender.from, "", sender.smtpHost)

	// Define the email headers and body
	headers := []string{
		"From: " + sender.from,
		"To: " + email.To,
		"Subject: " + email.Subject,
		"Content-Type: " + email.ContentType + "; charset=UTF-8",
	}

	// Combine the headers and body
	msg := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + email.Body)

	// Send the email
	err := smtp.SendMail(sender.smtpHost+":"+sender.smtpPort, auth, sender.from, []string{email.To}, msg)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", email.To)
	return nil
}

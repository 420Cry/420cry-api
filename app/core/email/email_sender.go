package core

import (
	EmailDomain "cry-api/app/domain/email"
	"cry-api/app/utils"
	"fmt"
	"log"
	"net/smtp"
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

// Send sends the email using the SMTP protocol
func (s *SMTPEmailSender) Send(email EmailDomain.EmailMessage) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	to := []string{email.To}
	msg := []byte(fmt.Sprintf("Subject: %s\r\nFrom: %s\r\nTo: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		utils.SanitizeInput(email.Subject), utils.SanitizeInput(email.From), utils.SanitizeInput(email.To), utils.SanitizeInput(email.Body)))

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, nil, email.From, to, msg)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent to: %s successfully", email.To)
	return nil
}

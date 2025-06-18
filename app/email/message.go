// Package mail provides structures and functions for composing email messages.
package mail

// EmailMessage represents the structure of an email with basic fields.
type EmailMessage struct {
	To          string // Recipient email address
	From        string // Sender email address
	Subject     string // Subject of the email
	Body        string // HTML body content of the email
	ContentType string // MIME content type (e.g., "text/html")
}

// NewEmailMessage creates and returns a new EmailMessage with default content type as "text/html".
func NewEmailMessage(to, from, subject, body string) EmailMessage {
	return EmailMessage{
		To:          to,
		From:        from,
		Subject:     subject,
		Body:        body,
		ContentType: "text/html", // Default to HTML content
	}
}

package domain

// EmailMessage represents an email message with its components.
type EmailMessage struct {
	To          string
	From        string
	Subject     string
	Body        string
	ContentType string
}

// NewEmailMessage creates a new EmailMessage instance.
func NewEmailMessage(to, from, subject, body string) EmailMessage {
	return EmailMessage{
		To:          to,
		From:        from,
		Subject:     subject,
		Body:        body,
		ContentType: "text/html",
	}
}

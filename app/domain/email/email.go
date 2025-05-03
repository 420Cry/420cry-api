package emaildomain

type Email struct {
	To          string
	Subject     string
	Body        string
	ContentType string
}

func NewEmail(to, subject, body string) Email {
	return Email{
		To:          to,
		Subject:     subject,
		Body:        body,
		ContentType: "text/html",
	}
}

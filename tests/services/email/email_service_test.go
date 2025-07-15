package tests

import (
	"errors"
	"testing"

	Email "cry-api/app/email"
	Services "cry-api/app/services/email"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getTestResetPasswordLink() string {
	return "https://example.com/reset"
}

func TestSendVerifyAccountEmail_Success(t *testing.T) {
	mockSender := new(testmocks.MockEmailSender)
	mockCreator := new(testmocks.MockEmailCreator)

	service := Services.NewEmailService(mockSender, mockCreator)

	to := "user@example.com"
	from := "no-reply@example.com"
	userName := "testuser"
	verificationLink := "https://example.com/verify"
	verificationTokens := "token123"

	// Expect sanitization call (optional, if utils.SanitizeInput is mocked)
	// Or you can trust it works and test only the email service logic here.

	expectedEmail := Email.EmailMessage{
		To:      to,
		From:    from,
		Subject: "Verify Your Account",
		Body:    "<html>Verification email body</html>",
	}

	mockCreator.
		On("CreateVerifyAccountEmail", to, from, userName, verificationLink, verificationTokens).
		Return(expectedEmail, nil).
		Once()

	mockSender.
		On("Send", expectedEmail).
		Return(nil).
		Once()

	err := service.SendVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens)
	assert.NoError(t, err)

	mockCreator.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

func TestSendVerifyAccountEmail_CreateEmailError(t *testing.T) {
	mockSender := new(testmocks.MockEmailSender)
	mockCreator := new(testmocks.MockEmailCreator)

	service := Services.NewEmailService(mockSender, mockCreator)

	mockCreator.
		On("CreateVerifyAccountEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(Email.EmailMessage{}, errors.New("template error")).
		Once()

	err := service.SendVerifyAccountEmail("to", "from", "user", "link", "token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template error")

	mockCreator.AssertExpectations(t)
}

func TestSendVerifyAccountEmail_SendEmailError(t *testing.T) {
	mockSender := new(testmocks.MockEmailSender)
	mockCreator := new(testmocks.MockEmailCreator)

	service := Services.NewEmailService(mockSender, mockCreator)

	expectedEmail := Email.EmailMessage{}

	mockCreator.
		On("CreateVerifyAccountEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expectedEmail, nil).
		Once()

	mockSender.
		On("Send", expectedEmail).
		Return(errors.New("send error")).
		Once()

	err := service.SendVerifyAccountEmail("to", "from", "user", "link", "token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "send error")

	mockCreator.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

func TestSendResetPasswordEmail_Success(t *testing.T) {
	mockSender := new(testmocks.MockEmailSender)
	mockCreator := new(testmocks.MockEmailCreator)

	service := Services.NewEmailService(mockSender, mockCreator)

	to := "user@example.com"
	from := "no-reply@example.com"
	userName := "testuser"
	resetPasswordLink := getTestResetPasswordLink()
	apiURL := "https://api.example.com"

	expectedEmail := Email.EmailMessage{
		To:      to,
		From:    from,
		Subject: "Reset Password Request",
		Body:    "<html>Reset password body</html>",
	}

	mockCreator.
		On("CreateResetPasswordRequestEmail", to, from, userName, resetPasswordLink, apiURL).
		Return(expectedEmail, nil).
		Once()

	mockSender.
		On("Send", expectedEmail).
		Return(nil).
		Once()

	err := service.SendResetPasswordEmail(to, from, userName, resetPasswordLink, apiURL)
	assert.NoError(t, err)

	mockCreator.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

func TestSendResetPasswordEmail_CreateEmailError(t *testing.T) {
	mockSender := new(testmocks.MockEmailSender)
	mockCreator := new(testmocks.MockEmailCreator)

	service := Services.NewEmailService(mockSender, mockCreator)

	mockCreator.
		On("CreateResetPasswordRequestEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(Email.EmailMessage{}, errors.New("template error")).
		Once()

	mockSender.
		On("Send", mock.Anything).
		Return(nil).
		Once()

	// According to your service code, the error from CreateResetPasswordRequestEmail
	// is logged but the function continues and tries to send the email anyway.

	err := service.SendResetPasswordEmail("to", "from", "user", "link", "apiURL")
	assert.NoError(t, err)

	mockCreator.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

func TestSendResetPasswordEmail_SendEmailError(t *testing.T) {
	mockSender := new(testmocks.MockEmailSender)
	mockCreator := new(testmocks.MockEmailCreator)

	service := Services.NewEmailService(mockSender, mockCreator)

	expectedEmail := Email.EmailMessage{}

	mockCreator.
		On("CreateResetPasswordRequestEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(expectedEmail, nil).
		Once()

	mockSender.
		On("Send", expectedEmail).
		Return(errors.New("send error")).
		Once()

	err := service.SendResetPasswordEmail("to", "from", "user", "link", "apiURL")
	assert.NoError(t, err) // your implementation logs error but returns nil

	mockCreator.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

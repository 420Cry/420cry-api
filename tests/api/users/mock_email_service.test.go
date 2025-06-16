// Package user_routes_test provides tests for user routes.
package user_routes_test

import (
	"github.com/stretchr/testify/mock"
)

// MockEmailService is a mock of EmailService for testing
type MockEmailService struct {
	mock.Mock
}

// SendVerifyAccountEmail mocks sending a verification email to a user with the provided recipient, sender, username, verification link, and token.
func (m *MockEmailService) SendVerifyAccountEmail(to, from, username, link, token string) error {
	args := m.Called(to, from, username, link, token)
	return args.Error(0)
}

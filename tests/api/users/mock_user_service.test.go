// Package user_routes_test provides tests for user routes.
package user_routes_test

import (
	UserDomain "cry-api/app/domain/users"

	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock UserService for testing
type MockUserService struct {
	mock.Mock
}

// CreateUser is a mock implementation of the UserService's CreateUser method.
// It simulates user creation by accepting fullname, username, email, and password,
// and returns a mocked User, a string (such as a token or message), and an error.
func (m *MockUserService) CreateUser(fullname, username, email, password string) (*UserDomain.User, string, error) {
	args := m.Called(fullname, username, email, password)
	return args.Get(0).(*UserDomain.User), args.String(1), args.Error(2)
}

// CheckUserByBothTokens is a mock method that simulates checking a user by both an authentication token and a verification token.
// It returns a pointer to a UserDomain.User and an error, mimicking the behavior of the actual service method for testing purposes.
func (m *MockUserService) CheckUserByBothTokens(token string, verificationToken string) (*UserDomain.User, error) {
	args := m.Called(token, verificationToken)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

// CheckEmailVerificationToken is a mock method for testing CheckEmailVerificationToken.
func (m *MockUserService) CheckEmailVerificationToken(token string) (*UserDomain.User, error) {
	args := m.Called(token)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

// CheckAccountVerificationToken is a mock method for testing CheckAccountVerificationToken.
func (m *MockUserService) CheckAccountVerificationToken(token string) (*UserDomain.User, error) {
	args := m.Called(token)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

// AuthenticateUser is a mock method for testing AuthenticateUser.
func (m *MockUserService) AuthenticateUser(username string, password string) (*UserDomain.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

// VerifyUserWithTokens mocks the verification of a user using the provided token and verificationToken.
// It returns a pointer to a UserDomain.User and an error.
func (m *MockUserService) VerifyUserWithTokens(token string, verificationToken string) (*UserDomain.User, error) {
	args := m.Called(token, verificationToken)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

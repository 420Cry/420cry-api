// Package tests provides tests for user routes.
package tests

import (
	UserDomain "cry-api/app/domain/users"

	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock UserService for testing
type MockUserService struct {
	mock.Mock
}

// CreateUser is a mock implementation of the UserService's CreateUser method.
func (m *MockUserService) CreateUser(fullname, username, email, password string) (*UserDomain.User, string, error) {
	args := m.Called(fullname, username, email, password)
	return args.Get(0).(*UserDomain.User), args.String(1), args.Error(2)
}

// CheckUserByBothTokens simulates checking a user by both an authentication token and a verification token.
func (m *MockUserService) CheckUserByBothTokens(token string, verificationToken string) (*UserDomain.User, error) {
	args := m.Called(token, verificationToken)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

// CheckEmailVerificationToken mocks CheckEmailVerificationToken.
func (m *MockUserService) CheckEmailVerificationToken(token string) (*UserDomain.User, error) {
	args := m.Called(token)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

// CheckAccountVerificationToken mocks CheckAccountVerificationToken.
func (m *MockUserService) CheckAccountVerificationToken(token string) (*UserDomain.User, error) {
	args := m.Called(token)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

// AuthenticateUser mocks AuthenticateUser.
func (m *MockUserService) AuthenticateUser(username string, password string) (*UserDomain.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

// VerifyUserWithTokens mocks VerifyUserWithTokens.
func (m *MockUserService) VerifyUserWithTokens(token string, verificationToken string) (*UserDomain.User, error) {
	args := m.Called(token, verificationToken)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

// CheckIfUserExists mocks CheckIfUserExists
func (m *MockUserService) CheckIfUserExists(email string) (*UserDomain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

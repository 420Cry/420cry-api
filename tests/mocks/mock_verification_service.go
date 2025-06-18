package mocks

import (
	UserModel "cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

// MockVerificationService mocks VerificationService methods
type MockVerificationService struct {
	mock.Mock
}

// VerifyUserWithTokens mocks VerifyUserWithTokens method from VerificationService
func (m *MockVerificationService) VerifyUserWithTokens(token, verificationToken string) (*UserModel.User, error) {
	args := m.Called(token, verificationToken)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// CheckAccountVerificationToken mocks CheckAccountVerificationToken method from VerificationService
func (m *MockVerificationService) CheckAccountVerificationToken(token string) (*UserModel.User, error) {
	args := m.Called(token)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// CheckUserByBothTokens mocks CheckUserByBothTokens method from VerificationService
func (m *MockVerificationService) CheckUserByBothTokens(emailVerificationToken, verificationToken string) (*UserModel.User, error) {
	args := m.Called(emailVerificationToken, verificationToken)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// CheckEmailVerificationToken mocks CheckEmailVerificationToken method from VerificationService
func (m *MockVerificationService) CheckEmailVerificationToken(emailVerificationToken string) (*UserModel.User, error) {
	args := m.Called(emailVerificationToken)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

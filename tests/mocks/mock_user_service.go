package mocks

import (
	UserModel "cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

// MockUserService mocks UserServiceInterface
type MockUserService struct {
	mock.Mock
}

// CreateUser mocks CreateUser from  UserService
func (m *MockUserService) CreateUser(fullname, username, email, password string) (*UserModel.User, error) {
	args := m.Called(fullname, username, email, password)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// AuthenticateUser mocks AuthenticateUser from UserService
func (m *MockUserService) AuthenticateUser(username, password string) (*UserModel.User, error) {
	args := m.Called(username, password)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// VerifyUserWithTokens mocks VerifyUserWithTokens from UserService
func (m *MockUserService) VerifyUserWithTokens(userToken, verifyToken string) (*UserModel.User, error) {
	args := m.Called(userToken, verifyToken)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// CheckAccountVerificationToken mocks CheckAccountVerificationToken from UserService
func (m *MockUserService) CheckAccountVerificationToken(token string) (*UserModel.User, error) {
	args := m.Called(token)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// CheckEmailVerificationToken mocks CheckEmailVerificationToken from UserService
func (m *MockUserService) CheckEmailVerificationToken(token string) (*UserModel.User, error) {
	args := m.Called(token)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// CheckUserByBothTokens mocks CheckUserByBothTokens from UserService
func (m *MockUserService) CheckUserByBothTokens(token, verificationToken string) (*UserModel.User, error) {
	args := m.Called(token, verificationToken)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

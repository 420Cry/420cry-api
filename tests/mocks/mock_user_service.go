package mocks

import (
	"github.com/stretchr/testify/mock"

	UserModel "cry-api/app/models"
	UserTypes "cry-api/app/types/users"
)

// MockUserService mocks UserServiceInterface
type MockUserService struct {
	mock.Mock
}

// CreateUser mocks CreateUser from UserService
func (m *MockUserService) CreateUser(fullname, username, email, password string) (*UserModel.User, error) {
	args := m.Called(fullname, username, email, password)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// GetUserByUUID mocks GetUserByUUID from UserService
func (m *MockUserService) GetUserByUUID(uuid string) (*UserModel.User, error) {
	args := m.Called(uuid)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// UpdateUser mocks UpdateUser from UserService
func (m *MockUserService) UpdateUser(user *UserModel.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// CheckIfUserExists mocks CheckIfUserExists from UserService
func (m *MockUserService) CheckIfUserExists(email string) (*UserModel.User, error) {
	args := m.Called(email)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// CheckUserByResetPasswordToken mocks CheckUserByResetPasswordToken from userService
func (m *MockUserService) CheckUserByResetPasswordToken(token string) (*UserModel.User, error) {
	args := m.Called(token)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// CreateResetPasswordToken mocks CreateResetPasswordToken from userService
func (m *MockUserService) CreateResetPasswordToken(foundUser *UserModel.User) (*UserModel.User, error) {
	args := m.Called(foundUser)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// HandleResetPassword mocks HandleResetPassword from userService
func (m *MockUserService) HandleResetPassword(foundUser *UserModel.User, req *UserTypes.IVerificationResetPasswordForm) error {
	args := m.Called(foundUser, req)
	return args.Error(1)
}

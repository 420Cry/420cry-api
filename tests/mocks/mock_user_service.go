package mocks

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	UserModel "cry-api/app/models"
	OAuthType "cry-api/app/types/oauth"
)

// MockUserService mocks UserServiceInterface
type MockUserService struct {
	mock.Mock
}

// CreateUser mocks CreateUser from UserService
func (m *MockUserService) CreateUser(fullname, username, email, password string, isVerified, isProfileCompleted bool) (*UserModel.User, error) {
	args := m.Called(fullname, username, email, password, isVerified, isProfileCompleted)
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

// FindUserByEmail mocks FindUserByEmail from UserService
func (m *MockUserService) FindUserByEmail(email string) (*UserModel.User, error) {
	args := m.Called(email)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// FindUserByID mocks FindUserByID from UserService
func (m *MockUserService) FindUserByID(id int) (*UserModel.User, error) {
	args := m.Called(id)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// FindUserTokenByPurpose mocks FindUserTokenByPurpose from UserService
func (m *MockUserService) FindUserTokenByPurpose(userID int, purpose string) (*UserModel.UserToken, error) {
	args := m.Called(userID, purpose)
	token, _ := args.Get(0).(*UserModel.UserToken)
	return token, args.Error(1)
}

// FindUserTokenByValueAndPurpose mocks FindUserTokenByValueAndPurpose from UserService
func (m *MockUserService) FindUserTokenByValueAndPurpose(tokenValue, purpose string) (*UserModel.UserToken, error) {
	args := m.Called(tokenValue, purpose)
	token, _ := args.Get(0).(*UserModel.UserToken)
	return token, args.Error(1)
}

// FindUserByUsername mocks FindUserByUsername from UserService
func (m *MockUserService) FindUserByUsername(username string) (*UserModel.User, error) {
	args := m.Called(username)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// CreateUserByGoogleAuth mocks CreateUserByGoogleAuth from UserService
func (m *MockUserService) CreateUserByGoogleAuth(googleUserInfo *OAuthType.IGoogleUserResponse, token *oauth2.Token) (*UserModel.User, error) {
	args := m.Called(googleUserInfo, token)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}
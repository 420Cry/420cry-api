package tests

import (
	"errors"
	"testing"

	UserService "cry-api/app/services/users"
	mocks "cry-api/tests/mocks"

	UserModel "cry-api/app/models"
	SignUpError "cry-api/app/types/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test creating a new user successfully
func TestUserService_CreateUser_NewUser_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserTokenRepo := new(mocks.MockUserTokenRepository)
	mockEmailSvc := new(mocks.MockEmailService)
	mockAuthService := new(mocks.MockAuthService)

	userSvc := UserService.NewUserService(mockUserRepo, mockUserTokenRepo, mockEmailSvc, mockAuthService)

	fullname := "John Doe"
	username := "johndoe"
	email := "john@example.com"
	password := "password123"

	mockUserRepo.On("FindByUsernameOrEmail", username, email).Return(nil, nil)
	mockUserRepo.On("Save", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := userSvc.CreateUser(fullname, username, email, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)

	mockUserRepo.AssertExpectations(t)
}

// Test creating user when user already exists returns conflict error
func TestUserService_CreateUser_UserExists_ReturnsConflict(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserTokenRepo := new(mocks.MockUserTokenRepository)
	mockEmailSvc := new(mocks.MockEmailService)
	mockAuthService := new(mocks.MockAuthService)

	userSvc := UserService.NewUserService(mockUserRepo, mockUserTokenRepo, mockEmailSvc, mockAuthService)

	existingUser := &UserModel.User{Username: "johndoe"}

	mockUserRepo.On("FindByUsernameOrEmail", "johndoe", "john@example.com").Return(existingUser, nil)

	user, err := userSvc.CreateUser("John Doe", "johndoe", "john@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, SignUpError.ErrUserConflict, err)

	mockUserRepo.AssertExpectations(t)
}

// Test FindByUsernameOrEmail error
func TestUserService_CreateUser_FindByUsernameOrEmail_Error(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserTokenRepo := new(mocks.MockUserTokenRepository)
	mockEmailSvc := new(mocks.MockEmailService)
	mockAuthService := new(mocks.MockAuthService)

	userSvc := UserService.NewUserService(mockUserRepo, mockUserTokenRepo, mockEmailSvc, mockAuthService)

	mockUserRepo.On("FindByUsernameOrEmail", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	user, err := userSvc.CreateUser("John Doe", "johndoe", "john@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, user)

	mockUserRepo.AssertExpectations(t)
}

// Test Save error
func TestUserService_CreateUser_Save_Error(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserTokenRepo := new(mocks.MockUserTokenRepository)
	mockEmailSvc := new(mocks.MockEmailService)
	mockAuthService := new(mocks.MockAuthService)

	userSvc := UserService.NewUserService(mockUserRepo, mockUserTokenRepo, mockEmailSvc, mockAuthService)

	mockUserRepo.On("FindByUsernameOrEmail", mock.Anything, mock.Anything).Return(nil, nil)
	mockUserRepo.On("Save", mock.AnythingOfType("*models.User")).Return(errors.New("save error"))

	user, err := userSvc.CreateUser("John Doe", "johndoe", "john@example.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, user)

	mockUserRepo.AssertExpectations(t)
}

// Test GetUserByUUID success
func TestUserService_GetUserByUUID_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserTokenRepo := new(mocks.MockUserTokenRepository)
	mockEmailSvc := new(mocks.MockEmailService)
	mockAuthService := new(mocks.MockAuthService)

	userSvc := UserService.NewUserService(mockUserRepo, mockUserTokenRepo, mockEmailSvc, mockAuthService)

	expectedUser := &UserModel.User{UUID: "uuid-1234"}

	mockUserRepo.On("FindByUUID", "uuid-1234").Return(expectedUser, nil)

	user, err := userSvc.GetUserByUUID("uuid-1234")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockUserRepo.AssertExpectations(t)
}

// Test GetUserByUUID returns error
func TestUserService_GetUserByUUID_Error(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserTokenRepo := new(mocks.MockUserTokenRepository)
	mockEmailSvc := new(mocks.MockEmailService)
	mockAuthService := new(mocks.MockAuthService)

	userSvc := UserService.NewUserService(mockUserRepo, mockUserTokenRepo, mockEmailSvc, mockAuthService)

	mockUserRepo.On("FindByUUID", "uuid-1234").Return(nil, errors.New("db error"))

	user, err := userSvc.GetUserByUUID("uuid-1234")

	assert.Error(t, err)
	assert.Nil(t, user)

	mockUserRepo.AssertExpectations(t)
}

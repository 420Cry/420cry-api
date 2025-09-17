package auth_test

import (
	"errors"
	"testing"

	UserModel "cry-api/app/models"
	AuthService "cry-api/app/services/auth"
	mocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_AuthenticateUser_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPasswordSvc := new(mocks.MockPasswordService)

	authSvc := AuthService.NewAuthService(mockUserRepo, mockPasswordSvc)

	user := &UserModel.User{
		Username:   "johndoe",
		Password:   "hashedpassword",
		IsVerified: true,
	}

	mockUserRepo.On("FindByUsername", "johndoe").Return(user, nil)
	mockPasswordSvc.On("CheckPassword", "hashedpassword", "password123").Return(nil)

	result, err := authSvc.AuthenticateUser("johndoe", "password123")

	assert.NoError(t, err)
	assert.Equal(t, user, result)

	mockUserRepo.AssertExpectations(t)
	mockPasswordSvc.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser_UserNotFound(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPasswordSvc := new(mocks.MockPasswordService)

	authSvc := AuthService.NewAuthService(mockUserRepo, mockPasswordSvc)

	mockUserRepo.On("FindByUsername", "unknown").Return(nil, nil)
	// No need to mock CheckPassword here because it won't be called

	result, err := authSvc.AuthenticateUser("unknown", "password123")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())

	mockUserRepo.AssertExpectations(t)
	// No expectations on password service here
}

func TestAuthService_AuthenticateUser_PasswordMismatch(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPasswordSvc := new(mocks.MockPasswordService)

	authSvc := AuthService.NewAuthService(mockUserRepo, mockPasswordSvc)

	user := &UserModel.User{
		Username:   "johndoe",
		Password:   "hashedpassword",
		IsVerified: true,
	}

	mockUserRepo.On("FindByUsername", "johndoe").Return(user, nil)
	mockPasswordSvc.On("CheckPassword", "hashedpassword", "wrongpassword").Return(errors.New("password mismatch"))

	result, err := authSvc.AuthenticateUser("johndoe", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid password", err.Error())

	mockUserRepo.AssertExpectations(t)
	mockPasswordSvc.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser_UserNotVerified(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockPasswordSvc := new(mocks.MockPasswordService)

	authSvc := AuthService.NewAuthService(mockUserRepo, mockPasswordSvc)

	user := &UserModel.User{
		Username:   "johndoe",
		Password:   "hashedpassword",
		IsVerified: false,
	}

	mockUserRepo.On("FindByUsername", "johndoe").Return(user, nil)
	mockPasswordSvc.On("CheckPassword", "hashedpassword", "password123").Return(nil)

	result, err := authSvc.AuthenticateUser("johndoe", "password123")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not verified", err.Error())

	mockUserRepo.AssertExpectations(t)
	mockPasswordSvc.AssertExpectations(t)
}

package tests

import (
	"errors"
	"testing"

	UserModel "cry-api/app/models"
	VerificationService "cry-api/app/services/users"
	mocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestVerificationService_VerifyUserWithTokens_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	emailToken := "email-token"
	verificationToken := "verify-token"

	user := &UserModel.User{
		IsVerified:               false,
		VerificationTokens:       verificationToken,
		AccountVerificationToken: &emailToken,
	}

	mockUserRepo.On("FindByAccountVerificationToken", emailToken).Return(user, nil)
	mockUserRepo.On("Save", user).Return(nil)

	result, err := svc.VerifyUserWithTokens(emailToken, verificationToken)
	assert.NoError(t, err)
	assert.True(t, result.IsVerified)
	assert.Empty(t, result.VerificationTokens)
	assert.Nil(t, result.AccountVerificationToken)

	mockUserRepo.AssertExpectations(t)
}

func TestVerificationService_VerifyUserWithTokens_ErrorFromCheckUserByBothTokens(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	mockUserRepo.On("FindByAccountVerificationToken", "email-token").Return(nil, errors.New("db error"))

	result, err := svc.VerifyUserWithTokens("email-token", "verify-token")
	assert.Error(t, err)
	assert.Nil(t, result)

	mockUserRepo.AssertExpectations(t)
}

func TestVerificationService_FindUserByAccountVerificationToken_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	token := "account-token"
	user := &UserModel.User{}

	mockUserRepo.On("FindByAccountVerificationToken", token).Return(user, nil)

	result, err := svc.FindUserByAccountVerificationToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	mockUserRepo.AssertExpectations(t)
}

func TestVerificationService_FindUserByAccountVerificationTokenn_UserNotFound(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	mockUserRepo.On("FindByAccountVerificationToken", "invalid-token").Return(nil, nil)

	result, err := svc.FindUserByAccountVerificationToken("invalid-token")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid account token", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestVerificationService_CheckUserByBothTokens_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	emailToken := "email-token"
	verificationToken := "verify-token"

	user := &UserModel.User{
		IsVerified:               false,
		VerificationTokens:       verificationToken,
		AccountVerificationToken: &emailToken,
	}
	mockUserRepo.On("FindByAccountVerificationToken", emailToken).Return(user, nil)

	result, err := svc.CheckUserByBothTokens(emailToken, verificationToken)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	mockUserRepo.AssertExpectations(t)
}

func TestVerificationService_CheckUserByBothTokens_VerificationTokenInvalid(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	mockUserRepo.On("FindByAccountVerificationToken", "email-token").Return(nil, nil)

	result, err := svc.CheckUserByBothTokens("email-token", "bad-token")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid verification token", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestVerificationService_CheckUserByBothTokens_TokenMismatch(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	user := &UserModel.User{
		AccountVerificationToken: nil,
	}
	mockUserRepo.On("FindByAccountVerificationToken", "email-token").Return(user, nil)

	result, err := svc.CheckUserByBothTokens("email-token", "verify-token")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "token does not match", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestVerificationService_CheckEmailVerificationToken_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	emailToken := "email-token"

	user := &UserModel.User{
		IsVerified:         false,
		VerificationTokens: "some-token",
	}

	mockUserRepo.On("FindByAccountVerificationToken", emailToken).Return(user, nil)
	mockUserRepo.On("Save", user).Return(nil)

	result, err := svc.CheckEmailVerificationToken(emailToken)
	assert.NoError(t, err)
	assert.True(t, result.IsVerified)
	assert.Empty(t, result.VerificationTokens)

	mockUserRepo.AssertExpectations(t)
}

func TestVerificationService_CheckEmailVerificationToken_InvalidToken(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	mockUserRepo.On("FindByAccountVerificationToken", "bad-token").Return(nil, nil)

	result, err := svc.CheckEmailVerificationToken("bad-token")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid verification token", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestVerificationService_CheckEmailVerificationToken_SaveFails(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	svc := VerificationService.NewVerificationService(mockUserRepo)

	user := &UserModel.User{
		IsVerified:         false,
		VerificationTokens: "some-token",
	}

	mockUserRepo.On("FindByAccountVerificationToken", "email-token").Return(user, nil)
	mockUserRepo.On("Save", user).Return(errors.New("save error"))

	result, err := svc.CheckEmailVerificationToken("email-token")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "save error", err.Error())

	mockUserRepo.AssertExpectations(t)
}

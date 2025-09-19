package tests

import (
	"errors"
	"testing"
	"time"

	UserModel "cry-api/app/models"
	services "cry-api/app/services/users"
	mocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestSave_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserTokenRepository)
	service := services.NewUserTokenService(mockRepo)

	token := &UserModel.UserToken{Token: "abc123"}
	mockRepo.On("Save", token).Return(nil)

	err := service.Save(token)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Save", token)
}

func TestFindValidToken_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserTokenRepository)
	service := services.NewUserTokenService(mockRepo)

	validToken := &UserModel.UserToken{
		Token:     "abc123",
		Purpose:   "reset",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	mockRepo.On("FindValidToken", "abc123", "reset").Return(validToken, nil)

	result, err := service.FindValidToken("abc123", "reset")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "abc123", result.Token)
}

func TestFindValidToken_Expired(t *testing.T) {
	mockRepo := new(mocks.MockUserTokenRepository)
	service := services.NewUserTokenService(mockRepo)

	expiredToken := &UserModel.UserToken{
		Token:     "expired",
		Purpose:   "reset",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	mockRepo.On("FindValidToken", "expired", "reset").Return(expiredToken, nil)

	result, err := service.FindValidToken("expired", "reset")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFindValidToken_Error(t *testing.T) {
	mockRepo := new(mocks.MockUserTokenRepository)
	service := services.NewUserTokenService(mockRepo)

	mockRepo.On("FindValidToken", "bad", "reset").Return(nil, errors.New("db error"))

	result, err := service.FindValidToken("bad", "reset")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFindLatestValidToken_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserTokenRepository)
	service := services.NewUserTokenService(mockRepo)

	token := &UserModel.UserToken{Token: "latest", Purpose: "verify"}
	mockRepo.On("FindLatestValidToken", 1, "verify").Return(token, nil)

	result, err := service.FindLatestValidToken(1, "verify")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "latest", result.Token)
}

func TestConsumeToken_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserTokenRepository)
	service := services.NewUserTokenService(mockRepo)

	mockRepo.On("ConsumeToken", 1, "abc123", "reset").Return(nil)

	err := service.ConsumeToken(1, "abc123", "reset")
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "ConsumeToken", 1, "abc123", "reset")
}

func TestDeleteExpired_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserTokenRepository)
	service := services.NewUserTokenService(mockRepo)

	mockRepo.On("DeleteExpired").Return(nil)

	err := service.DeleteExpired()
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "DeleteExpired")
}

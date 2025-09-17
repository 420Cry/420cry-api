package user_controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controller "cry-api/app/controllers/users"
	UserModel "cry-api/app/models"
	UserTypes "cry-api/app/types/users"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVerifyResetPasswordToken_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	mockPasswordService := new(testmocks.MockPasswordService)

	userController := &controller.UserController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
		PasswordService:  mockPasswordService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "validtoken123",
		NewPassword:        "newsecurepassword",
	}
	bodyBytes, _ := json.Marshal(input)

	// Mock the reset token
	userToken := &UserModel.UserToken{
		Token:     input.ResetPasswordToken,
		UserID:    42,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockUserTokenService.On("FindValidToken", input.ResetPasswordToken, mock.Anything).Return(userToken, nil)
	mockUserTokenService.On("ConsumeToken", userToken.UserID, userToken.Token, mock.Anything).Return(nil)

	// Mock the user
	user := &UserModel.User{
		ID:         42,
		IsVerified: true,
	}
	mockUserService.On("FindUserByID", userToken.UserID).Return(user, nil)
	mockUserService.On("UpdateUser", user).Return(nil)

	// Mock password hashing
	mockPasswordService.On("HashPassword", input.NewPassword).Return("hashedpassword", nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	var resp map[string]bool
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.True(t, resp["success"])

	mockUserTokenService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
	mockPasswordService.AssertExpectations(t)
}

func TestVerifyResetPasswordToken_InvalidJSON(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	mockPasswordService := new(testmocks.MockPasswordService)

	userController := &controller.UserController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
		PasswordService:  mockPasswordService,
	}

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader([]byte(`{invalid-json}`)))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid JSON format", resp["message"])
}

func TestVerifyResetPasswordToken_TokenNotFound(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	mockPasswordService := new(testmocks.MockPasswordService)

	userController := &controller.UserController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
		PasswordService:  mockPasswordService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "invalidtoken",
		NewPassword:        "password",
	}
	bodyBytes, _ := json.Marshal(input)

	mockUserTokenService.On("FindValidToken", input.ResetPasswordToken, mock.Anything).Return(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid or expired reset password token", resp["message"])
}

func TestVerifyResetPasswordToken_UserNotVerified(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	mockPasswordService := new(testmocks.MockPasswordService)

	userController := &controller.UserController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
		PasswordService:  mockPasswordService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "validtoken",
		NewPassword:        "password",
	}
	bodyBytes, _ := json.Marshal(input)

	userToken := &UserModel.UserToken{
		Token:     input.ResetPasswordToken,
		UserID:    42,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockUserTokenService.On("FindValidToken", input.ResetPasswordToken, mock.Anything).Return(userToken, nil)

	user := &UserModel.User{
		ID:         42,
		IsVerified: false,
	}
	mockUserService.On("FindUserByID", userToken.UserID).Return(user, nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "User is not verified", resp["message"])
}

func TestVerifyResetPasswordToken_HashPasswordFailure(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	mockPasswordService := new(testmocks.MockPasswordService)

	userController := &controller.UserController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
		PasswordService:  mockPasswordService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "validtoken123",
		NewPassword:        "newpassword",
	}
	bodyBytes, _ := json.Marshal(input)

	userToken := &UserModel.UserToken{
		Token:     input.ResetPasswordToken,
		UserID:    42,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockUserTokenService.On("FindValidToken", input.ResetPasswordToken, mock.Anything).Return(userToken, nil)

	user := &UserModel.User{
		ID:         42,
		IsVerified: true,
	}
	mockUserService.On("FindUserByID", userToken.UserID).Return(user, nil)

	mockPasswordService.On("HashPassword", input.NewPassword).Return("", errors.New("hash error"))

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to hash password", resp["message"])
}

package tests

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
	mockEmailService := new(testmocks.MockEmailService)
	mockPasswordService := new(testmocks.MockPasswordService)

	userController := &controller.UserController{
		UserService:     mockUserService,
		EmailService:    mockEmailService,
		PasswordService: mockPasswordService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "validtoken123",
		NewPassword:        "newsecurepassword",
	}
	bodyBytes, _ := json.Marshal(input)

	now := time.Now().Add(-30 * time.Minute)
	dummyUser := &UserModel.User{
		ResetPasswordToken:          input.ResetPasswordToken,
		ResetPasswordTokenCreatedAt: &now,
		IsVerified:                  true,
	}

	mockUserService.
		On("FindUserByResetPasswordToken", input.ResetPasswordToken).
		Return(dummyUser, nil)

	mockPasswordService.
		On("HashPassword", input.NewPassword).
		Return("hashedpassword", nil)

	mockUserService.
		On("UpdateUser", mock.MatchedBy(func(u *UserModel.User) bool {
			// Relaxed matcher: only check the reset token is cleared, password is non-empty (hashed)
			return u.Password == "hashedpassword" && u.ResetPasswordToken == "" && u.ResetPasswordTokenCreatedAt == nil
		})).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var respBody map[string]bool
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.True(t, respBody["success"])

	mockUserService.AssertExpectations(t)
	mockPasswordService.AssertExpectations(t)
}

func TestVerifyResetPasswordToken_InvalidJSON(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	invalidJSON := []byte(`{invalid-json}`)

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(invalidJSON))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid JSON Format", respBody["message"])
}

func TestVerifyResetPasswordToken_UserNotFound(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "invalidtoken",
		NewPassword:        "password",
	}
	bodyBytes, _ := json.Marshal(input)

	mockUserService.
		On("FindUserByResetPasswordToken", input.ResetPasswordToken).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Cannot find user", respBody["message"])
}

func TestVerifyResetPasswordToken_UserNotVerified(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "validtoken",
		NewPassword:        "password",
	}
	bodyBytes, _ := json.Marshal(input)

	dummyUser := &UserModel.User{
		ResetPasswordToken: input.ResetPasswordToken,
		IsVerified:         false,
	}

	mockUserService.
		On("FindUserByResetPasswordToken", input.ResetPasswordToken).
		Return(dummyUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "User is not verified", respBody["message"])
}

func TestVerifyResetPasswordToken_TokenExpired(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "validtoken",
		NewPassword:        "password",
	}
	bodyBytes, _ := json.Marshal(input)

	// Token created more than 1 hour ago => expired
	oldTime := time.Now().Add(-2 * time.Hour)

	dummyUser := &UserModel.User{
		ResetPasswordToken:          input.ResetPasswordToken,
		ResetPasswordTokenCreatedAt: &oldTime,
		IsVerified:                  true,
	}

	mockUserService.
		On("FindUserByResetPasswordToken", input.ResetPasswordToken).
		Return(dummyUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "The token has been expired", respBody["message"])
}

func TestVerifyResetPasswordToken_HashPasswordFailure(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockPasswordService := new(testmocks.MockPasswordService)

	userController := &controller.UserController{
		UserService:     mockUserService,
		PasswordService: mockPasswordService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "validtoken123",
		NewPassword:        "newsecurepassword",
	}
	bodyBytes, _ := json.Marshal(input)

	now := time.Now().Add(-30 * time.Minute)
	dummyUser := &UserModel.User{
		ResetPasswordToken:          input.ResetPasswordToken,
		ResetPasswordTokenCreatedAt: &now,
		IsVerified:                  true,
	}

	mockUserService.
		On("FindUserByResetPasswordToken", input.ResetPasswordToken).
		Return(dummyUser, nil)

	mockPasswordService.
		On("HashPassword", input.NewPassword).
		Return("", errors.New("hash error"))

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Cannot hash password", respBody["message"])

	mockUserService.AssertExpectations(t)
	mockPasswordService.AssertExpectations(t)
}

func TestVerifyResetPasswordToken_UpdateUserFailure(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockPasswordService := new(testmocks.MockPasswordService)

	userController := &controller.UserController{
		UserService:     mockUserService,
		PasswordService: mockPasswordService,
	}

	input := UserTypes.IVerificationResetPasswordForm{
		ResetPasswordToken: "validtoken123",
		NewPassword:        "newsecurepassword",
	}
	bodyBytes, _ := json.Marshal(input)

	now := time.Now().Add(-30 * time.Minute)
	dummyUser := &UserModel.User{
		ResetPasswordToken:          input.ResetPasswordToken,
		ResetPasswordTokenCreatedAt: &now,
		IsVerified:                  true,
	}

	mockUserService.
		On("FindUserByResetPasswordToken", input.ResetPasswordToken).
		Return(dummyUser, nil)

	mockPasswordService.
		On("HashPassword", input.NewPassword).
		Return("hashedpassword", nil)

	mockUserService.
		On("UpdateUser", mock.Anything).
		Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/verify-reset-password-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyResetPasswordToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to update user", respBody["message"])

	mockUserService.AssertExpectations(t)
	mockPasswordService.AssertExpectations(t)
}

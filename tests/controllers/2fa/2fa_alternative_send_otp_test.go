// Package tests provides tests for 2FA routes.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controller "cry-api/app/controllers/2fa"
	UserModel "cry-api/app/models"
	TwoFactorTypes "cry-api/app/types/2fa"
	TokenTypes "cry-api/app/types/token_purpose"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAlternativeSendOtp_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	mockEmailService := new(testmocks.MockEmailService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
		EmailService:     mockEmailService,
	}

	input := TwoFactorTypes.ITwoFactorAlternativeRequest{
		Email: "john@example.com",
	}

	bodyBytes, _ := json.Marshal(input)

	dummyUser := &UserModel.User{
		ID:         1,
		Email:      input.Email,
		Username:   "johndoe",
		IsVerified: true,
	}

	dummyToken := &UserModel.UserToken{
		UserID:    dummyUser.ID,
		Token:     "dummy-otp-token",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// Channel to signal async email sending
	done := make(chan struct{})

	// Mock: find user
	mockUserService.
		On("FindUserByEmail", input.Email).
		Return(dummyUser, nil)

	// Mock: no existing OTP token
	mockUserTokenService.
		On("FindLatestValidToken", dummyUser.ID, string(TokenTypes.TwoFactorAuthAlternativeOTP)).
		Return(nil, nil)

	// Mock: saving new OTP
	mockUserTokenService.
		On("Save", mock.AnythingOfType("*models.UserToken")).
		Return(nil).
		Run(func(args mock.Arguments) {
			token := args.Get(0).(*UserModel.UserToken)
			token.Token = dummyToken.Token
		})

	// Mock: sending OTP asynchronously
	mockEmailService.
		On("SendTwoFactorAlternativeEmail",
			dummyUser.Email,
			mock.Anything,
			dummyUser.Username,
			mock.Anything,
			5,
		).
		Return(nil).
		Run(func(_ mock.Arguments) {
			close(done) // signal that email was sent
		})

	// Make HTTP request
	req := httptest.NewRequest(http.MethodPost, "/2fa/alternative", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)

	controller.AlternativeSendOtp(c)

	// Wait for async email to be called (or timeout)
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected SendTwoFactorAlternativeEmail to be called, but it wasn’t")
	}

	// Verify response
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var respBody map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, true, respBody["success"])
	assert.Equal(t, "OTP sent successfully", respBody["message"])

	// Assert all mocks were called as expected
	mockUserService.AssertExpectations(t)
	mockUserTokenService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestAlternativeSendOtp_InvalidJSON(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	mockEmailService := new(testmocks.MockEmailService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
		EmailService:     mockEmailService,
	}

	invalidJSON := []byte(`{invalid-json}`)

	req := httptest.NewRequest(http.MethodPost, "/2fa/alternative", bytes.NewReader(invalidJSON))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	controller.AlternativeSendOtp(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid request payload")

	mockUserService.AssertNotCalled(t, "FindUserByEmail", mock.Anything)
}

func TestAlternativeSendOtp_UserNotFound(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	mockEmailService := new(testmocks.MockEmailService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
		EmailService:     mockEmailService,
	}

	input := TwoFactorTypes.ITwoFactorAlternativeRequest{
		Email: "notfound@example.com",
	}

	bodyBytes, _ := json.Marshal(input)

	mockUserService.On("FindUserByEmail", input.Email).Return(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/2fa/alternative", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	controller.AlternativeSendOtp(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "User not found")

	mockUserService.AssertExpectations(t)
}

func TestAlternativeSendOtp_UserNotVerified(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	mockEmailService := new(testmocks.MockEmailService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
		EmailService:     mockEmailService,
	}

	input := TwoFactorTypes.ITwoFactorAlternativeRequest{
		Email: "unverified@example.com",
	}

	user := &UserModel.User{
		ID:         2,
		Email:      input.Email,
		Username:   "unverified",
		IsVerified: false,
	}

	bodyBytes, _ := json.Marshal(input)

	mockUserService.On("FindUserByEmail", input.Email).Return(user, nil)

	req := httptest.NewRequest(http.MethodPost, "/2fa/alternative", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	controller.AlternativeSendOtp(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusForbidden, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "User not verified")

	mockUserService.AssertExpectations(t)
}

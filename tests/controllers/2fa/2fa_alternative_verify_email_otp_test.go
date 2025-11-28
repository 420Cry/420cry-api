// Package tests provides tests for 2FA routes.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "cry-api/app/controllers/2fa"
	UserModel "cry-api/app/models"
	JWT "cry-api/app/services/jwt"
	TwoFactorTypes "cry-api/app/types/2fa"
	TokenType "cry-api/app/types/token_purpose"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestAlternativeVerifyOTP_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	input := TwoFactorTypes.ITwoFactorVerifyRequest{
		UserUUID: "uuid-1234",
		OTP:      "valid-otp",
	}

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}

	user := &UserModel.User{
		ID:           1,
		UUID:         input.UserUUID,
		Email:        "john@example.com",
		TwoFAEnabled: true,
	}

	token := &UserModel.UserToken{
		UserID: user.ID,
		Token:  "valid-otp",
	}

	mockUserService.On("GetUserByUUID", input.UserUUID).Return(user, nil)
	mockUserTokenService.On("FindLatestValidToken", user.ID, string(TokenType.TwoFactorAuthAlternativeOTP)).
		Return(token, nil)
	mockUserTokenService.On("ConsumeToken", user.ID, input.OTP, string(TokenType.TwoFactorAuthAlternativeOTP)).
		Return(nil)

	// Mock JWT
	oldGenerateJWT := JWT.GenerateJWT
	defer func() { JWT.GenerateJWT = oldGenerateJWT }()
	JWT.GenerateJWT = func(_, _ string, _, _ bool) (string, error) {
		return "mocked.jwt.token", nil
	}

	req := httptest.NewRequest(http.MethodPost, "/2fa/verify-alternative", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)

	controller.AlternativeVerifyOTP(c)

	res := w.Result()
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			t.Fatalf("failed to close response body: %v", cerr)
		}
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var respBody map[string]string
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	assert.Equal(t, "mocked.jwt.token", respBody["jwt"])

	mockUserService.AssertExpectations(t)
	mockUserTokenService.AssertExpectations(t)
}

func TestAlternativeVerifyOTP_InvalidJSON(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	req := httptest.NewRequest(http.MethodPost, "/2fa/verify-alternative", bytes.NewReader([]byte(`{invalid-json}`)))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)

	controller.AlternativeVerifyOTP(c)

	res := w.Result()
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			t.Fatalf("failed to close response body: %v", cerr)
		}
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	assert.Contains(t, respBody["error"], "Invalid request")
}

func TestAlternativeVerifyOTP_MissingFields(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	tests := []struct {
		name       string
		payload    TwoFactorTypes.ITwoFactorVerifyRequest
		wantErrMsg string
	}{
		{"Missing UUID", TwoFactorTypes.ITwoFactorVerifyRequest{OTP: "otp"}, "User UUID is required"},
		{"Missing OTP", TwoFactorTypes.ITwoFactorVerifyRequest{UserUUID: "uuid-123"}, "OTP is required for verification"},
	}

	for _, tt := range tests {
		bodyBytes, err := json.Marshal(tt.payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/2fa/verify-alternative", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()
		c := TestUtils.GetGinContext(w, req)

		controller.AlternativeVerifyOTP(c)

		res := w.Result()
		defer func() {
			if cerr := res.Body.Close(); cerr != nil {
				t.Fatalf("failed to close response body: %v", cerr)
			}
		}()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		var respBody map[string]string
		if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		assert.Contains(t, respBody["error"], tt.wantErrMsg)
	}
}

func TestAlternativeVerifyOTP_UserNotFound(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	input := TwoFactorTypes.ITwoFactorVerifyRequest{
		UserUUID: "uuid-1234",
		OTP:      "otp",
	}

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}

	mockUserService.On("GetUserByUUID", input.UserUUID).Return(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/2fa/verify-alternative", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)

	controller.AlternativeVerifyOTP(c)

	res := w.Result()
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			t.Fatalf("failed to close response body: %v", cerr)
		}
	}()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	var respBody map[string]string
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	assert.Contains(t, respBody["error"], "User not found")

	mockUserService.AssertExpectations(t)
}

func TestAlternativeVerifyOTP_User2FANotEnabled(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	input := TwoFactorTypes.ITwoFactorVerifyRequest{
		UserUUID: "uuid-1234",
		OTP:      "otp",
	}

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}

	user := &UserModel.User{
		ID:           1,
		UUID:         input.UserUUID,
		TwoFAEnabled: false,
	}

	mockUserService.On("GetUserByUUID", input.UserUUID).Return(user, nil)

	req := httptest.NewRequest(http.MethodPost, "/2fa/verify-alternative", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)

	controller.AlternativeVerifyOTP(c)

	res := w.Result()
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			t.Fatalf("failed to close response body: %v", cerr)
		}
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if val, ok := respBody["error"]; ok {
		assert.Contains(t, val, "User has not enabled 2FA")
	} else {
		t.Fatalf("response body missing 'error' field")
	}

	mockUserService.AssertExpectations(t)
}

func TestAlternativeVerifyOTP_InvalidOrExpiredOTP(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	input := TwoFactorTypes.ITwoFactorVerifyRequest{
		UserUUID: "uuid-1234",
		OTP:      "wrong-otp",
	}

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}

	user := &UserModel.User{
		ID:           1,
		UUID:         input.UserUUID,
		TwoFAEnabled: true,
	}

	mockUserService.On("GetUserByUUID", input.UserUUID).Return(user, nil)

	mockUserTokenService.On("FindLatestValidToken", user.ID, string(TokenType.TwoFactorAuthAlternativeOTP)).
		Return(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/2fa/verify-alternative", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)

	controller.AlternativeVerifyOTP(c)

	res := w.Result()
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			t.Fatalf("failed to close response body: %v", cerr)
		}
	}()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var respBody map[string]string
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Contains(t, respBody["error"], "Invalid or expired OTP")
	mockUserService.AssertExpectations(t)
	mockUserTokenService.AssertExpectations(t)
}

func TestAlternativeVerifyOTP_WrongOTP(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)

	controller := &controller.TwoFactorController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	input := TwoFactorTypes.ITwoFactorVerifyRequest{
		UserUUID: "uuid-1234",
		OTP:      "wrong-otp",
	}

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}

	user := &UserModel.User{
		ID:           1,
		UUID:         input.UserUUID,
		TwoFAEnabled: true,
	}

	existingToken := &UserModel.UserToken{
		UserID: user.ID,
		Token:  "correct-otp",
	}

	mockUserService.On("GetUserByUUID", input.UserUUID).Return(user, nil)
	mockUserTokenService.On("FindLatestValidToken", user.ID, string(TokenType.TwoFactorAuthAlternativeOTP)).
		Return(existingToken, nil)

	req := httptest.NewRequest(http.MethodPost, "/2fa/verify-alternative", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)

	controller.AlternativeVerifyOTP(c)

	res := w.Result()
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			t.Fatalf("failed to close response body: %v", cerr)
		}
	}()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var respBody map[string]string
	if err := json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	assert.Contains(t, respBody["error"], "Invalid or expired OTP")
	mockUserService.AssertExpectations(t)
	mockUserTokenService.AssertExpectations(t)
}

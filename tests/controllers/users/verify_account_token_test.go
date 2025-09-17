package user_controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controller "cry-api/app/controllers/users"
	UserModel "cry-api/app/models"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestVerifyAccountToken_Success(t *testing.T) {
	mockUserTokenService := new(testmocks.MockUserTokenService)
	userController := &controller.UserController{
		UserTokenService: mockUserTokenService,
	}

	token := "valid-token-123"
	userToken := &UserModel.UserToken{
		Token:     token,
		UserID:    42,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Purpose:   "account_verification",
	}

	mockUserTokenService.
		On("FindValidToken", token, "account_verification").
		Return(userToken, nil)

	bodyBytes, _ := json.Marshal(map[string]string{"token": token})
	req := httptest.NewRequest(http.MethodPost, "/verify-account-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyAccountToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&resp)
	assert.Equal(t, float64(userToken.UserID), resp["user_id"])
	assert.True(t, resp["valid"].(bool))

	mockUserTokenService.AssertExpectations(t)
}

func TestVerifyAccountToken_InvalidJSON(t *testing.T) {
	mockUserTokenService := new(testmocks.MockUserTokenService)
	userController := &controller.UserController{
		UserTokenService: mockUserTokenService,
	}

	req := httptest.NewRequest(http.MethodPost, "/verify-account-token", bytes.NewReader([]byte(`{invalid-json}`)))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyAccountToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp map[string]string
	_ = json.NewDecoder(res.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "Invalid request body")
}

func TestVerifyAccountToken_TokenNotFound(t *testing.T) {
	mockUserTokenService := new(testmocks.MockUserTokenService)
	userController := &controller.UserController{
		UserTokenService: mockUserTokenService,
	}

	token := "nonexistent-token"

	mockUserTokenService.
		On("FindValidToken", token, "account_verification").
		Return((*UserModel.UserToken)(nil), nil)

	bodyBytes, _ := json.Marshal(map[string]string{"token": token})
	req := httptest.NewRequest(http.MethodPost, "/verify-account-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyAccountToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp map[string]string
	_ = json.NewDecoder(res.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "Token is invalid or expired")

	mockUserTokenService.AssertExpectations(t)
}

func TestVerifyAccountToken_ExpiredToken(t *testing.T) {
	mockUserTokenService := new(testmocks.MockUserTokenService)
	userController := &controller.UserController{
		UserTokenService: mockUserTokenService,
	}

	token := "expired-token"
	userToken := &UserModel.UserToken{
		Token:     token,
		UserID:    42,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		Purpose:   "account_verification",
	}

	mockUserTokenService.
		On("FindValidToken", token, "account_verification").
		Return(userToken, nil)

	bodyBytes, _ := json.Marshal(map[string]string{"token": token})
	req := httptest.NewRequest(http.MethodPost, "/verify-account-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyAccountToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp map[string]string
	_ = json.NewDecoder(res.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "Token has expired")

	mockUserTokenService.AssertExpectations(t)
}

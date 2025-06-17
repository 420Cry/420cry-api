package tests

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

	"github.com/stretchr/testify/assert"
)

func TestVerifyAccountToken_Success(t *testing.T) {
	mockUserService := new(MockUserService)

	userController := &controller.UserController{
		UserService: mockUserService,
	}

	token := "valid-token-123"
	user := &UserModel.User{
		Token:                      &token,
		VerificationTokenCreatedAt: time.Now(),
	}

	mockUserService.On("CheckAccountVerificationToken", token).Return(user, nil)

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

	var resp map[string]bool
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.True(t, resp["valid"])

	mockUserService.AssertExpectations(t)
}

func TestVerifyAccountToken_InvalidJSON(t *testing.T) {
	mockUserService := new(MockUserService)

	userController := &controller.UserController{
		UserService: mockUserService,
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
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "Invalid request body")
}

func TestVerifyAccountToken_UserNotFound(t *testing.T) {
	mockUserService := new(MockUserService)

	userController := &controller.UserController{
		UserService: mockUserService,
	}

	token := "nonexistent-token"
	mockUserService.On("CheckAccountVerificationToken", token).Return((*UserModel.User)(nil), assert.AnError)

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
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "Token is invalid or expired")

	mockUserService.AssertExpectations(t)
}

func TestVerifyAccountToken_TokenMismatchOrExpired(t *testing.T) {
	mockUserService := new(MockUserService)

	userController := &controller.UserController{
		UserService: mockUserService,
	}

	// Simulate mismatch or expiration
	requestToken := "valid-token-123"
	storedToken := "different-token"
	user := &UserModel.User{
		Token:                      &storedToken,
		VerificationTokenCreatedAt: time.Now().Add(-25 * time.Hour), // expired
	}

	mockUserService.On("CheckAccountVerificationToken", requestToken).Return(user, nil)

	bodyBytes, _ := json.Marshal(map[string]string{"token": requestToken})
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
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "Token is invalid or expired")

	mockUserService.AssertExpectations(t)
}

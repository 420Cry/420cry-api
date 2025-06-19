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
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestVerifyAccountToken_Success(t *testing.T) {
	mockAuthService := new(testmocks.MockAuthService)
	mockVerificationService := new(testmocks.MockVerificationService)
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:         mockUserService,
		EmailService:        mockEmailService,
		VerificationService: mockVerificationService,
		AuthService:         mockAuthService,
	}

	token := "valid-token-123"
	user := &UserModel.User{
		AccountVerificationToken:   &token,
		VerificationTokenCreatedAt: time.Now(),
	}

	mockVerificationService.On("CheckAccountVerificationToken", token).Return(user, nil)

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

	mockVerificationService.AssertExpectations(t)
}

func TestVerifyAccountToken_InvalidJSON(t *testing.T) {
	mockAuthService := new(testmocks.MockAuthService)
	mockVerificationService := new(testmocks.MockVerificationService)
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:         mockUserService,
		EmailService:        mockEmailService,
		VerificationService: mockVerificationService,
		AuthService:         mockAuthService,
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
	mockAuthService := new(testmocks.MockAuthService)
	mockVerificationService := new(testmocks.MockVerificationService)
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:         mockUserService,
		EmailService:        mockEmailService,
		VerificationService: mockVerificationService,
		AuthService:         mockAuthService,
	}

	token := "nonexistent-token"
	mockVerificationService.On("CheckAccountVerificationToken", token).Return((*UserModel.User)(nil), assert.AnError)

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

	mockVerificationService.AssertExpectations(t)
}

func TestVerifyAccountToken_TokenMismatchOrExpired(t *testing.T) {
	mockAuthService := new(testmocks.MockAuthService)
	mockVerificationService := new(testmocks.MockVerificationService)
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:         mockUserService,
		EmailService:        mockEmailService,
		VerificationService: mockVerificationService,
		AuthService:         mockAuthService,
	}

	// Simulate mismatch or expiration
	requestToken := "valid-token-123"
	storedToken := "different-token"
	user := &UserModel.User{
		AccountVerificationToken:   &storedToken,
		VerificationTokenCreatedAt: time.Now().Add(-25 * time.Hour), // expired
	}

	mockVerificationService.On("CheckAccountVerificationToken", requestToken).Return(user, nil)

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

	mockVerificationService.AssertExpectations(t)
}

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "cry-api/app/controllers/users"
	UserModel "cry-api/app/models"
	UserTypes "cry-api/app/types/users"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestVerifyEmailToken_Success(t *testing.T) {
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

	reqBody := UserTypes.IVerificationTokenCheckRequest{
		UserToken:   "user-token-123",
		VerifyToken: "verify-token-abc",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	dummyUser := &UserModel.User{
		IsVerified: true,
	}

	mockVerificationService.
		On("VerifyUserWithTokens", reqBody.UserToken, reqBody.VerifyToken).
		Return(dummyUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-email-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyEmailToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp map[string]bool
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.True(t, resp["verified"])

	mockVerificationService.AssertExpectations(t)
}

func TestVerifyEmailToken_InvalidJSON(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodPost, "/verify-email-token", bytes.NewReader([]byte(`{invalid-json}`)))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyEmailToken(c)

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

func TestVerifyEmailToken_VerificationFails(t *testing.T) {
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

	reqBody := UserTypes.IVerificationTokenCheckRequest{
		UserToken:   "user-token-123",
		VerifyToken: "verify-token-abc",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockVerificationService.
		On("VerifyUserWithTokens", reqBody.UserToken, reqBody.VerifyToken).
		Return((*UserModel.User)(nil), assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/verify-email-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyEmailToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], assert.AnError.Error())

	mockVerificationService.AssertExpectations(t)
}

// Package tests provides tests for user routes.
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "cry-api/app/controllers/users"
	UserModel "cry-api/app/models"
	SignInError "cry-api/app/types/errors"
	UserTypes "cry-api/app/types/users"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignIn_Success(t *testing.T) {
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

	input := UserTypes.IUserSigninRequest{
		Username: "johndoe",
		Password: "securepassword",
	}

	bodyBytes, _ := json.Marshal(input)

	dummyUser := &UserModel.User{
		UUID:     "uuid-1234",
		Fullname: "John Doe",
		Email:    "john@example.com",
		Username: input.Username,
	}

	// Setup mock expectations
	mockAuthService.
		On("AuthenticateUser", input.Username, input.Password).
		Return(dummyUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.SignIn(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var respBody map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)

	// Validate JWT presence and user details in response
	assert.Contains(t, respBody, "jwt")
	userData := respBody["user"].(map[string]interface{})
	assert.Equal(t, dummyUser.UUID, userData["uuid"])
	assert.Equal(t, dummyUser.Fullname, userData["fullname"])
	assert.Equal(t, dummyUser.Email, userData["email"])
	assert.Equal(t, dummyUser.Username, userData["username"])

	mockUserService.AssertExpectations(t)
}

func TestSignIn_InvalidJSON(t *testing.T) {
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

	invalidJSON := []byte(`{invalid-json}`) // malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(invalidJSON))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.SignIn(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid JSON")

	mockAuthService.AssertNotCalled(t, "AuthenticateUser", mock.Anything, mock.Anything)
}

func TestSignIn_UserNotFound(t *testing.T) {
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

	input := UserTypes.IUserSigninRequest{
		Username: "nonexistentuser",
		Password: "anyPassword",
	}

	bodyBytes, _ := json.Marshal(input)

	// Return ErrUserNotFound
	mockAuthService.
		On("AuthenticateUser", input.Username, input.Password).
		Return((*UserModel.User)(nil), SignInError.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.SignIn(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid email or password")

	mockAuthService.AssertExpectations(t)
}

func TestSignIn_InvalidPassword(t *testing.T) {
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

	input := UserTypes.IUserSigninRequest{
		Username: "existinguser",
		Password: "wrongpassword",
	}

	bodyBytes, _ := json.Marshal(input)

	// Return ErrInvalidPassword
	mockAuthService.
		On("AuthenticateUser", input.Username, input.Password).
		Return((*UserModel.User)(nil), SignInError.ErrInvalidPassword)

	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.SignIn(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid email or password")

	mockAuthService.AssertExpectations(t)
}

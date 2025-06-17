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
	UserTypes "cry-api/app/types/users"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignIn_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService) // needed for Handler struct

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
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
	mockUserService.
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
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
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

	mockUserService.AssertNotCalled(t, "AuthenticateUser", mock.Anything, mock.Anything)
}

func TestSignIn_AuthenticationFails(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserSigninRequest{
		Username: "wronguser",
		Password: "wrongpassword",
	}

	bodyBytes, _ := json.Marshal(input)

	// Simulate authentication failure
	mockUserService.
		On("AuthenticateUser", input.Username, input.Password).
		Return((*UserModel.User)(nil), assert.AnError)

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

	mockUserService.AssertExpectations(t)
}

// Package tests provides tests for user routes.
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestSignup_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserSignupRequest{
		Fullname: "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "securepassword",
	}
	bodyBytes, _ := json.Marshal(input)

	dummyToken := "verify123"
	dummyAccountToken := "account123"

	dummyUser := &UserModel.User{
		Email:                    input.Email,
		Username:                 input.Username,
		VerificationTokens:       dummyToken,
		AccountVerificationToken: &dummyAccountToken,
	}

	done := make(chan struct{})

	mockUserService.
		On("CreateUser", input.Fullname, input.Username, input.Email, input.Password).
		Return(dummyUser, nil)

	mockEmailService.
		On("SendVerifyAccountEmail",
			dummyUser.Email,
			mock.Anything, // from
			dummyUser.Username,
			mock.Anything, // verification link
			dummyUser.VerificationTokens,
		).
		Return(nil).
		Run(func(_ mock.Arguments) {
			close(done)
		})

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.Signup(c)

	<-done

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var respBody map[string]bool
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.True(t, respBody["success"])

	mockUserService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestSignup_InvalidJSON(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)
	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	invalidJSON := []byte(`{invalid-json}`) // malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(invalidJSON))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.Signup(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid JSON")
}

func TestSignup_UserCreationFails(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserSignupRequest{
		Fullname: "Jane Doe",
		Username: "janedoe",
		Email:    "jane@example.com",
		Password: "password123",
	}
	bodyBytes, _ := json.Marshal(input)

	mockUserService.
		On("CreateUser", input.Fullname, input.Username, input.Email, input.Password).
		Return(nil, fmt.Errorf("user exists")) // updated: no second string value

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.Signup(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.NotEqual(t, http.StatusCreated, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "Could not create user", respBody["error"])

	mockUserService.AssertExpectations(t)
	mockEmailService.AssertNotCalled(
		t,
		"SendVerifyAccountEmail",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	)
}

func TestSignup_EmptyRequestBody(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.Signup(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid JSON")
}

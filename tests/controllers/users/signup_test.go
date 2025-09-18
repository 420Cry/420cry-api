package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controller "cry-api/app/controllers/users"
	UserModel "cry-api/app/models"
	SignUpError "cry-api/app/types/errors"
	UserTypes "cry-api/app/types/users"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignup_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)
	mockUserTokenService := new(testmocks.MockUserTokenService)

	userController := &controller.UserController{
		UserService:      mockUserService,
		EmailService:     mockEmailService,
		UserTokenService: mockUserTokenService,
	}

	input := UserTypes.IUserSignupRequest{
		Fullname: "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "securepassword",
	}
	bodyBytes, _ := json.Marshal(input)

	dummyUser := &UserModel.User{
		ID:       1,
		Email:    input.Email,
		Username: input.Username,
	}

	done := make(chan struct{})

	// Mock: CreateUser
	mockUserService.
		On("CreateUser", input.Fullname, input.Username, input.Email, input.Password).
		Return(dummyUser, nil)

	// Mock: SaveUserToken for both link and OTP tokens
	mockUserTokenService.
		On("Save", mock.AnythingOfType("*models.UserToken")).
		Return(nil)

	// Mock: email service
	mockEmailService.
		On("SendVerifyAccountEmail",
			dummyUser.Email,
			mock.Anything,
			dummyUser.Username,
			mock.Anything,
			mock.Anything,
		).
		Return(nil).
		Run(func(_ mock.Arguments) {
			close(done) // signal email sent
		})

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()
	c := TestUtils.GetGinContext(w, req)

	userController.Signup(c)

	// wait for the async email goroutine
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected SendVerifyAccountEmail to be called, but it wasn’t")
	}

	// verify response
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var respBody map[string]bool
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.True(t, respBody["success"])

	// assert all expectations
	mockUserService.AssertExpectations(t)
	mockUserTokenService.AssertExpectations(t)
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

func TestSignup_UserConflict(t *testing.T) {
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

	// Return ErrUserConflict
	mockUserService.
		On("CreateUser", input.Fullname, input.Username, input.Email, input.Password).
		Return(nil, SignUpError.ErrUserConflict)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.Signup(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusConflict, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "User already exists", respBody["error"])

	mockEmailService.AssertNotCalled(t, "SendVerifyAccountEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
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

	// Return generic error
	mockUserService.
		On("CreateUser", input.Fullname, input.Username, input.Email, input.Password).
		Return(nil, fmt.Errorf("db error"))

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

	mockEmailService.AssertNotCalled(t, "SendVerifyAccountEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
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

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

func TestHandleResetPasswordRequest_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IVerificationResetPasswordRequest{
		Email: "test@example.com",
	}
	bodyBytes, _ := json.Marshal(input)

	dummyUser := &UserModel.User{
		Email:      input.Email,
		Username:   "testuser",
		IsVerified: true,
	}

	mockUserService.
		On("FindUserByEmail", input.Email).
		Return(dummyUser, nil)

	mockUserService.
		On("UpdateUser", mock.AnythingOfType("*models.User")).
		Return(nil)

	done := make(chan struct{})

	mockEmailService.
		On("SendResetPasswordEmail",
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

	req := httptest.NewRequest(http.MethodPost, "/reset-password-request", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.HandleResetPasswordRequest(c)

	<-done // wait for SendResetPasswordEmail call

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var respBody map[string]bool
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.True(t, respBody["success"])

	mockUserService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestHandleResetPasswordRequest_InvalidJSON(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	invalidJSON := []byte(`{invalid-json}`)

	req := httptest.NewRequest(http.MethodPost, "/reset-password-request", bytes.NewReader(invalidJSON))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.HandleResetPasswordRequest(c)

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

func TestHandleResetPasswordRequest_UserNotFound(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IVerificationResetPasswordRequest{
		Email: "notfound@example.com",
	}
	bodyBytes, _ := json.Marshal(input)

	mockUserService.
		On("FindUserByEmail", input.Email).
		Return(nil, nil) // user not found

	req := httptest.NewRequest(http.MethodPost, "/reset-password-request", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.HandleResetPasswordRequest(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "user not found", respBody["error"])
}

func TestHandleResetPasswordRequest_UserNotVerified(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IVerificationResetPasswordRequest{
		Email: "unverified@example.com",
	}
	bodyBytes, _ := json.Marshal(input)

	dummyUser := &UserModel.User{
		Email:      input.Email,
		Username:   "unverifieduser",
		IsVerified: false,
	}

	mockUserService.
		On("FindUserByEmail", input.Email).
		Return(dummyUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/reset-password-request", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.HandleResetPasswordRequest(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusForbidden, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "user not verified", respBody["error"])
}

func TestHandleResetPasswordRequest_InternalErrorFindingUser(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockEmailService := new(testmocks.MockEmailService)

	userController := &controller.UserController{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IVerificationResetPasswordRequest{
		Email: "error@example.com",
	}
	bodyBytes, _ := json.Marshal(input)

	mockUserService.
		On("FindUserByEmail", input.Email).
		Return(nil, fmt.Errorf("db failure"))

	req := httptest.NewRequest(http.MethodPost, "/reset-password-request", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.HandleResetPasswordRequest(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Equal(t, "internal error", respBody["error"])
}

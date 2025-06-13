// Package user_routes_test provides tests for user routes.
package user_routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	users "cry-api/app/api/routes/users"
	UserDomain "cry-api/app/domain/users"
	UserTypes "cry-api/app/types/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignup_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)

	handler := &users.Handler{
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

	dummyUser := &UserDomain.User{
		Email:              input.Email,
		Username:           input.Username,
		VerificationTokens: "verify123",
	}

	// Channel to signal SendVerifyAccountEmail was called
	done := make(chan struct{})

	// Setup mock expectations
	mockUserService.
		On("CreateUser", input.Fullname, input.Username, input.Email, input.Password).
		Return(dummyUser, "token123", nil)

	mockEmailService.
		On("SendVerifyAccountEmail",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(nil).
		Run(func(_ mock.Arguments) {
			close(done)
		})

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	// Wait for email sending to be called (since it's async)
	<-done

	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
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
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	handler := &users.Handler{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	invalidJSON := []byte(`{invalid-json}`) // malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(invalidJSON))
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	res := w.Result()
	if err := res.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid JSON")
}

func TestSignup_UserCreationFails(t *testing.T) {
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	handler := &users.Handler{
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

	// Simulate error from CreateUser
	mockUserService.
		On("CreateUser", input.Fullname, input.Username, input.Email, input.Password).
		Return((*UserDomain.User)(nil), "", fmt.Errorf("user exists"))

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	res := w.Result()
	if err := res.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}

	assert.NotEqual(t, http.StatusCreated, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "user exists")

	mockUserService.AssertExpectations(t)
	// Email service should NOT be called on failure
	mockEmailService.AssertNotCalled(t, "SendVerifyAccountEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestSignup_EmptyRequestBody(t *testing.T) {
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)
	handler := &users.Handler{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader([]byte{})) // empty body
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	res := w.Result()
	if err := res.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}

	// Should fail with invalid JSON (empty body)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid JSON")
}

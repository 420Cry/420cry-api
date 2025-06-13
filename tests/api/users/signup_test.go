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

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(fullname, username, email, password string) (*UserDomain.User, string, error) {
	args := m.Called(fullname, username, email, password)
	return args.Get(0).(*UserDomain.User), args.String(1), args.Error(2)
}

func (m *MockUserService) CheckUserByBothTokens(token string, verificationToken string) (*UserDomain.User, error) {
	args := m.Called(token, verificationToken)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

func (m *MockUserService) CheckEmailVerificationToken(token string) (*UserDomain.User, error) {
	args := m.Called(token)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

func (m *MockUserService) CheckAccountVerificationToken(token string) (*UserDomain.User, error) {
	args := m.Called(token)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

func (m *MockUserService) AuthenticateUser(username string, password string) (*UserDomain.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

func (m *MockUserService) VerifyUserWithTokens(token string, verificationToken string) (*UserDomain.User, error) {
	args := m.Called(token, verificationToken)
	return args.Get(0).(*UserDomain.User), args.Error(1)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendVerifyAccountEmail(to, from, username, link, token string) error {
	args := m.Called(to, from, username, link, token)
	return args.Error(0)
}
func TestSignup_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)

	handler := &users.Handler{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	// Dummy input
	input := UserTypes.UserSignupRequest{
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
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.Signup(w, req)

	res := w.Result()
	if err := res.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}

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

	input := UserTypes.UserSignupRequest{
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

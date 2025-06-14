package user_routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	users "cry-api/app/api/routes/users"
	UserDomain "cry-api/app/domain/users"
	UserTypes "cry-api/app/types/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignIn_Success(t *testing.T) {
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService) // unused here but needed for Handler struct

	handler := &users.Handler{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	input := UserTypes.IUserSigninRequest{
		Username: "johndoe",
		Password: "securepassword",
	}

	bodyBytes, _ := json.Marshal(input)

	dummyUser := &UserDomain.User{
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

	handler.SignIn(w, req)

	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
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
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)

	handler := &users.Handler{
		UserService:  mockUserService,
		EmailService: mockEmailService,
	}

	invalidJSON := []byte(`{invalid-json}`) // malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(invalidJSON))
	w := httptest.NewRecorder()

	handler.SignIn(w, req)

	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid JSON")

	mockUserService.AssertNotCalled(t, "AuthenticateUser", mock.Anything, mock.Anything)
}

func TestSignIn_AuthenticationFails(t *testing.T) {
	mockUserService := new(MockUserService)
	mockEmailService := new(MockEmailService)

	handler := &users.Handler{
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
		Return((*UserDomain.User)(nil), assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.SignIn(w, req)

	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Fatalf("failed to close response body: %v", err)
		}
	}()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	var respBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&respBody)
	assert.NoError(t, err)
	assert.Contains(t, respBody["error"], "Invalid email or password")

	mockUserService.AssertExpectations(t)
}

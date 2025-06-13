package user_routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	users "cry-api/app/api/routes/users"
	UserDomain "cry-api/app/domain/users"

	"github.com/stretchr/testify/assert"
)

func TestVerifyAccountToken_Success(t *testing.T) {
	mockUserService := new(MockUserService)

	handler := &users.Handler{
		UserService: mockUserService,
	}

	token := "valid-token-123"
	user := &UserDomain.User{
		Token:                      token,
		VerificationTokenCreatedAt: time.Now(),
	}

	mockUserService.On("CheckAccountVerificationToken", token).Return(user, nil)

	bodyBytes, _ := json.Marshal(map[string]string{"token": token})

	req := httptest.NewRequest(http.MethodPost, "/verify-account-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.VerifyAccountToken(w, req)

	res := w.Result()
	if err := res.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp map[string]bool
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.True(t, resp["valid"])

	mockUserService.AssertExpectations(t)
}

func TestVerifyAccountToken_InvalidJSON(t *testing.T) {
	mockUserService := new(MockUserService)

	handler := &users.Handler{
		UserService: mockUserService,
	}

	req := httptest.NewRequest(http.MethodPost, "/verify-account-token", bytes.NewReader([]byte(`{invalid-json}`)))
	w := httptest.NewRecorder()

	handler.VerifyAccountToken(w, req)

	res := w.Result()
	if err := res.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "Invalid request body")
}

func TestVerifyAccountToken_UserNotFound(t *testing.T) {
	mockUserService := new(MockUserService)

	handler := &users.Handler{
		UserService: mockUserService,
	}

	token := "nonexistent-token"
	mockUserService.On("CheckAccountVerificationToken", token).Return((*UserDomain.User)(nil), assert.AnError)

	bodyBytes, _ := json.Marshal(map[string]string{"token": token})

	req := httptest.NewRequest(http.MethodPost, "/verify-account-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.VerifyAccountToken(w, req)

	res := w.Result()
	if err := res.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "Token is invalid or expired")

	mockUserService.AssertExpectations(t)
}

func TestVerifyAccountToken_TokenMismatchOrExpired(t *testing.T) {
	mockUserService := new(MockUserService)

	handler := &users.Handler{
		UserService: mockUserService,
	}

	token := "valid-token-123"
	// Token in user is different or created too far in the past
	user := &UserDomain.User{
		Token:                      "different-token",
		VerificationTokenCreatedAt: time.Now().Add(-25 * time.Hour),
	}

	mockUserService.On("CheckAccountVerificationToken", token).Return(user, nil)

	bodyBytes, _ := json.Marshal(map[string]string{"token": token})

	req := httptest.NewRequest(http.MethodPost, "/verify-account-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.VerifyAccountToken(w, req)

	res := w.Result()
	if err := res.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "Token is invalid or expired")

	mockUserService.AssertExpectations(t)
}

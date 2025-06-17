package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	users "cry-api/app/api/routes/users"
	UserDomain "cry-api/app/domain/users"
	UserTypes "cry-api/app/types/users"
	TestUtils "cry-api/app/utils/tests"

	"github.com/stretchr/testify/assert"
)

func TestVerifyEmailToken_Success(t *testing.T) {
	mockUserService := new(MockUserService)

	handler := &users.Handler{
		UserService: mockUserService,
	}

	reqBody := UserTypes.IVerificationTokenCheckRequest{
		UserToken:   "user-token-123",
		VerifyToken: "verify-token-abc",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	dummyUser := &UserDomain.User{
		IsVerified: true,
	}

	mockUserService.
		On("VerifyUserWithTokens", reqBody.UserToken, reqBody.VerifyToken).
		Return(dummyUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-email-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	handler.VerifyEmailToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var resp map[string]bool
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.True(t, resp["verified"])

	mockUserService.AssertExpectations(t)
}

func TestVerifyEmailToken_InvalidJSON(t *testing.T) {
	mockUserService := new(MockUserService)

	handler := &users.Handler{
		UserService: mockUserService,
	}

	req := httptest.NewRequest(http.MethodPost, "/verify-email-token", bytes.NewReader([]byte(`{invalid-json}`)))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	handler.VerifyEmailToken(c)

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
	mockUserService := new(MockUserService)

	handler := &users.Handler{
		UserService: mockUserService,
	}

	reqBody := UserTypes.IVerificationTokenCheckRequest{
		UserToken:   "user-token-123",
		VerifyToken: "verify-token-abc",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	mockUserService.
		On("VerifyUserWithTokens", reqBody.UserToken, reqBody.VerifyToken).
		Return((*UserDomain.User)(nil), assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/verify-email-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	handler.VerifyEmailToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], assert.AnError.Error())

	mockUserService.AssertExpectations(t)
}

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controller "cry-api/app/controllers/users"
	UserModel "cry-api/app/models"
	UserTypes "cry-api/app/types/users"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVerifyEmailToken_Success(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)

	userController := &controller.UserController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	reqBody := UserTypes.IVerificationTokenCheckRequest{
		UserToken:   "user-token-123",
		VerifyToken: "verify-token-abc",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// 1️⃣ Mock the long-link token
	userTokenObj := &UserModel.UserToken{
		Token:     reqBody.UserToken,
		UserID:    42,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockUserService.
		On("FindUserTokenByValueAndPurpose", reqBody.UserToken, mock.Anything).
		Return(userTokenObj, nil)

	// 2️⃣ Mock the OTP token
	otpTokenObj := &UserModel.UserToken{
		Token:     reqBody.VerifyToken,
		UserID:    42,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	mockUserTokenService.
		On("FindLatestValidToken", userTokenObj.UserID, mock.Anything).
		Return(otpTokenObj, nil)

	// 3️⃣ Mock ConsumeToken calls
	mockUserTokenService.On("ConsumeToken", userTokenObj.UserID, userTokenObj.Token, mock.Anything).Return(nil)
	mockUserTokenService.On("ConsumeToken", otpTokenObj.UserID, otpTokenObj.Token, mock.Anything).Return(nil)

	// 4️⃣ Mock updating user
	dummyUser := &UserModel.User{ID: userTokenObj.UserID}
	mockUserService.On("FindUserByID", userTokenObj.UserID).Return(dummyUser, nil)
	mockUserService.On("UpdateUser", dummyUser).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-email-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyEmailToken(c)

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
	mockUserTokenService.AssertExpectations(t)
}

func TestVerifyEmailToken_InvalidJSON(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	userController := &controller.UserController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	req := httptest.NewRequest(http.MethodPost, "/verify-email-token", bytes.NewReader([]byte(`{invalid-json}`)))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyEmailToken(c)

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
	mockUserService := new(testmocks.MockUserService)
	mockUserTokenService := new(testmocks.MockUserTokenService)
	userController := &controller.UserController{
		UserService:      mockUserService,
		UserTokenService: mockUserTokenService,
	}

	reqBody := UserTypes.IVerificationTokenCheckRequest{
		UserToken:   "user-token-123",
		VerifyToken: "verify-token-abc",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Simulate missing long-link token
	mockUserService.
		On("FindUserTokenByValueAndPurpose", reqBody.UserToken, mock.Anything).
		Return(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/verify-email-token", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	c := TestUtils.GetGinContext(w, req)
	userController.VerifyEmailToken(c)

	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["error"], "invalid or expired account verification link")

	mockUserService.AssertExpectations(t)
	mockUserTokenService.AssertExpectations(t)
}

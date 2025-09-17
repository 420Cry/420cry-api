package two_factor_controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "cry-api/app/controllers/2fa"
	UserModel "cry-api/app/models"
	JWT "cry-api/app/services/jwt"
	testmocks "cry-api/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVerifyOTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(testmocks.MockUserService)
	mockAuthService := new(testmocks.MockAuthService)

	controller := &controller.TwoFactorController{
		UserService: mockUserService,
		AuthService: mockAuthService,
	}

	// Helper to perform requests
	performRequest := func(body any) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBytes, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest("POST", "/verify-otp", bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.VerifyOTP(c)
		return w
	}

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("POST", "/verify-otp", bytes.NewReader([]byte("{invalid-json")))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.VerifyOTP(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid request"}`, w.Body.String())
	})

	t.Run("Missing UserUUID", func(t *testing.T) {
		resp := performRequest(map[string]string{"otp": "123456"})
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"User UUID is required"}`, resp.Body.String())
	})

	t.Run("Missing OTP", func(t *testing.T) {
		resp := performRequest(map[string]string{"userUUID": "user-123"})
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"OTP is required for verification"}`, resp.Body.String())
	})

	t.Run("User retrieval error", func(t *testing.T) {
		mockUserService.On("GetUserByUUID", "user-123").Return(nil, errors.New("db error")).Once()

		resp := performRequest(map[string]string{"userUUID": "user-123", "otp": "123456"})
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"Failed to get user"}`, resp.Body.String())

		mockUserService.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		mockUserService.On("GetUserByUUID", "user-123").Return(nil, nil).Once()

		resp := performRequest(map[string]string{"userUUID": "user-123", "otp": "123456"})
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error":"User not found"}`, resp.Body.String())

		mockUserService.AssertExpectations(t)
	})

	t.Run("OTP verification failure", func(t *testing.T) {
		user := &UserModel.User{
			UUID:         "user-123",
			Email:        "user@example.com",
			TwoFASecret:  stringPtr("secret123"),
			TwoFAEnabled: false,
		}

		mockUserService.On("GetUserByUUID", "user-123").Return(user, nil).Once()
		mockAuthService.On("VerifyOTP", "secret123", "wrong-otp").Return(false, nil).Once()

		resp := performRequest(map[string]string{"userUUID": "user-123", "otp": "wrong-otp"})
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.JSONEq(t, `{"error":"Invalid OTP"}`, resp.Body.String())

		mockUserService.AssertExpectations(t)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Successful OTP verification and 2FA enable", func(t *testing.T) {
		user := &UserModel.User{
			UUID:         "user-123",
			Email:        "user@example.com",
			TwoFASecret:  stringPtr("secret123"),
			TwoFAEnabled: false,
		}

		mockUserService.On("GetUserByUUID", "user-123").Return(user, nil).Once()
		mockAuthService.On("VerifyOTP", "secret123", "valid-otp").Return(true, nil).Once()
		mockUserService.On("UpdateUser", mock.Anything).Return(nil).Once()

		// Mock JWT.GenerateJWT globally (override it)
		oldGenerateJWT := JWT.GenerateJWT
		defer func() { JWT.GenerateJWT = oldGenerateJWT }()
		JWT.GenerateJWT = func(_, _ string, _, _ bool) (string, error) {
			return "mocked.jwt.token", nil
		}

		resp := performRequest(map[string]string{"userUUID": "user-123", "otp": "valid-otp"})
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.JSONEq(t, `{"jwt":"mocked.jwt.token"}`, resp.Body.String())

		mockUserService.AssertExpectations(t)
		mockAuthService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Failed to update user when enabling 2FA", func(t *testing.T) {
		user := &UserModel.User{
			UUID:         "user-123",
			Email:        "user@example.com",
			TwoFASecret:  stringPtr("secret123"),
			TwoFAEnabled: false,
		}

		mockUserService.On("GetUserByUUID", "user-123").Return(user, nil).Once()
		mockAuthService.On("VerifyOTP", "secret123", "valid-otp").Return(true, nil).Once()
		mockUserService.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(errors.New("update failed"))

		resp := performRequest(map[string]string{"userUUID": "user-123", "otp": "valid-otp"})
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"Failed to enable 2FA"}`, resp.Body.String())

		mockUserService.AssertExpectations(t)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("JWT generation failure", func(t *testing.T) {
		user := &UserModel.User{
			UUID:         "user-123",
			Email:        "user@example.com",
			TwoFASecret:  stringPtr("secret123"),
			TwoFAEnabled: true,
		}

		mockUserService.On("GetUserByUUID", "user-123").Return(user, nil).Once()
		mockAuthService.On("VerifyOTP", "secret123", "valid-otp").Return(true, nil).Once()

		// Mock JWT.GenerateJWT globally
		oldGenerateJWT := JWT.GenerateJWT
		defer func() { JWT.GenerateJWT = oldGenerateJWT }()
		JWT.GenerateJWT = func(_, _ string, _, _ bool) (string, error) {
			return "", errors.New("jwt error")
		}

		resp := performRequest(map[string]string{"userUUID": "user-123", "otp": "valid-otp"})
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"Failed to generate JWT"}`, resp.Body.String())

		mockUserService.AssertExpectations(t)
		mockAuthService.AssertExpectations(t)
	})
}

// helper to create pointer to string
func stringPtr(s string) *string {
	return &s
}

package tests

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
	types "cry-api/app/types/2fa"
	testmocks "cry-api/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVerifySetUpOTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(testmocks.MockUserService)
	mockAuthService := new(testmocks.MockAuthService)

	ctrl := &controller.TwoFactorController{
		UserService: mockUserService,
		AuthService: mockAuthService,
	}

	performRequest := func(body any) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBytes, _ := json.Marshal(body)
		c.Request, _ = http.NewRequest("POST", "/verify-setup-otp", bytes.NewReader(jsonBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.VerifySetUpOTP(c)
		return w
	}

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("POST", "/verify-setup-otp", bytes.NewReader([]byte("{invalid-json")))
		c.Request.Header.Set("Content-Type", "application/json")

		ctrl.VerifySetUpOTP(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid request"}`, w.Body.String())
	})

	t.Run("Missing UserUUID", func(t *testing.T) {
		resp := performRequest(types.ITwoFactorSetupRequest{
			OTP:    stringPtr("123456"),
			Secret: stringPtr("secret123"),
		})
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"User UUID is required"}`, resp.Body.String())
	})

	t.Run("Missing OTP", func(t *testing.T) {
		resp := performRequest(types.ITwoFactorSetupRequest{
			UserUUID: "user-123",
			Secret:   stringPtr("secret123"),
		})
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"OTP is required for verification"}`, resp.Body.String())
	})

	t.Run("Missing Secret", func(t *testing.T) {
		resp := performRequest(types.ITwoFactorSetupRequest{
			UserUUID: "user-123",
			OTP:      stringPtr("123456"),
		})
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, `{"error":"TOTP secret is required"}`, resp.Body.String())
	})

	t.Run("OTP verification failure", func(t *testing.T) {
		mockAuthService.On("VerifyOTP", "secret123", "wrong-otp").Return(false, nil).Once()

		resp := performRequest(types.ITwoFactorSetupRequest{
			UserUUID: "user-123",
			OTP:      stringPtr("wrong-otp"),
			Secret:   stringPtr("secret123"),
		})
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.JSONEq(t, `{"error":"Invalid OTP"}`, resp.Body.String())

		mockAuthService.AssertExpectations(t)
	})

	t.Run("OTP verification error", func(t *testing.T) {
		mockAuthService.On("VerifyOTP", "secret123", "some-otp").Return(false, errors.New("otp error")).Once()

		resp := performRequest(types.ITwoFactorSetupRequest{
			UserUUID: "user-123",
			OTP:      stringPtr("some-otp"),
			Secret:   stringPtr("secret123"),
		})
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.JSONEq(t, `{"error":"Invalid OTP"}`, resp.Body.String())

		mockAuthService.AssertExpectations(t)
	})

	t.Run("User retrieval error", func(t *testing.T) {
		mockAuthService.On("VerifyOTP", "secret123", "valid-otp").Return(true, nil).Once()
		mockUserService.On("GetUserByUUID", "user-123").Return(nil, errors.New("db error")).Once()

		resp := performRequest(types.ITwoFactorSetupRequest{
			UserUUID: "user-123",
			OTP:      stringPtr("valid-otp"),
			Secret:   stringPtr("secret123"),
		})
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"Failed to get user"}`, resp.Body.String())

		mockAuthService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		mockAuthService.On("VerifyOTP", "secret123", "valid-otp").Return(true, nil).Once()
		mockUserService.On("GetUserByUUID", "user-123").Return(nil, nil).Once()

		resp := performRequest(types.ITwoFactorSetupRequest{
			UserUUID: "user-123",
			OTP:      stringPtr("valid-otp"),
			Secret:   stringPtr("secret123"),
		})
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, `{"error":"User not found"}`, resp.Body.String())

		mockAuthService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Failed to update user when enabling 2FA", func(t *testing.T) {
		user := &UserModel.User{
			UUID:         "user-123",
			Email:        "user@example.com",
			TwoFAEnabled: false,
		}

		mockAuthService.On("VerifyOTP", "secret123", "valid-otp").Return(true, nil).Once()
		mockUserService.On("GetUserByUUID", "user-123").Return(user, nil).Once()
		mockUserService.On("UpdateUser", mock.Anything).Return(errors.New("update failed")).Once()

		resp := performRequest(types.ITwoFactorSetupRequest{
			UserUUID: "user-123",
			OTP:      stringPtr("valid-otp"),
			Secret:   stringPtr("secret123"),
		})
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"Failed to enable 2FA"}`, resp.Body.String())

		mockAuthService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})

	t.Run("JWT generation failure", func(t *testing.T) {
		user := &UserModel.User{
			UUID:         "user-123",
			Email:        "user@example.com",
			TwoFAEnabled: true,
		}

		mockAuthService.On("VerifyOTP", "secret123", "valid-otp").Return(true, nil).Once()
		mockUserService.On("GetUserByUUID", "user-123").Return(user, nil).Once()

		oldGenerateJWT := JWT.GenerateJWT
		defer func() { JWT.GenerateJWT = oldGenerateJWT }()
		JWT.GenerateJWT = func(_, _ string, _, _ bool) (string, error) {
			return "", errors.New("jwt error")
		}

		resp := performRequest(types.ITwoFactorSetupRequest{
			UserUUID: "user-123",
			OTP:      stringPtr("valid-otp"),
			Secret:   stringPtr("secret123"),
		})

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.JSONEq(t, `{"error":"Failed to generate JWT"}`, resp.Body.String())

		mockAuthService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Successful verification and 2FA enable", func(t *testing.T) {
		user := &UserModel.User{
			UUID:         "user-123",
			Email:        "user@example.com",
			Fullname:     "Test User",
			Username:     "testuser",
			TwoFAEnabled: false,
		}

		mockAuthService.On("VerifyOTP", "secret123", "valid-otp").Return(true, nil).Once()
		mockUserService.On("GetUserByUUID", "user-123").Return(user, nil).Once()
		mockUserService.On("UpdateUser", mock.Anything).Return(nil).Once()

		oldGenerateJWT := JWT.GenerateJWT
		defer func() { JWT.GenerateJWT = oldGenerateJWT }()
		JWT.GenerateJWT = func(_, _ string, _, _ bool) (string, error) {
			return "mocked.jwt.token", nil
		}

		resp := performRequest(types.ITwoFactorSetupRequest{
			UserUUID: "user-123",
			OTP:      stringPtr("valid-otp"),
			Secret:   stringPtr("secret123"),
		})

		assert.Equal(t, http.StatusOK, resp.Code)
		expectedBody := `{
			"jwt": "mocked.jwt.token",
			"user": {
				"uuid": "user-123",
				"fullname": "Test User",
				"email": "user@example.com",
				"username": "testuser",
				"twoFAEnabled": true
			}
		}`
		assert.JSONEq(t, expectedBody, resp.Body.String())

		mockAuthService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})
}

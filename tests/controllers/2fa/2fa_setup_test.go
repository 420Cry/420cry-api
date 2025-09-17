package two_factor_controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	controller "cry-api/app/controllers/2fa"
	UserModel "cry-api/app/models"
	TwoFactorType "cry-api/app/types/2fa"
	TestUtils "cry-api/app/utils/tests"
	testmocks "cry-api/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTwoFactorController_Setup(t *testing.T) {
	mockUserService := new(testmocks.MockUserService)
	mockAuthService := new(testmocks.MockAuthService)
	mockTwoFactorService := new(testmocks.MockTwoFactorService)

	twoFactorController := &controller.TwoFactorController{
		UserService:      mockUserService,
		AuthService:      mockAuthService,
		TwoFactorService: mockTwoFactorService,
	}

	makeRequest := func(body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/2fa/setup", bytes.NewReader(b))
		w := httptest.NewRecorder()
		c := TestUtils.GetGinContext(w, req)
		return c, w
	}

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/2fa/setup", bytes.NewReader([]byte("{invalid-json")))
		w := httptest.NewRecorder()
		c := TestUtils.GetGinContext(w, req)

		twoFactorController.Setup(c)

		resp := w.Result()
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Missing UserUUID", func(t *testing.T) {
		c, w := makeRequest(TwoFactorType.ITwoFactorSetupRequest{UserUUID: ""})
		twoFactorController.Setup(c)

		resp := w.Result()
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userUUID := "missing-uuid"
		mockUserService.On("GetUserByUUID", userUUID).Return(nil, nil).Once()

		c, w := makeRequest(TwoFactorType.ITwoFactorSetupRequest{UserUUID: userUUID})
		twoFactorController.Setup(c)

		resp := w.Result()
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockUserService.AssertExpectations(t)
	})

	t.Run("User Has Existing 2FA Secret - Success", func(t *testing.T) {
		userUUID := "user-uuid"
		secret := "existing-secret"
		user := &UserModel.User{
			Email:       "user@example.com",
			TwoFASecret: &secret,
		}

		mockUserService.On("GetUserByUUID", userUUID).Return(user, nil).Once()
		mockTwoFactorService.On("GenerateOtpauthURL", user.Email, secret).Return("otpauth://mockurl").Once()
		mockTwoFactorService.On("GenerateQRCodeBase64", "otpauth://mockurl").Return("mockQRcodeBase64", nil).Once()

		c, w := makeRequest(TwoFactorType.ITwoFactorSetupRequest{UserUUID: userUUID})
		twoFactorController.Setup(c)

		resp := w.Result()
		defer func() {
			_ = resp.Body.Close()
		}()

		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Logf("Response status: %d, body: %s", resp.StatusCode, string(bodyBytes))

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody map[string]string
		err := json.Unmarshal(bodyBytes, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, secret, respBody["secret"])
		assert.Equal(t, "mockQRcodeBase64", respBody["qrCode"])

		mockUserService.AssertExpectations(t)
		mockTwoFactorService.AssertExpectations(t)
	})

	t.Run("User Has No 2FA Secret - Generate New Secret Success", func(t *testing.T) {
		userUUID := "new-secret-user"
		user := &UserModel.User{
			Email: "newuser@example.com",
		}
		secret := "new-secret"
		otpauthURL := "otpauth://newmockurl"

		mockUserService.On("GetUserByUUID", userUUID).Return(user, nil).Once()
		mockTwoFactorService.On("GenerateTOTP", user.Email).Return(secret, otpauthURL, nil).Once()
		mockUserService.On("UpdateUser", mock.MatchedBy(func(u *UserModel.User) bool {
			return u.TwoFASecret != nil && *u.TwoFASecret == secret
		})).Return(nil).Once()
		mockTwoFactorService.On("GenerateQRCodeBase64", otpauthURL).Return("newMockQRcodeBase64", nil).Once()

		c, w := makeRequest(TwoFactorType.ITwoFactorSetupRequest{UserUUID: userUUID})
		twoFactorController.Setup(c)

		resp := w.Result()
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody map[string]string
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, secret, respBody["secret"])
		assert.Equal(t, "newMockQRcodeBase64", respBody["qrCode"])

		mockUserService.AssertExpectations(t)
		mockTwoFactorService.AssertExpectations(t)
	})

	t.Run("GenerateTOTP Failure", func(t *testing.T) {
		userUUID := "fail-totp"
		user := &UserModel.User{Email: "fail@example.com"}

		mockUserService.On("GetUserByUUID", userUUID).Return(user, nil).Once()
		mockTwoFactorService.On("GenerateTOTP", user.Email).Return("", "", errors.New("fail generate totp")).Once()

		c, w := makeRequest(TwoFactorType.ITwoFactorSetupRequest{UserUUID: userUUID})
		twoFactorController.Setup(c)

		resp := w.Result()
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockUserService.AssertExpectations(t)
		mockTwoFactorService.AssertExpectations(t)
	})

	t.Run("UpdateUser Failure", func(t *testing.T) {
		userUUID := "fail-update"
		user := &UserModel.User{Email: "failupdate@example.com"}
		secret := "some-secret"
		otpauthURL := "otpauth://someurl"

		mockUserService.On("GetUserByUUID", userUUID).Return(user, nil).Once()
		mockTwoFactorService.On("GenerateTOTP", user.Email).Return(secret, otpauthURL, nil).Once()
		mockUserService.On("UpdateUser", mock.Anything).Return(errors.New("db update error")).Once()

		c, w := makeRequest(TwoFactorType.ITwoFactorSetupRequest{UserUUID: userUUID})
		twoFactorController.Setup(c)

		resp := w.Result()
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockUserService.AssertExpectations(t)
		mockTwoFactorService.AssertExpectations(t)
	})

	t.Run("GenerateQRCodeBase64 Failure", func(t *testing.T) {
		userUUID := "fail-qrcode"
		user := &UserModel.User{Email: "failqrcode@example.com"}
		secret := "secret-for-qrcode"
		otpauthURL := "otpauth://failqrcode"

		mockUserService.On("GetUserByUUID", userUUID).Return(user, nil).Once()
		mockTwoFactorService.On("GenerateTOTP", user.Email).Return(secret, otpauthURL, nil).Once()
		mockUserService.On("UpdateUser", mock.Anything).Return(nil).Once()
		mockTwoFactorService.On("GenerateQRCodeBase64", otpauthURL).Return("", errors.New("qrcode error")).Once()

		c, w := makeRequest(TwoFactorType.ITwoFactorSetupRequest{UserUUID: userUUID})
		twoFactorController.Setup(c)

		resp := w.Result()
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockUserService.AssertExpectations(t)
		mockTwoFactorService.AssertExpectations(t)
	})
}

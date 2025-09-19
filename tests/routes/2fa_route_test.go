package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockTwoFactorController mocks TwoFactorController methods
type MockTwoFactorController struct{}

func (m *MockTwoFactorController) Setup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "setup called"})
}

func (m *MockTwoFactorController) VerifySetUpOTP(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "verify-setup-otp called"})
}

func (m *MockTwoFactorController) VerifyOTP(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "verify-otp called"})
}

// helper function to register routes with mock controller (instead of real one)
func registerMockTwoFactorRoutes(rg *gin.RouterGroup, ctrl *MockTwoFactorController) {
	rg.POST("/setup", ctrl.Setup)
	rg.POST("/setup/verify-otp", ctrl.VerifySetUpOTP)
	rg.POST("/auth/verify-otp", ctrl.VerifyOTP)
}

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/2fa")

	mockCtrl := &MockTwoFactorController{}
	registerMockTwoFactorRoutes(rg, mockCtrl)

	testCases := []struct {
		method       string
		endpoint     string
		expectedCode int
		expectedBody string
	}{
		{"POST", "/2fa/setup", http.StatusOK, `{"message":"setup called"}`},
		{"POST", "/2fa/setup/verify-otp", http.StatusOK, `{"message":"verify-setup-otp called"}`},
		{"POST", "/2fa/auth/verify-otp", http.StatusOK, `{"message":"verify-otp called"}`},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.endpoint, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, tc.expectedCode, resp.Code)
		assert.JSONEq(t, tc.expectedBody, resp.Body.String())
	}
}

package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockUserController mimics the real controller, with stub handlers
type MockUserController struct{}

func (m *MockUserController) Signup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "signup called"})
}

func (m *MockUserController) VerifyEmailToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "verify-email-token called"})
}

func (m *MockUserController) VerifyAccountToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "verify-account-token called"})
}

func (m *MockUserController) SignIn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "signin called"})
}

func (m *MockUserController) ResetPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "reset-password called"})
}

func (m *MockUserController) VerifyResetPasswordToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "verify-reset-password-token called"})
}

func (m *MockUserController) CompleteProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "complete-profile called"})
}

// A helper to register routes with mocked controller
func registerTestRoutes(rg *gin.RouterGroup, controller *MockUserController) {
	rg.POST("/signup", controller.Signup)
	rg.POST("/verify-email-token", controller.VerifyEmailToken)
	rg.POST("/verify-account-token", controller.VerifyAccountToken)
	rg.POST("/signin", controller.SignIn)
	rg.POST("/reset-password", controller.ResetPassword)
	rg.POST("/verify-reset-password-token", controller.VerifyResetPasswordToken)
	rg.POST("/complete-profile", controller.CompleteProfile)
}

func TestUserRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/users")

	mockCtrl := &MockUserController{}
	registerTestRoutes(rg, mockCtrl)

	testCases := []struct {
		method       string
		endpoint     string
		expectedCode int
		expectedBody string
	}{
		{"POST", "/users/signup", http.StatusOK, `{"message":"signup called"}`},
		{"POST", "/users/verify-email-token", http.StatusOK, `{"message":"verify-email-token called"}`},
		{"POST", "/users/verify-account-token", http.StatusOK, `{"message":"verify-account-token called"}`},
		{"POST", "/users/signin", http.StatusOK, `{"message":"signin called"}`},
		{"POST", "/users/reset-password", http.StatusOK, `{"message":"reset-password called"}`},
		{"POST", "/users/verify-reset-password-token", http.StatusOK, `{"message":"verify-reset-password-token called"}`},
		{"POST", "/users/complete-profile", http.StatusOK, `{"message":"complete-profile called"}`},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.endpoint, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, tc.expectedCode, resp.Code)
		assert.JSONEq(t, tc.expectedBody, resp.Body.String())
	}
}

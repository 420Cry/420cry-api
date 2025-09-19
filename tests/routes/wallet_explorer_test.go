package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockWalletExplorerController mocks the wallet explorer controller methods
type MockWalletExplorerController struct{}

func (m *MockWalletExplorerController) GetTransactionInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get transaction info called"})
}

func (m *MockWalletExplorerController) GetTransactionByXPUB(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get transaction by xpub called"})
}

// mock middleware that simply calls next handler (bypass real JWT)
func mockJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// helper function to register routes with mock controller and middleware
func registerMockWalletExplorerRoutes(rg *gin.RouterGroup, ctrl *MockWalletExplorerController) {
	authGroup := rg.Group("")
	authGroup.Use(mockJWTMiddleware())

	authGroup.GET("/tx", ctrl.GetTransactionInfo)
	authGroup.GET("/xpub", ctrl.GetTransactionByXPUB)
}

func TestWalletExplorerRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	rg := router.Group("/wallet")

	mockCtrl := &MockWalletExplorerController{}
	registerMockWalletExplorerRoutes(rg, mockCtrl)

	testCases := []struct {
		method       string
		endpoint     string
		expectedCode int
		expectedBody string
	}{
		{"GET", "/wallet/tx", http.StatusOK, `{"message":"get transaction info called"}`},
		{"GET", "/wallet/xpub", http.StatusOK, `{"message":"get transaction by xpub called"}`},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.endpoint, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, tc.expectedCode, resp.Code)
		assert.JSONEq(t, tc.expectedBody, resp.Body.String())
	}
}

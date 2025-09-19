package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cry-api/app/middleware"
	JwtServices "cry-api/app/services/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func generateTestJWT(t *testing.T, twoFAEnabled bool, twoFAVerified bool) string {
	// Set a secret for signing
	err := os.Setenv("JWT_SECRET", "testsecretkey1234567890")
	assert.NoError(t, err)

	token, err := JwtServices.GenerateJWT(
		"test-uuid-1234",
		"user@example.com",
		twoFAEnabled,
		twoFAVerified,
	)
	assert.NoError(t, err)
	return token
}

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := func(c *gin.Context) {
		claims, exists := c.Get("user")
		assert.True(t, exists)
		userClaims := claims.(*JwtServices.Claims)
		assert.Equal(t, "test-uuid-1234", userClaims.UUID)
		assert.Equal(t, "user@example.com", userClaims.Email)
		c.JSON(http.StatusOK, gin.H{"message": "authorized"})
	}

	tests := []struct {
		name            string
		authHeader      string
		expectedCode    int
		expectedBody    string
		expectUserInCtx bool
	}{
		{
			name:         "Missing Authorization header",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Authorization header missing"}`,
		},
		{
			name:         "Malformed Authorization header",
			authHeader:   "Basic abcdefg",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Authorization header format must be Bearer {token}"}`,
		},
		{
			name:         "Invalid token",
			authHeader:   "Bearer invalid.token.here",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"Invalid or expired token"}`,
		},
		{
			name:            "Valid token",
			authHeader:      "Bearer " + generateTestJWT(t, false, false),
			expectedCode:    http.StatusOK,
			expectedBody:    `{"message":"authorized"}`,
			expectUserInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test router with the middleware and test handler
			r := gin.New()
			r.Use(middleware.JWTAuthMiddleware())
			r.GET("/protected", handler)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

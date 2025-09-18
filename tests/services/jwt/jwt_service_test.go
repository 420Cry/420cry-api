package tests

import (
	"os"
	"strings"
	"testing"
	"time"

	services "cry-api/app/services/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	// Setup a test secret for signing tokens
	if err := os.Setenv("JWT_SECRET", "testsecretkey1234567890"); err != nil {
		t.Fatalf("failed to set JWT_SECRET: %v", err)
	}
	services.SetJWTSecret([]byte(os.Getenv("JWT_SECRET")))

	uuid := "123e4567-e89b-12d3-a456-426614174000"
	email := "user@example.com"

	t.Run("Without 2FA enabled", func(t *testing.T) {
		tokenString, err := services.GenerateJWT(uuid, email, false, false)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
		assert.True(t, strings.Count(tokenString, ".") == 2) // basic JWT structure check

		// Parse token to check claims
		token, err := jwt.ParseWithClaims(tokenString, &services.Claims{}, func(_ *jwt.Token) (interface{}, error) {
			return services.GetJWTSecret(), nil
		})
		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(*services.Claims)
		assert.True(t, ok)
		assert.Equal(t, uuid, claims.UUID)
		assert.Equal(t, email, claims.Email)
		assert.False(t, claims.TwoFAEnabled)

		// Expiry should be ~10 minutes from now
		assert.WithinDuration(t, time.Now().Add(10*time.Minute), claims.ExpiresAt.Time, time.Minute)
	})

	t.Run("With 2FA enabled", func(t *testing.T) {
		tokenString, err := services.GenerateJWT(uuid, email, true, true)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
		assert.True(t, strings.Count(tokenString, ".") == 2)

		token, err := jwt.ParseWithClaims(tokenString, &services.Claims{}, func(_ *jwt.Token) (interface{}, error) {
			return services.GetJWTSecret(), nil
		})
		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(*services.Claims)
		assert.True(t, ok)
		assert.Equal(t, uuid, claims.UUID)
		assert.Equal(t, email, claims.Email)
		assert.True(t, claims.TwoFAEnabled)
		assert.True(t, claims.TwoFAVerified)

		// Expiry should be ~7 days from now
		assert.WithinDuration(t, time.Now().Add(7*24*time.Hour), claims.ExpiresAt.Time, time.Minute*5)
	})
}

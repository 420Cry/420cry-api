// Package services provides JWT functions.
package services

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret         []byte
	jwtSecretOnce     sync.Once
	secretInitialized bool
)

// getJWTSecret returns the JWT secret, initializing it on first use
func getJWTSecret() []byte {
	jwtSecretOnce.Do(func() {
		secret := os.Getenv("JWT_SECRET")
		if len(secret) == 0 {
			// In test mode, allow empty secret (tests can set it via SetJWTSecret)
			if os.Getenv("APP_ENV") != "test" {
				log.Fatal("JWT_SECRET is not set; refusing to start")
			}
		} else {
			jwtSecret = []byte(secret)
		}
		secretInitialized = true
	})
	return jwtSecret
}

// Claims defines the custom JWT claims used for authentication,
// embedding standard registered claims along with user-specific fields
// such as UUID and Email.
type Claims struct {
	UUID                    string `json:"uuid"`
	Email                   string `json:"email"`
	TwoFAEnabled            bool   `json:"twoFAEnabled"`
	TwoFASetUpSkippedForNow bool   `json:"twoFASetUpSkippedForNow"`
	TwoFAVerified           bool   `json:"twoFAVerified"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token
var GenerateJWT = func(uuid, email string, twoFAEnabled bool, twoFAVerified bool) (string, error) {
	var expiryDuration time.Duration
	if twoFAEnabled {
		expiryDuration = 7 * 24 * time.Hour
	} else {
		expiryDuration = 10 * time.Minute
	}

	claims := Claims{
		UUID:         uuid,
		Email:        email,
		TwoFAEnabled: twoFAEnabled,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   uuid,
		},
		TwoFAVerified: twoFAVerified,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTSecret())
}

// SetJWTSecret is a setter for the JWT secret (used for testing ONLY)
// This should be called before any JWT operations in tests
func SetJWTSecret(secret []byte) {
	jwtSecret = secret
	secretInitialized = true
}

// GetJWTSecret is a getter for the JWT secret
func GetJWTSecret() []byte {
	if secretInitialized {
		return jwtSecret
	}
	return getJWTSecret()
}

// Package jwtcore provides core JWT functions.
package jwtcore

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// CustomClaims defines the custom JWT claims used for authentication,
// embedding standard registered claims along with user-specific fields
// such as UUID and Email.
type CustomClaims struct {
	UUID         string `json:"uuid"`
	Email        string `json:"email"`
	TwoFAEnabled bool   `json:"twoFAEnabled"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(uuid, email string, twoFAEnabled bool) (string, error) {
	var expiryDuration time.Duration
	if twoFAEnabled {
		expiryDuration = 7 * 24 * time.Hour // 7 days
	} else {
		expiryDuration = 10 * time.Minute // 10 minutes for pre-2FA grace period
	}

	claims := CustomClaims{
		UUID:         uuid,
		Email:        email,
		TwoFAEnabled: twoFAEnabled,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   uuid,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

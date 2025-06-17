// Package services provides JWT functions.
package services

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
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token
func GenerateJWT(uuid, email string) (string, error) {
	claims := CustomClaims{
		UUID:  uuid,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Token valid for 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   uuid,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

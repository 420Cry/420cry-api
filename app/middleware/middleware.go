// Package middleware provides HTTP middleware for the application.
package middleware

import (
	"net/http"
	"strings"

	services "cry-api/app/services/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware verifies JWT tokens in Authorization header.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Expect format: Bearer <token>
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenStr := tokenParts[1]

		// Parse and verify token
		token, err := jwt.ParseWithClaims(tokenStr, &services.Claims{}, func(_ *jwt.Token) (any, error) {
			return services.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store claims in context for later use and enforce 2FA completion when enabled
		claims, ok := token.Claims.(*services.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if claims.TwoFAEnabled && !claims.TwoFAVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "Two-factor authentication required"})
			c.Abort()
			return
		}

		c.Set("user", claims)

		c.Next()
	}
}

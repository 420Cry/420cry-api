// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"log"
	"net/http"

	JWT "cry-api/app/services/jwt"
	Types "cry-api/app/types/2fa"
	TokenType "cry-api/app/types/token_purpose"

	"github.com/gin-gonic/gin"
)

// AlternativeVerifyOTP validates the alternative OTP (from user's email) and returns a new JWT if successful.
func (h *TwoFactorController) AlternativeVerifyOTP(c *gin.Context) {
	var req Types.ITwoFactorVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.UserUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User UUID is required"})
		return
	}

	if req.OTP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP is required for verification"})
		return
	}

	// Fetch user
	user, err := h.UserService.GetUserByUUID(req.UserUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !user.TwoFAEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has not enabled 2FA"})
		return
	}

	// Verify OTP
	existingToken, err := h.UserTokenService.FindLatestValidToken(user.ID, string(TokenType.TwoFactorAuthAlternativeOTP))
	if err != nil {
		log.Printf("error finding OTP token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	if existingToken == nil || existingToken.Token != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Generate JWT
	newJWT, err := JWT.GenerateJWT(user.UUID, user.Email, user.TwoFAEnabled, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Respond with new JWT
	c.JSON(http.StatusOK, gin.H{
		"jwt": newJWT,
	})
}

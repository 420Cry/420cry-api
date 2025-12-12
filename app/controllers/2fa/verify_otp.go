// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"net/http"

	JWT "cry-api/app/services/jwt"
	TwoFactorTypes "cry-api/app/types/2fa"

	"github.com/gin-gonic/gin"
)

// VerifyOTP validates the OTP and returns a new JWT if successful.
func (h *TwoFactorController) VerifyOTP(c *gin.Context) {
	var req TwoFactorTypes.ITwoFactorVerifyRequest
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

	// Ensure secret exists before verifying
	if user.TwoFASecret == nil || *user.TwoFASecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is not set up for this user"})
		return
	}

	// Verify OTP using provided secret
	isValid, err := h.AuthService.VerifyOTP(*user.TwoFASecret, req.OTP)
	if err != nil || !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	// Enable 2FA flag if not already set
	if !user.TwoFAEnabled {
		user.TwoFAEnabled = true
		if err := h.UserService.UpdateUser(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable 2FA"})
			return
		}
	}

	// Generate JWT with 2FA verified
	jwt, err := JWT.GenerateJWT(user.UUID, user.Email, user.TwoFAEnabled, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Respond with new JWT and user info
	c.JSON(http.StatusOK, gin.H{
		"jwt": jwt,
	})
}

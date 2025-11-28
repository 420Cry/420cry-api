// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"net/http"

	TwoFactorTypes "cry-api/app/types/2fa"

	"github.com/gin-gonic/gin"
)

// Setup generates a 2FA secret and QR code for the authenticated user.
func (h *TwoFactorController) Setup(c *gin.Context) {
	var req TwoFactorTypes.ITwoFactorSetupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.UserUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User UUID is required"})
		return
	}

	user, err := h.UserService.GetUserByUUID(req.UserUUID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.TwoFASecret != nil && *user.TwoFASecret != "" {
		// Use interface method instead of package function
		otpauthURL := h.TwoFactorService.GenerateOtpauthURL(user.Email, *user.TwoFASecret)

		qrCode, err := h.TwoFactorService.GenerateQRCodeBase64(otpauthURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"secret": *user.TwoFASecret,
			"qrCode": qrCode,
		})
		return
	}

	// Generate new secret
	secret, otpauthURL, err := h.TwoFactorService.GenerateTOTP(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate 2FA secret"})
		return
	}

	user.TwoFASecret = &secret
	if err := h.UserService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save 2FA secret"})
		return
	}

	qrCode, err := h.TwoFactorService.GenerateQRCodeBase64(otpauthURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret": secret,
		"qrCode": qrCode,
	})
}

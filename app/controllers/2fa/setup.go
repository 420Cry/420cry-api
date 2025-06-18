// Package controllers handles HTTP requests and responses,
package controllers

import (
	"cry-api/app/factories"
	TwoFactorType "cry-api/app/types/2fa"
	"log"

	"github.com/gin-gonic/gin"
)

// Setup generates a 2FA secret and QR code for the authenticated user.
func (h *TwoFactorController) Setup(c *gin.Context) {
	var req TwoFactorType.ITwoFactorSetupRequest

	// Parse JSON body into req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if req.UserUUID == "" {
		c.JSON(400, gin.H{"error": "User UUID is required"})
		return
	}

	// Get user by UUID
	user, err := h.UserService.GetUserByUUID(req.UserUUID)
	if err != nil || user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Generate TOTP secret and otpauth URL
	secret, otpauthURL, err := factories.GenerateTOTP(user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate 2FA secret"})
		return
	}

	// Generate base64 QR code
	qrCode, err := factories.GenerateQRCodeBase64(otpauthURL)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate QR code"})
		return
	}
	log.Printf("HIHIHI")
	// Return the secret and the QR code image
	c.JSON(200, gin.H{
		"secret": secret,
		"qrCode": qrCode,
	})
}

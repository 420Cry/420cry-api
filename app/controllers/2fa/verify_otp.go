// Package controllers handles HTTP requests and responses,
package controllers

import (
	types "cry-api/app/types/2fa"

	"github.com/gin-gonic/gin"
)

// VerifyOTP Validates OTP and store a 2FA secret
func (h *TwoFactorController) VerifyOTP(c *gin.Context) {
	var req types.ITwoFactorSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Validate UUID presence
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

	// If OTP provided, verify it

	// No OTP provided - maybe respond with QR code & secret or error
	c.JSON(400, gin.H{"error": "OTP is required for verification"})
}

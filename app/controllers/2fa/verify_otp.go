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

	if req.UserUUID == "" {
		c.JSON(400, gin.H{"error": "User UUID is required"})
		return
	}

	if req.OTP == nil || *req.OTP == "" {
		c.JSON(400, gin.H{"error": "OTP is required for verification"})
		return
	}

	valid, err := h.AuthService.VerifyOTP(req.UserUUID, *req.OTP)
	if err != nil || !valid {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "OTP verified successfully"})
}

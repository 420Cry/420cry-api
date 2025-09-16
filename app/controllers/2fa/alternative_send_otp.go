// Package controllers handles HTTP requests and responses.
package controllers

import (
	"github.com/gin-gonic/gin"
)

// AlternativeSendOtp generates a one-time alternative login OTP and sends it via email.
func (h *TwoFactorController) AlternativeSendOtp(c *gin.Context) {
}

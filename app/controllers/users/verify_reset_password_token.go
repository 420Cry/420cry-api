package controllers

import (
	"net/http"
	"time"

	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

// VerifyResetPasswordToken verifies the token with passwords and save it intø the database
func (h *UserController) VerifyResetPasswordToken(c *gin.Context) {
	var req UserTypes.IVerificationResetPasswordForm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON Format"})
		return
	}

	user, err := h.UserService.FindUserByResetPasswordToken(req.ResetPasswordToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot find user"})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is not verified"})
		return
	}

	if user.ResetPasswordTokenCreatedAt == nil || time.Since(*user.ResetPasswordTokenCreatedAt) > time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The token has been expired"})
		return
	}

	// Use the injected PasswordService here:
	hashedPassword, err := h.PasswordService.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot hash password"})
		return
	}

	user.ResetPasswordTokenCreatedAt = nil
	user.ResetPasswordToken = ""
	user.Password = hashedPassword
	err = h.UserService.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

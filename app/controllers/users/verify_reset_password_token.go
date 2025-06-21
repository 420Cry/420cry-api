package controllers

import (
	"net/http"
	"time"

	PasswordService "cry-api/app/services/password"
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

	// Check user by token to have the user holding the token
	user, err := h.UserService.FindUserByResetPasswordToken(req.ResetPasswordToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot find user"})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is not verified"})
	}

	// Save new password and reset token
	if user.ResetPasswordTokenCreatedAt == nil || time.Since(*user.ResetPasswordTokenCreatedAt) > time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The token has been expired"})
		return
	}

	// 	// Create PasswordService instance
	passwordService := PasswordService.NewPasswordService()
	hashedPassword, err := passwordService.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Can not hash password"})
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

package controllers

import (
	UserTypes "cry-api/app/types/users"
	"net/http"
	"time"

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
	foundUser, err := h.UserService.CheckUserByResetPasswordToken(req.ResetPasswordToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot find user"})
		return
	}

	if !foundUser.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is not verified"})
	}

	// Save new password and reset token
	if foundUser.ResetPasswordTokenCreatedAt == nil || time.Since(*foundUser.ResetPasswordTokenCreatedAt) > time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The token has been expired"})
		return
	}

	if err := h.UserService.HandleResetPassword(foundUser, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})

}

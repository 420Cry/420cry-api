package controllers

import (
	UserTypes "cry-api/app/types/users"
	"log"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// VerifyResetPassword verifies the token with passwords and save it intø the database
func (h *UserController) VerifyResetPassword(c *gin.Context) {
	token := c.Param("token")
	log.Printf("Reset password token param: %s", token)

	var req UserTypes.IVerificationResetPasswordForm

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Format"})
		return
	}

	if req.NewPassword == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing fields"})
		return
	}

	if req.NewPassword != req.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Check user by token to have the user holding the token
	foundUser, err := h.UserService.CheckUserByResetPasswordToken(token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot find user"})
		return
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

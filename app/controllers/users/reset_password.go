// Package controllers handles HTTP requests and responses,
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"cry-api/app/config"
	UserModel "cry-api/app/models"
	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

// HandleResetPasswordRequest sends the password request link with reset password token if the user exists
func (h *UserController) HandleResetPasswordRequest(c *gin.Context) {
	var req UserTypes.IVerificationResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Format"})
		return
	}

	foundUser, err := h.UserService.CheckIfUserExists(req.Email)
	if err != nil || foundUser == nil || !foundUser.IsVerified {
		log.Printf("error finding user or user is not verified: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error finding user or user is not verified"})
		return
	}

	var shouldCreateNewToken bool
	if foundUser.ResetPasswordTokenCreatedAt == nil {
		shouldCreateNewToken = true
	} else {
		shouldCreateNewToken = time.Since(*foundUser.ResetPasswordTokenCreatedAt) > time.Hour
	}

	var userToUse *UserModel.User
	if shouldCreateNewToken {
		userToUse, err = h.UserService.CreateResetPasswordToken(foundUser)
		if err != nil {
			log.Printf("error saving reset password token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}
	} else {
		userToUse = foundUser
	}

	go func(user *UserModel.User) {
		cfg := config.Get()

		resetPasswordLink := fmt.Sprintf("%s/auth/reset-password/%s", cfg.CryAppURL, user.ResetPasswordToken)

		err := h.EmailService.SendResetPasswordEmail(
			user.Email,
			cfg.NoReplyEmail,
			user.Username,
			resetPasswordLink,
			cfg.CryAPIURL,
		)
		if err != nil {
			log.Printf("Failed to send reset password email to %s: %v", user.Email, err)
		}
	}(userToUse)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

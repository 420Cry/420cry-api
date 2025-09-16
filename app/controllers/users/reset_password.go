// Package controllers handles HTTP requests and responses
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"cry-api/app/config"
	"cry-api/app/factories"
	UserModel "cry-api/app/models"
	types "cry-api/app/types/token_purpose"
	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

// HandleResetPasswordRequest sends a password reset link with a token if the user exists.
func (h *UserController) HandleResetPasswordRequest(c *gin.Context) {
	var req UserTypes.IVerificationResetPasswordRequest

	// 1️⃣ Parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// 2️⃣ Find user by email
	user, err := h.UserService.FindUserByEmail(req.Email)
	if err != nil {
		log.Printf("internal error finding user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not verified"})
		return
	}

	// 3️⃣ Check for existing valid reset password token
	existingToken, err := h.UserService.FindUserTokenByPurpose(user.ID, string(types.ResetPassword))
	if err != nil {
		log.Printf("error checking existing reset token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	var resetTokenObj *UserModel.UserToken
	if existingToken == nil || existingToken.ExpiresAt.Before(time.Now()) {
		// 4️⃣ Generate new reset password token
		resetTokenObj, err = factories.NewUserToken(user.ID, string(types.ResetPassword), time.Hour, factories.LongLink)
		if err != nil {
			log.Printf("failed to generate reset password token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}

		// 5️⃣ Save new token
		if err := h.UserTokenService.Save(resetTokenObj); err != nil {
			log.Printf("error saving reset password token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save token"})
			return
		}
	} else {
		// Use the existing token if still valid
		resetTokenObj = existingToken
	}

	// 6️⃣ Send reset password email asynchronously
	go func(user *UserModel.User, token *UserModel.UserToken) {
		cfg := config.Get()
		resetPasswordLink := fmt.Sprintf("%s/auth/reset-password/%s", cfg.CryAppURL, token.Token)

		if err := h.EmailService.SendResetPasswordEmail(
			user.Email,
			cfg.NoReplyEmail,
			user.Username,
			resetPasswordLink,
			cfg.CryAPIURL,
		); err != nil {
			log.Printf("Failed to send reset password email to %s: %v", user.Email, err)
		}
	}(user, resetTokenObj)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

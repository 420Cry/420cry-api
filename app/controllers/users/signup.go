// Package controllers handles HTTP requests and responses,
package controllers

import (
	"fmt"
	"log"
	"net/http"

	"cry-api/app/config"
	UserModel "cry-api/app/models"
	EnvTypes "cry-api/app/types/env"
	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

/*
Signup handles user registration requests.
It decodes the incoming JSON request into a UserDomain.User struct,
creates a new user using the userService, and sends a verification email asynchronously.
Responds with a success status if user creation is successful.
*/
func (h *UserController) Signup(c *gin.Context) {
	cfg := config.Get()

	var input UserTypes.IUserSignupRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	createdUser, err := h.UserService.CreateUser(input.Fullname, input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Send email asynchronously
	go func(user *UserModel.User, cfg *EnvTypes.EnvConfig) {
		if user.AccountVerificationToken == nil {
			log.Printf("Cannot send email: account verification token is nil")
			return
		}

		verificationLink := fmt.Sprintf("%s/auth/signup/verify?token=%s", cfg.CryAppURL, *user.AccountVerificationToken)
		err := h.EmailService.SendVerifyAccountEmail(
			user.Email,
			cfg.NoReplyEmail,
			user.Username,
			verificationLink,
			user.VerificationTokens,
		)
		if err != nil {
			log.Printf("Failed to send verification email to %s: %v", user.Email, err)
		}
	}(createdUser, cfg)

	c.JSON(http.StatusCreated, gin.H{"success": true})
}

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

	createdUser, token, err := h.UserService.CreateUser(input.Fullname, input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(mapUserCreationErrorToStatusCode(err.Error()), gin.H{"error": err.Error()})
		return
	}

	go func(user *UserModel.User, token string, cfg *EnvTypes.EnvConfig) {
		verificationLink := fmt.Sprintf("%s/auth/signup/verify?token=%s", cfg.CryAppURL, token)
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
	}(createdUser, token, cfg)

	c.JSON(http.StatusCreated, gin.H{"success": true})
}

// mapUserCreationErrorToStatusCode maps the error message to an HTTP status code
func mapUserCreationErrorToStatusCode(errMessage string) int {
	switch errMessage {
	case "username is already taken":
		return http.StatusConflict
	case "email is already taken":
		return http.StatusConflict
	case "failed to generate signup token":
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	SignUpError "cry-api/app/types/errors"
	types "cry-api/app/types/token_purpose"

	"cry-api/app/config"
	UserModel "cry-api/app/models"
	EnvTypes "cry-api/app/types/env"
	UserTypes "cry-api/app/types/users"

	"cry-api/app/factories"

	"github.com/gin-gonic/gin"
)

/*
Signup handles user registration requests.
It decodes the incoming JSON request, creates a new user,
generates an account verification link token and OTP, saves them,
and sends a verification email containing both asynchronously.
*/
func (h *UserController) Signup(c *gin.Context) {
	cfg := config.Get()

	var input UserTypes.IUserSignupRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// 1️⃣ Create the user
	createdUser, err := h.UserService.CreateUser(input.Fullname, input.Username, input.Email, input.Password)
	if errors.Is(err, SignUpError.ErrUserConflict) {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// 2️⃣ Generate account verification link token (long link)
	linkToken, err := factories.NewUserToken(
		createdUser.ID,
		string(types.AccountVerification),
		24*time.Hour,
		factories.LongLink,
	)
	if err != nil {
		log.Printf("failed to generate account verification link token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate verification token"})
		return
	}

	// 3️⃣ Generate OTP token
	otpToken, err := factories.NewUserToken(
		createdUser.ID,
		string(types.AccountVerificationOTP),
		10*time.Minute,
		factories.OTP,
	)
	if err != nil {
		log.Printf("failed to generate OTP token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate OTP"})
		return
	}

	// 4️⃣ Save both tokens
	if err := h.UserTokenService.Save(linkToken); err != nil {
		log.Printf("failed to save account verification link token for user %d: %v", createdUser.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save verification token"})
		return
	}
	if err := h.UserTokenService.Save(otpToken); err != nil {
		log.Printf("failed to save OTP token for user %d: %v", createdUser.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save OTP"})
		return
	}

	// 5️⃣ Send email asynchronously with both link and OTP
	go func(user *UserModel.User, linkToken, otpToken *UserModel.UserToken, cfg *EnvTypes.EnvConfig) {
		verificationLink := fmt.Sprintf("%s/auth/signup/verify?token=%s", cfg.CryAppURL, linkToken.Token)
		err := h.EmailService.SendVerifyAccountEmail(
			user.Email,
			cfg.NoReplyEmail,
			user.Username,
			verificationLink,
			otpToken.Token,
		)
		if err != nil {
			log.Printf("Failed to send verification email to %s: %v", user.Email, err)
		}
	}(createdUser, linkToken, otpToken, cfg)

	c.JSON(http.StatusCreated, gin.H{"success": true})
}

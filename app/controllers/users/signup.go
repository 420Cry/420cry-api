// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"cry-api/app/config"
	"cry-api/app/logger"
	"cry-api/app/middleware"
	UserModel "cry-api/app/models"
	EnvTypes "cry-api/app/types/env"
	app_errors "cry-api/app/types/errors"
	types "cry-api/app/types/token_purpose"
	"cry-api/app/validators"

	"cry-api/app/factories"

	"github.com/gin-gonic/gin"
)

/*
Signup handles user registration requests.
It validates the incoming request, creates a new user,
generates an account verification link token and OTP, saves them,
and sends a verification email containing both asynchronously.
*/
func (h *UserController) Signup(c *gin.Context) {
	cfg := config.Get()
	logger := logger.GetLogger()

	// Validate request input
	input, err := validators.ValidateUserSignup(c)
	if err != nil {
		logger.WithError(err).Warn("User signup validation failed")
		middleware.AbortWithError(c, err)
		return
	}

	// Create the user
	isVerified := false
	isProfileCompleted := true
	createdUser, err := h.UserService.CreateUser(input.Fullname, input.Username, input.Email, input.Password, isVerified, isProfileCompleted)
	if errors.Is(err, app_errors.ErrUserConflict) {
		logger.WithField("email", input.Email).Warn("User signup failed - user already exists")
		middleware.AbortWithError(c, app_errors.ErrUserConflict)
		return
	}
	if err != nil {
		logger.WithError(err).Error("Failed to create user")
		middleware.AbortWithError(c, app_errors.NewInternalServerError("Could not create user"))
		return
	}

	// Generate account verification link token (long link)
	linkToken, err := factories.NewUserToken(
		createdUser.ID,
		string(types.AccountVerification),
		24*time.Hour,
		factories.LongLink,
	)
	if err != nil {
		logger.WithError(err).Error("Failed to generate account verification link token")
		middleware.AbortWithError(c, app_errors.ErrTokenGeneration)
		return
	}

	// Generate OTP token
	otpToken, err := factories.NewUserToken(
		createdUser.ID,
		string(types.AccountVerificationOTP),
		10*time.Minute,
		factories.OTP,
	)
	if err != nil {
		logger.WithError(err).Error("Failed to generate OTP token")
		middleware.AbortWithError(c, app_errors.ErrTokenGeneration)
		return
	}

	// Save both tokens
	if err := h.UserTokenService.Save(linkToken); err != nil {
		logger.WithError(err).WithField("user_id", createdUser.ID).Error("Failed to save account verification link token")
		middleware.AbortWithError(c, app_errors.ErrDatabaseError)
		return
	}
	if err := h.UserTokenService.Save(otpToken); err != nil {
		logger.WithError(err).WithField("user_id", createdUser.ID).Error("Failed to save OTP token")
		middleware.AbortWithError(c, app_errors.ErrDatabaseError)
		return
	}

	// Send email asynchronously with both link and OTP
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
			logger.WithError(err).WithField("email", user.Email).Error("Failed to send verification email")
		} else {
			logger.WithField("email", user.Email).Info("Verification email sent successfully")
		}
	}(createdUser, linkToken, otpToken, cfg)

	logger.WithField("user_id", createdUser.ID).Info("User signup completed successfully")
	c.JSON(http.StatusCreated, gin.H{"success": true})
}

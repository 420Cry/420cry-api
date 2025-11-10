// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"net/http"

	"cry-api/app/logger"
	"cry-api/app/middleware"
	services "cry-api/app/services/jwt"
	app_errors "cry-api/app/types/errors"
	UserTypes "cry-api/app/types/users"
	"cry-api/app/validators"

	"github.com/gin-gonic/gin"
)

/*
UpdateAccountName handles requests to update a user's username.
It validates the incoming request, extracts the user from JWT context,
and updates the user's fullname in the database.
*/
func (h *UserController) UpdateAccountName(c *gin.Context) {
	logger := logger.GetLogger()

	// Validate request input
	var input UserTypes.IUserUpdateAccountNameRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.WithError(err).Warn("User account name update validation failed")
		middleware.AbortWithError(c, app_errors.ErrInvalidInput)
		return
	}

	// Validate username content (length, allowed characters)
	if err := validators.ValidateUsername(input.AccountName); err != nil {
		logger.WithError(err).WithField("username", input.AccountName).Warn("Username validation failed")
		middleware.AbortWithError(c, err)
		return
	}

	// Get user from JWT context
	userClaims, exists := c.Get("user")
	if !exists {
		logger.Warn("User claims not found in context")
		middleware.AbortWithError(c, app_errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	claims, ok := userClaims.(*services.Claims)
	if !ok {
		logger.Warn("Invalid user claims format")
		middleware.AbortWithError(c, app_errors.NewUnauthorizedError("Invalid user claims"))
		return
	}

	// Get user from database
	user, err := h.UserService.GetUserByUUID(claims.UUID)
	if err != nil {
		logger.WithError(err).WithField("user_uuid", claims.UUID).Error("Failed to find user")
		middleware.AbortWithError(c, app_errors.NewInternalServerError("User not found"))
		return
	}

	// Check if the new username is already in use by another user
	existingUser, err := h.UserService.FindUserByUsername(input.AccountName)
	if err != nil {
		logger.WithError(err).WithField("username", input.AccountName).Error("Failed to check username availability")
		middleware.AbortWithError(c, app_errors.NewInternalServerError("Failed to check username availability"))
		return
	}

	if existingUser != nil && existingUser.UUID != user.UUID {
		logger.WithField("username", input.AccountName).Warn("Username update failed - username already in use")
		middleware.AbortWithError(c, app_errors.NewConflictError("username", "Username is already in use"))
		return
	}

	// Update user's username
	user.Username = input.AccountName
	if err := h.UserService.UpdateUser(user); err != nil {
		logger.WithError(err).WithField("user_uuid", claims.UUID).Error("Failed to update user username")
		middleware.AbortWithError(c, app_errors.NewInternalServerError("Failed to update username"))
		return
	}

	logger.WithField("user_uuid", claims.UUID).Info("User username updated successfully")
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Username updated successfully",
	})
}

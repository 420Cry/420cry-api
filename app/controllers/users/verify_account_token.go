// Package controllers handles HTTP requests and responses,
package controllers

import (
	"net/http"
	"time"

	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

// VerifyAccountToken checks if the provided account verification token is valid and not expired.
// It expects a JSON body with a "token" field, retrieves the user associated with the token,
// and ensures the token matches and was created within the last 24 hours.
func (h *UserController) VerifyAccountToken(c *gin.Context) {
	var req UserTypes.IUserVerifyAccountTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	token := req.Token
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	user, err := h.UserService.CheckAccountVerificationToken(token)
	if err != nil || user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is invalid or expired"})
		return
	}

	timeLimit := time.Now().Add(-24 * time.Hour)
	if user.AccountVerificationToken == nil || *user.AccountVerificationToken != token || user.VerificationTokenCreatedAt.Before(timeLimit) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is invalid or expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

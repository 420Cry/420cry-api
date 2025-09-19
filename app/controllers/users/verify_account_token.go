// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"log"
	"net/http"
	"time"

	types "cry-api/app/types/token_purpose"
	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

// VerifyAccountToken checks if the provided account verification long-link token is valid and not expired.
// It expects a JSON body with a "token" field and verifies it via UserTokenService.
func (h *UserController) VerifyAccountToken(c *gin.Context) {
	var req UserTypes.IUserVerifyAccountTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tokenValue := req.Token
	if tokenValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	// 1️⃣ Find the long-link account_verification token
	userTokenObj, err := h.UserTokenService.FindValidToken(tokenValue, string(types.AccountVerification))
	if err != nil {
		log.Printf("error finding token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if userTokenObj == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is invalid or expired"})
		return
	}

	// 2️⃣ Ensure the token is not older than 24 hours (optional if your token repo handles expires_at)
	if userTokenObj.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token has expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "user_id": userTokenObj.UserID})
}

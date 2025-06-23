// Package controllers handles HTTP requests and responses,
package controllers

import (
	"net/http"

	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

/*
VerifyEmailToken checks the validity of the email verification token. (This function is used to verify the email address of a user during the signup process.)
*/
func (h *UserController) VerifyEmailToken(c *gin.Context) {
	var req UserTypes.IVerificationTokenCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.VerificationService.VerifyUserWithTokens(req.UserToken, req.VerifyToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verified": user.IsVerified})
}

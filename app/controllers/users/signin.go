// Package controllers handles HTTP requests and responses,
package controllers

import (
	"net/http"

	JWT "cry-api/app/services/jwt"
	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

// SignIn method. auth + JWT
func (h *UserController) SignIn(c *gin.Context) {
	var req UserTypes.IUserSigninRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	user, err := h.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	jwt, err := JWT.GenerateJWT(user.UUID, user.Email, user.TwoFAEnabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jwt": jwt,
		"user": gin.H{
			"uuid":         user.UUID,
			"fullname":     user.Fullname,
			"email":        user.Email,
			"username":     user.Username,
			"twoFAEnabled": user.TwoFAEnabled,
		},
	})
}

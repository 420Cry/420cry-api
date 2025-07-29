package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	JWT "cry-api/app/services/jwt"
	UserTypes "cry-api/app/types/users"
)

func (h *UserController) HandleCompleteProfile(c *gin.Context) {
	var req UserTypes.IUserSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	existingUser, err := h.UserService.FindUserByEmail(req.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "This user does not exist"})
		return
	}

	if existingUser.IsProfileCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile already completed"})
		return
	}

	existingUser.Fullname = req.Fullname
	existingUser.Username = req.Username

	hashedPassword, err := h.PasswordService.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot hash password token"})
		return
	}

	existingUser.Password = hashedPassword
	existingUser.IsProfileCompleted = true

	if err := h.UserService.UpdateUser(existingUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot update user's information"})
		return
	}

	jwt, err := JWT.GenerateJWT(existingUser.UUID, existingUser.Email, existingUser.TwoFAEnabled, existingUser.TwoFAEnabled)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot generate token"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"jwt": jwt,
		"user": gin.H{
			"uuid":         existingUser.UUID,
			"fullname":     existingUser.Fullname,
			"email":        existingUser.Email,
			"username":     existingUser.Username,
			"twoFAEnabled": existingUser.TwoFAEnabled,
		},
	})
}

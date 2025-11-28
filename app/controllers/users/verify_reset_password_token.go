package controllers

import (
	"log"
	"net/http"
	"time"

	types "cry-api/app/types/token_purpose"
	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

/*
VerifyResetPasswordToken verifies a reset password token.
It checks that the token is valid (purpose = "reset_password"), not expired, and not consumed,
then allows the user to update their password and marks the token as consumed.
*/
func (h *UserController) VerifyResetPasswordToken(c *gin.Context) {
	var req UserTypes.IVerificationResetPasswordForm
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON format"})
		return
	}

	// 1️⃣ Find token with purpose "reset_password"
	userToken, err := h.UserTokenService.FindValidToken(req.ResetPasswordToken, string(types.ResetPassword))
	if err != nil {
		log.Printf("error finding reset password token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	if userToken == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid or expired reset password token"})
		return
	}

	// 2️⃣ Find user associated with the token
	user, err := h.UserService.FindUserByID(userToken.UserID)
	if err != nil || user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}
	if !user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is not verified"})
		return
	}

	// 3️⃣ Optional: enforce a max age
	if userToken.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Reset password token has expired"})
		return
	}

	// 4️⃣ Hash the new password
	hashedPassword, err := h.PasswordService.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	// 5️⃣ Update user password
	user.Password = hashedPassword
	if err := h.UserService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user password"})
		return
	}

	// 6️⃣ Consume the token so it cannot be reused
	if err := h.UserTokenService.ConsumeToken(user.ID, userToken.Token, string(types.ResetPassword)); err != nil {
		log.Printf("failed to consume reset password token: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

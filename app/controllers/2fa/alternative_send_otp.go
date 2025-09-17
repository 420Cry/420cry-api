// Package controllers handles HTTP requests and responses.
package controllers

import (
	"log"
	"net/http"
	"time"

	"cry-api/app/config"
	"cry-api/app/factories"
	UserModel "cry-api/app/models"
	TwoFactorType "cry-api/app/types/2fa"
	types "cry-api/app/types/token_purpose"

	"github.com/gin-gonic/gin"
)

// AlternativeSendOtp generates a one-time alternative login OTP for 2FA
// and sends it to the user via email. The OTP is stored as a user token with
// purpose TwoFactorAuthAlternativeOTP and expires in 5 minutes.
func (h *TwoFactorController) AlternativeSendOtp(c *gin.Context) {
	var req TwoFactorType.ITwoFactorAlternativeRequest

	// 1️⃣ Parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AlternativeSendOtp] invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	log.Printf("[AlternativeSendOtp] request received for email: %s", req.Email)

	// 2️⃣ Find user
	user, err := h.UserService.FindUserByEmail(req.Email)
	if err != nil || user == nil {
		log.Printf("[AlternativeSendOtp] user not found: %s", req.Email)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if !user.IsVerified {
		log.Printf("[AlternativeSendOtp] user not verified: %s", req.Email)
		c.JSON(http.StatusForbidden, gin.H{"error": "User not verified"})
		return
	}

	// 3️⃣ Check if an unexpired OTP already exists
	existingToken, err := h.UserTokenService.FindLatestValidToken(user.ID, string(types.TwoFactorAuthAlternativeOTP))
	if err != nil {
		log.Printf("[AlternativeSendOtp] error checking existing token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// 4️⃣ Create a new OTP if no valid token exists
	if existingToken == nil {
		otpToken, err := factories.NewUserToken(user.ID, string(types.TwoFactorAuthAlternativeOTP), 5*time.Minute, factories.OTP)
		if err != nil {
			log.Printf("[AlternativeSendOtp] could not generate OTP: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate OTP"})
			return
		}

		if err := h.UserTokenService.Save(otpToken); err != nil {
			log.Printf("[AlternativeSendOtp] could not save OTP: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save OTP"})
			return
		}

		// 5️⃣ Send OTP asynchronously
		go func(u *UserModel.User, token string) {
			cfg := config.Get()
			if err := h.EmailService.SendTwoFactorAlternativeEmail(
				u.Email,
				cfg.NoReplyEmail,
				u.Username,
				token,
				5,
			); err != nil {
				log.Printf("[AlternativeSendOtp] failed to send email to %s: %v", u.Email, err)
			}
		}(user, otpToken.Token)
	} else {
		log.Printf("[AlternativeSendOtp] existing valid OTP found for user: %s", req.Email)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "OTP sent successfully"})
}

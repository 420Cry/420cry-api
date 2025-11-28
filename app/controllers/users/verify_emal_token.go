// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"log"
	"net/http"

	types "cry-api/app/types/token_purpose"
	UserTypes "cry-api/app/types/users"

	"github.com/gin-gonic/gin"
)

/*
VerifyEmailToken checks the validity of the email verification token.
It verifies both the long-link token and the OTP token for account verification.
*/
func (h *UserController) VerifyEmailToken(c *gin.Context) {
	var req UserTypes.IVerificationTokenCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 1️⃣ Find the user token by long-link token
	userTokenObj, err := h.UserService.FindUserTokenByValueAndPurpose(req.UserToken, string(types.AccountVerification))
	if err != nil {
		log.Printf("error finding user token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if userTokenObj == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired account verification link"})
		return
	}

	// 2️⃣ Find the OTP token for the same user
	otpTokenObj, err := h.UserTokenService.FindLatestValidToken(userTokenObj.UserID, string(types.AccountVerificationOTP))
	if err != nil {
		log.Printf("error finding OTP token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if otpTokenObj == nil || otpTokenObj.Token != req.VerifyToken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired verification code"})
		return
	}

	// 3️⃣ Mark both tokens as consumed
	if err := h.UserTokenService.ConsumeToken(userTokenObj.UserID, userTokenObj.Token, string(types.AccountVerification)); err != nil {
		log.Printf("failed to consume long-link token: %v", err)
	}
	if err := h.UserTokenService.ConsumeToken(otpTokenObj.UserID, otpTokenObj.Token, string(types.AccountVerificationOTP)); err != nil {
		log.Printf("failed to consume OTP token: %v", err)
	}

	// 4️⃣ Update user as verified
	user, err := h.UserService.FindUserByID(userTokenObj.UserID)
	if err != nil {
		log.Printf("error fetching user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	user.IsVerified = true
	if err := h.UserService.UpdateUser(user); err != nil {
		log.Printf("error updating user verification status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verified": true})
}

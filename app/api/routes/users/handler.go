package users

import (
	"cry-api/app/config"
	EmailCore "cry-api/app/core/email"
	UserCore "cry-api/app/core/users"
	UserDomain "cry-api/app/domain/users"
	EmailServices "cry-api/app/services/email"
	UserServices "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"
	UserTypes "cry-api/app/types/users"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Handler struct {
	userService  *UserServices.UserService
	emailService *EmailServices.EmailService
}

func NewHandler(db *gorm.DB, cfg *EnvTypes.EnvConfig) *Handler {
	userRepo := UserCore.NewGormUserRepository(db)
	emailSender := EmailCore.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailService := EmailServices.NewEmailService(emailSender)
	userService := UserServices.NewUserService(userRepo, emailService)

	return &Handler{userService: userService, emailService: emailService}
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var user UserDomain.User
	cfg := config.Get()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	createdUser, token, err := h.userService.CreateUser(user.Fullname, user.Username, user.Email, user.Password)
	if err != nil {
		RespondError(w, mapUserCreationErrorToStatusCode(err.Error()), err.Error())
		return
	}

	go h.sendVerificationEmail(createdUser, token, cfg)

	RespondJSON(w, http.StatusCreated, map[string]bool{
		"success": true,
	})
}

func (h *Handler) sendVerificationEmail(user *UserDomain.User, token string, cfg *EnvTypes.EnvConfig) {
	verificationLink := fmt.Sprintf("%s/auth/signup/verify?token=%s", cfg.CryAppURL, token)
	err := h.emailService.SendVerifyAccountEmail(user.Email, cfg.NoReplyEmail, user.Username, verificationLink, user.VerificationTokens)
	if err != nil {
		log.Printf("Failed to send verification email to %s: %v", user.Email, err)
	} else {
		log.Printf("Verification email sent to %s", user.Email)
	}
}

func (h *Handler) VerifyEmailToken(w http.ResponseWriter, r *http.Request) {
	var req UserTypes.VerificationTokenCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserToken == "" {
		RespondError(w, http.StatusBadRequest, "User verification token is required")
		return
	}

	if req.VerifyToken == "" {
		RespondError(w, http.StatusBadRequest, "VerifyToken is required")
		return
	}

	// Now verify both token + OTP
	user, err := h.userService.CheckUserByBothTokens(req.UserToken, req.VerifyToken)
	if err != nil || user == nil {
		RespondError(w, http.StatusBadRequest, "Token or OTP is invalid or expired")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]bool{"verified": user.IsVerified})
}

func (h *Handler) Test(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, map[string]bool{"loggedIn": false})
}

func (h *Handler) VerifyAccountToken(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token := req["token"]

	if token == "" {
		RespondError(w, http.StatusBadRequest, "Token is required")
		return
	}
	// Retrieve the user using the provided token
	user, err := h.userService.CheckAccountVerificationToken(token)
	if err != nil || user == nil {
		RespondError(w, http.StatusBadRequest, "Token is invalid or expired")
		return
	}

	// Check if the token matches and the created date is within the last 24 hours
	timeLimit := time.Now().Add(-24 * time.Hour)
	if user.Token != token || user.VerificationTokenCreatedAt.Before(timeLimit) {
		RespondError(w, http.StatusBadRequest, "Token is invalid or expired")
		return
	}

	// If token is valid and within the 24-hour window
	RespondJSON(w, http.StatusOK, map[string]bool{"valid": true})
}

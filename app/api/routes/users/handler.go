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
	"cry-api/app/utils"
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

	if req.Token == "" {
		RespondError(w, http.StatusBadRequest, "Token is required")
		return
	}

	user, err := h.userService.CheckEmailVerificationToken(req.Token)
	if err != nil || user == nil {
		RespondError(w, http.StatusBadRequest, "Token is invalid or expired")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]bool{"verified": user.IsVerified})
}

func (h *Handler) Test(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, map[string]bool{"loggedIn": false})
}

func (h *Handler) VerificationAccountToken(w http.ResponseWriter, r *http.Request) {
	var req UserTypes.VerificationTokenCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		RespondError(w, http.StatusBadRequest, "Token is required")
		return
	}

	// Retrieve the user using the provided token
	user, err := h.userService.CheckAccountVerificationToken(req.Token)
	if err != nil || user == nil {
		RespondError(w, http.StatusBadRequest, "Token is invalid or expired")
		return
	}

	// Check if the token matches and the created date is within the last 24 hours
	timeLimit := time.Now().Add(-24 * time.Hour)
	if user.Token != req.Token || user.CreatedAt.Before(timeLimit) {
		RespondError(w, http.StatusBadRequest, "Token is invalid or expired")
		return
	}

	// If token is valid and within the 24-hour window
	RespondJSON(w, http.StatusOK, map[string]bool{"valid": true})
}

func (h *Handler) HandlePasswordRequest(w http.ResponseWriter, r *http.Request) {
	// Check if user exists (Thinking of splitting into 2 middlewares)
	var req UserTypes.VerificationResetPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" {
		RespondError(w, http.StatusBadRequest, "Email is required")
		return
	}

	user, err := h.userService.CheckIfUserExists(req.Email)

	if err != nil || user == nil {
		RespondJSON(w, http.StatusOK, map[string]bool{"success": true})
		return
	}

	// Generate random token for email sender
	cfg := config.Get()
	resetPasswordToken, err := utils.GenerateRandomToken()

	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Cannot generate this token")
		return
	}

	// Send the email
	go h.sendVerificationEmail(user, resetPasswordToken, cfg)

	// Response with status
	RespondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *Handler) SendResetPasswordEmail(user *UserDomain.User, token string, cfg *EnvTypes.EnvConfig) error {
	resetPasswordLink := fmt.Sprintf("%s/auth/reset-password/%s", cfg.CryAppURL, token)

	err := h.emailService.SendResetPasswordEmail(user.Email, cfg.NoReplyEmail, user.Username, resetPasswordLink)

	if err != nil {
		log.Printf("Error sending email")
	} else {
		log.Printf("Complete sending email")
	}

	return nil
}

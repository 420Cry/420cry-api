package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"cry-api/app/config"
	EmailCore "cry-api/app/core/email"
	JWT "cry-api/app/core/jwt"
	UserCore "cry-api/app/core/users"
	UserDomain "cry-api/app/domain/users"
	EmailServices "cry-api/app/services/email"
	UserServices "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"
	UserTypes "cry-api/app/types/users"

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
	cfg := config.Get()

	// Read raw body for logging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to read request body")
		return
	}

	// Restore body so it can be decoded
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Decode into input struct
	var input UserTypes.UserSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Call your service with fields (service will call NewUser and hash password)
	createdUser, token, err := h.userService.CreateUser(input.Fullname, input.Username, input.Email, input.Password)
	if err != nil {
		RespondError(w, mapUserCreationErrorToStatusCode(err.Error()), err.Error())
		return
	}

	// Send verification email asynchronously
	go h.SendVerificationEmail(createdUser, token, cfg)

	RespondJSON(w, http.StatusCreated, map[string]bool{
		"success": true,
	})
}

func (h *Handler) SendVerificationEmail(user *UserDomain.User, token string, cfg *EnvTypes.EnvConfig) {
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

	user, err := h.userService.VerifyUserWithTokens(req.UserToken, req.VerifyToken)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]bool{"verified": user.IsVerified})
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	user, err := h.userService.AuthenticateUser(req.UserName, req.Password)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := JWT.GenerateJWT(user.UUID, user.Email)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return full user object with token
	response := map[string]any{
		"user": map[string]interface{}{
			"jwt":      token,
			"uuid":     user.UUID,
			"fullname": user.Fullname,
			"email":    user.Email,
			"username": user.Username,
		},
	}

	RespondJSON(w, http.StatusOK, response)
}

func (h *Handler) VerifyAccountToken(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	token := req["token"]
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

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

// Handler handles HTTP requests related to user operations
type Handler struct {
	UserService  UserServices.UserServiceInterface
	EmailService EmailServices.EmailServiceInterface
}

/*
NewHandler initializes and returns a new Handler instance with its dependencies.
It sets up the user repository, email sender, email service, and user service
using the provided GORM database connection and environment configuration.

Parameters:
  - db:   A pointer to the GORM database instance.
  - cfg:  A pointer to the environment configuration.

Returns:
  - A pointer to the initialized Handler.
*/
func NewHandler(db *gorm.DB, cfg *EnvTypes.EnvConfig) *Handler {
	userRepo := UserCore.NewGormUserRepository(db)
	emailSender := EmailCore.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailService := EmailServices.NewEmailService(emailSender)
	userService := UserServices.NewUserService(userRepo, emailService)

	return &Handler{UserService: userService, EmailService: emailService}
}

/*
Signup handles user registration requests.
It decodes the incoming JSON request into a UserDomain.User struct,
creates a new user using the userService, and sends a verification email asynchronously.
Responds with a success status if user creation is successful.
*/
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

	// Create user and get verification token
	createdUser, token, err := h.UserService.CreateUser(input.Fullname, input.Username, input.Email, input.Password)
	if err != nil {
		RespondError(w, mapUserCreationErrorToStatusCode(err.Error()), err.Error())
		return
	}

	// Send verification email asynchronously (merged logic here)
	go func(user *UserDomain.User, token string, cfg *EnvTypes.EnvConfig) {
		verificationLink := fmt.Sprintf("%s/auth/signup/verify?token=%s", cfg.CryAppURL, token)

		err := h.EmailService.SendVerifyAccountEmail(
			user.Email,
			cfg.NoReplyEmail,
			user.Username,
			verificationLink,
			user.VerificationTokens,
		)
		if err != nil {
			log.Printf("Failed to send verification email to %s: %v", user.Email, err)
		} else {
			log.Printf("Verification email sent to %s", user.Email)
		}
	}(createdUser, token, cfg)

	RespondJSON(w, http.StatusCreated, map[string]bool{"success": true})
}

/*
VerifyEmailToken checks the validity of the email verification token. (This function is used to verify the email address of a user during the signup process.)
*/
func (h *Handler) VerifyEmailToken(w http.ResponseWriter, r *http.Request) {
	var req UserTypes.VerificationTokenCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.UserService.VerifyUserWithTokens(req.UserToken, req.VerifyToken)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]bool{"verified": user.IsVerified})
}

// SignIn method. auth + JWT
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	user, err := h.UserService.AuthenticateUser(req.UserName, req.Password)
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

// VerifyAccountToken checks if the provided account verification token is valid and not expired.
// It expects a JSON body with a "token" field, retrieves the user associated with the token,
// and ensures the token matches and was created within the last 24 hours.
func (h *Handler) VerifyAccountToken(w http.ResponseWriter, r *http.Request) {
	var req UserTypes.UserVerifyAccountTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	fmt.Printf("VerifyAccountToken request: %+v\n", req)

	token := req.Token
	if token == "" {
		RespondError(w, http.StatusBadRequest, "Token is required")
		return
	}

	user, err := h.UserService.CheckAccountVerificationToken(token)
	if err != nil || user == nil {
		RespondError(w, http.StatusBadRequest, "Token is invalid or expired")
		return
	}

	timeLimit := time.Now().Add(-24 * time.Hour)
	if user.Token == nil || *user.Token != token || user.VerificationTokenCreatedAt.Before(timeLimit) {
		RespondError(w, http.StatusBadRequest, "Token is invalid or expired")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]bool{"valid": true})
}

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

// Handler handles HTTP requests related to user operations
type Handler struct {
	userService  *UserServices.UserService
	emailService *EmailServices.EmailService
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

	return &Handler{userService: userService, emailService: emailService}
}

/*
Signup handles user registration requests.
It decodes the incoming JSON request into a UserDomain.User struct,
creates a new user using the userService, and sends a verification email asynchronously.
Responds with a success status if user creation is successful.
*/
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

/*
sendVerificationEmail sends a verification email to the specified user asynchronously.
It constructs a verification link using the application's URL and the provided token,
then uses the emailService to send the verification email. Any errors encountered
during the sending process are logged.

Parameters:
  - user:  Pointer to the UserDomain.User to whom the email will be sent.
  - token: The verification token to be included in the verification link.
  - cfg:   Pointer to the application's environment configuration.
*/

func (h *Handler) sendVerificationEmail(user *UserDomain.User, token string, cfg *EnvTypes.EnvConfig) {
	verificationLink := fmt.Sprintf("%s/auth/signup/verify?token=%s", cfg.CryAppURL, token)
	err := h.emailService.SendVerifyAccountEmail(user.Email, cfg.NoReplyEmail, user.Username, verificationLink, user.VerificationTokens)
	if err != nil {
		log.Printf("Failed to send verification email to %s: %v", user.Email, err)
	} else {
		log.Printf("Verification email sent to %s", user.Email)
	}
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

// Test responds with a JSON object indicating the user is not logged in.
func (h *Handler) Test(w http.ResponseWriter, _ *http.Request) {
	RespondJSON(w, http.StatusOK, map[string]bool{"loggedIn": false})
}

/*
VerificationAccountToken checks the validity of the account verification token. (TODO: This need to be refactored to OPT)
*/
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

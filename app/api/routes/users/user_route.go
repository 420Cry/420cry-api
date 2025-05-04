package routes

import (
	EmailApplication "cry-api/app/application/email"
	UserApplication "cry-api/app/application/users"
	"cry-api/app/config"
	EmailCore "cry-api/app/core/email"
	UserCore "cry-api/app/core/users"
	UserDomain "cry-api/app/domain/users"
	types "cry-api/app/types/env"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService  *UserApplication.UserService
	emailService *EmailApplication.EmailService
}

// NewUserHandler creates a new UserHandler with the required dependencies
func NewUserHandler(db *gorm.DB, cfg *types.EnvConfig) *UserHandler {
	userService, emailService := initializeServices(db, cfg)
	return &UserHandler{
		userService:  userService,
		emailService: emailService,
	}
}

// initializeServices initializes the services needed for the user routes
func initializeServices(db *gorm.DB, cfg *types.EnvConfig) (*UserApplication.UserService, *EmailApplication.EmailService) {
	// Create core layer repository and services
	userRepo := UserCore.NewGormUserRepository(db)
	emailSender := EmailCore.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailService := EmailApplication.NewEmailService(emailSender)
	userService := UserApplication.NewUserService(userRepo, emailService)

	return userService, emailService
}

// handleError is a helper function for sending error responses
func handleError(w http.ResponseWriter, errMessage string, statusCode int) {
	http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errMessage), statusCode)
}

// Signup handles the user signup request
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var user UserDomain.User
	cfg := config.Get()

	// Decode JSON request body into user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handleError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Call the application service to create the user
	createdUser, token, err := h.userService.CreateUser(user.Fullname, user.Username, user.Email, user.Password)
	if err != nil {
		handleError(w, err.Error(), mapUserCreationErrorToStatusCode(err.Error()))
		return
	}

	// Trigger sending the verification email asynchronously
	go func() {
		verificationLink := fmt.Sprintf(cfg.CryAppURL+"/auth/signup/verify?token=%s", token)
		verificationTokens := createdUser.VerificationTokens
		err := h.emailService.SendVerifyAccountEmail(createdUser.Email, cfg.NoReplyEmail, createdUser.Username, verificationLink, verificationTokens)
		if err != nil {
			log.Printf("Failed to send verification email to %s: %v", createdUser.Email, err)
		} else {
			log.Printf("Verification email sent to %s", createdUser.Email)
		}
	}()

	// Respond with success message, excluding the token in the response
	response := map[string]string{"message": "User created successfully", "uuid": createdUser.UUID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handleError(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// Test handles the /users/test route for testing purposes
func (h *UserHandler) Test(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]bool{"loggedIn": false} // false for now
	if err := json.NewEncoder(w).Encode(response); err != nil {
		handleError(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
	}
}

// Users sets up the user-related routes with the /users prefix
func Users(r *mux.Router, db *gorm.DB) {
	cfg := config.Get()
	userHandler := NewUserHandler(db, cfg)
	// Route for creating a new user
	r.HandleFunc("/signup", userHandler.Signup).Methods("POST")
	// Test route for /users/test
	r.HandleFunc("/test", userHandler.Test).Methods("GET")
}

// mapErrorToStatusCode maps the error message to an HTTP status code
func mapUserCreationErrorToStatusCode(errMessage string) int {
	switch errMessage {
	case "username is already taken":
		return http.StatusConflict
	case "email is already taken":
		return http.StatusConflict
	case "failed to generate signup token":
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

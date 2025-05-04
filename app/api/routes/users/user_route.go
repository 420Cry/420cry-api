package routes

import (
	EmailApplication "cry-api/app/application/email"
	UserApplication "cry-api/app/application/users"
	"cry-api/app/config"
	EmailCore "cry-api/app/core/email"
	UserCore "cry-api/app/core/users"
	UserDomain "cry-api/app/domain/users"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Users sets up the user-related routes with the /users prefix
func Users(r *mux.Router, db *gorm.DB) {
	// Create core layer repository
	userRepo := UserCore.NewGormUserRepository(db)
	cfg := config.Get()
	emailSender := EmailCore.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailService := EmailApplication.NewEmailService(emailSender)
	userService := UserApplication.NewUserService(userRepo, emailService)

	// Route for creating a new user
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		var user UserDomain.User

		// Decode JSON request body into user struct
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
			return
		}
		log.Printf("Received user: %+v", user)

		// Call the application service to create the user
		createdUser, token, err := userService.CreateUser(user.Username, user.Email, user.Password)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			errorMessage := err.Error()

			// Handle specific error cases
			switch errorMessage {
			case "username is already taken":
				http.Error(w, `{"error": "Username is already taken"}`, http.StatusConflict)
				return
			case "email is already taken":
				http.Error(w, `{"error": "Email is already taken"}`, http.StatusConflict)
				return
			case "failed to generate signup token":
				http.Error(w, `{"error": "Error generating signup token"}`, http.StatusInternalServerError)
				return
			default:
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
				return
			}
		}

		// Trigger sending the verification email
		go func() {
			verificationLink := fmt.Sprintf("https://example.com/verify?token=%s", token)
			err := emailService.SendVerifyAccountEmail(createdUser.Email, cfg.NoReplyEmail, createdUser.Username, verificationLink)
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
			http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
			return
		}
	}).Methods("POST")

	// Test route for /users/test
	r.HandleFunc("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]bool{"loggedIn": false} // false for now
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")
}

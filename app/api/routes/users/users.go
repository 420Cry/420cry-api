package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	controllers "cry-api/app/api/controllers/users"
	models "cry-api/app/api/models/users"
)

// Users sets up the user-related routes with the /users prefix
func Users(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		// Decode JSON request body into user struct
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
			return
		}
		log.Printf("Received user: %+v", user)

		// Create user in DB
		_, err := controllers.CreateUser(db, user)
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

		// Respond with success message without the token
		response := map[string]string{"message": "User created successfully"}
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
		response := map[string]bool{"loggedIn": false} // Simulate a logged-in user
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")
}

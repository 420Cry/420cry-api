package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	controllers "cry-api/app/api/controllers/users"
	types "cry-api/app/types/users"
)

// Users sets up the user-related routes
func Users(r *mux.Router, db *gorm.DB) {
	// Create User
	r.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
		var user types.User

		// Decode JSON request body into user struct
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		log.Println("Inside create-user endpoint")
		log.Printf("Received user: %+v", user)

		// Create user in DB
		createdUser, err := controllers.CreateUser(db, user)
		if err != nil {
			// Handle specific errors and return more detailed messages
			switch {
			case err.Error() == "duplicate username or UUID":
				http.Error(w, "Username or UUID already exists", http.StatusConflict)
			case err.Error() == "failed to generate signup token":
				http.Error(w, "Error generating signup token", http.StatusInternalServerError)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Return created user as JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(createdUser); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}).Methods("POST")
}

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

// TestRoute sets up a simple test endpoint
func TestRoute(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("Hello, World!")); err != nil {
			http.Error(w, "Unable to write response", http.StatusInternalServerError)
		}
	}).Methods("GET")
}

// Users sets up the user-related routes
func Users(r *mux.Router, db *gorm.DB) {
	TestRoute(r)

	// Create User
	r.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
		var user types.User

		// Decode JSON request body into user struct
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		log.Println("Inside create-user endpoint")
		log.Printf("Received user: %+v", user)

		// Create user in DB
		createdUser, err := controllers.CreateUser(db, user)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
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

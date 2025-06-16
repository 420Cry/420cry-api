package users

import (
	"cry-api/app/config"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterRoutes registers the user-related HTTP routes with the provided router.
func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	cfg := config.Get()
	handler := NewHandler(db, cfg)

	r.HandleFunc("/signup", handler.Signup).Methods("POST")
	r.HandleFunc("/verify-email-token", handler.VerifyEmailToken).Methods("POST")
	r.HandleFunc("/verify-account-token", handler.VerifyAccountToken).Methods("POST")
	r.HandleFunc("/signin", handler.SignIn).Methods("POST")
}

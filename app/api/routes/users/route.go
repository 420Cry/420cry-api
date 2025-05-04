package users

import (
	"cry-api/app/config"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	cfg := config.Get()
	handler := NewHandler(db, cfg)

	r.HandleFunc("/signup", handler.Signup).Methods("POST")
	r.HandleFunc("/verify-email-token", handler.VerificationTokenCheck).Methods("GET")
	r.HandleFunc("/test", handler.Test).Methods("GET")
}

// Package api provides the main entry point for registering all API routes.
package api

import (
	users "cry-api/app/api/routes/users"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterAllRoutes sets up all API routes under /users using the provided router and database.
func RegisterAllRoutes(r *mux.Router, db *gorm.DB) {
	usersRouter := r.PathPrefix("/users").Subrouter()
	users.RegisterRoutes(usersRouter, db)
}

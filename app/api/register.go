package api

import (
	UserRoute "cry-api/app/api/routes/users"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterAllRoutes(r *mux.Router, db *gorm.DB) {
	usersRouter := r.PathPrefix("/users").Subrouter()
	UserRoute.Users(usersRouter, db)
}

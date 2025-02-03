package api

import (
	routes "cry-api/app/api/routes/users"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterAllRoutes(r *mux.Router, db *gorm.DB) {
	routes.Users(r, db)
}

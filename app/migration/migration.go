package main

import (
	models "cry-api/app/api/models/users"
	database "cry-api/app/database"
	"log"
)

func main() {
	dbConn, err := database.GetDBConnection()

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Run AutoMigrate for the User model (create tables if not already there)
	err = dbConn.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("Auto-migration failed: ", err)
	}
}

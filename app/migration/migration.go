package main

import (
	database "cry-api/app/database"
	types "cry-api/app/types/users"
	"log"
)

func main() {
	dbConn, err := database.GetDBConnection()

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Run AutoMigrate for the User model (create tables if not already there)
	err = dbConn.AutoMigrate(&types.User{})

	if err != nil {
		log.Fatal("Auto-migration failed: ", err)
	}
}

// MAIN FOR MIGRATION
package main

import (
	"log"

	database "cry-api/app/database"
	UserModel "cry-api/app/models"
)

func main() {
	dbConn, err := database.GetDBConnection()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Run AutoMigrate for the User model (create tables if not already there)
	err = dbConn.AutoMigrate(&UserModel.User{})
	if err != nil {
		log.Fatal("Auto-migration failed: ", err)
	}
}

// MAIN FOR MIGRATION
package main

import (
	"log"

	database "cry-api/app/database"
	Models "cry-api/app/models"
)

func main() {
	dbConn, err := database.GetDBConnection()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Run AutoMigrate for the User model (create tables if not already there)
	err = dbConn.AutoMigrate(&Models.User{}, &Models.Oauth_Accounts{})
	if err != nil {
		log.Fatal("Auto-migration failed: ", err)
	}
}

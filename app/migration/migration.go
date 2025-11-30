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

	// Run AutoMigrate for the User and UserToken models
	err = dbConn.AutoMigrate(&UserModel.User{}, &UserModel.UserToken{}, &UserModel.Oauth_Accounts{})
	if err != nil {
		log.Fatal("Auto-migration failed: ", err)
	}

	log.Println("Migration completed successfully")
}

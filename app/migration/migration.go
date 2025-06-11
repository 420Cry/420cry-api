// MAIN FOR MIGRATION
package main

import (
	database "cry-api/app/database"
	UserDomain "cry-api/app/domain/users"
	"log"
)

func main() {
	dbConn, err := database.GetDBConnection()

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Run AutoMigrate for the User model (create tables if not already there)
	err = dbConn.AutoMigrate(&UserDomain.User{})

	if err != nil {
		log.Fatal("Auto-migration failed: ", err)
	}
}

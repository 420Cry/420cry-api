package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	routes "api/app/api/routes/users"
	"api/app/config"
	"api/app/database"

	"github.com/gorilla/mux"
)

func main() {
	// Load the configuration settings
	cfg := config.Get()
	dbConn, err := database.GetDBConnection()

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Set up the router
	r := mux.NewRouter()

	// Pass the underlying *gorm.DB to the Users route using GetDB()
	routes.Users(r, dbConn.GetDB())

	// Log when the server is starting
	log.Println("Server started on port " + strconv.Itoa(cfg.APIPort))

	// Define the HTTP server
	server := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.APIPort),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start the server and check for errors
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

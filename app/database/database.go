package database

import (
	"cry-api/app/config"
	types "cry-api/app/types/database"
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDatabase function to initialize a new database connection
func NewDatabase(dsn string) (*types.Database, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Ping the database to check the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB from gorm: %v", err)
	}

	// Attempt to ping the database
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Log success
	log.Println("Successfully connected to the database!")

	return &types.Database{DB: db}, nil
}

// GetDBConnection function to load configuration and return a database connection
func GetDBConnection() (*types.Database, error) {
	cfg := config.Get()

	// Database connection string
	dbAddress := cfg.DB
	dbPort := cfg.DBPort
	dbTable := cfg.DBTable
	dsn := "root:@tcp(" + dbAddress + ":" + strconv.Itoa(dbPort) + ")/" + dbTable + "?charset=utf8&parseTime=True&loc=Local"

	// Get the database connection
	dbConn, err := NewDatabase(dsn)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

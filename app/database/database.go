package database

import (
	"cry-api/app/config"
	types "cry-api/app/types/database"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDatabase initializes a new database connection
func NewDatabase(dsn string) (*types.Database, error) {
	// Log the DSN for debugging purposes
	log.Println("Database DSN:", dsn)

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

// GetDBConnection loads configuration and returns a database connection
func GetDBConnection() (*types.Database, error) {
	cfg := config.Get()

	// Database connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUserName, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBDatabase)

	// Get the database connection
	dbConn, err := NewDatabase(dsn)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

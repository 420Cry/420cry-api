// Package database provides functionality for initializing and managing the application's database connection
// using GORM and MySQL. It includes functions to create a new database connection and retrieve a configured
// database instance based on application settings.
package database

import (
	"fmt"

	Config "cry-api/app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDatabase initializes a new database connection
func NewDatabase(dsn string) (*Database, error) {
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

	return &Database{DB: db}, nil
}

// GetDBConnection loads configuration and returns a database connection
func GetDBConnection() (*Database, error) {
	cfg := Config.Get()

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

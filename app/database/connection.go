package database

import (
	"time"

	"gorm.io/gorm"
)

// Database wraps a gorm.DB instance with enhanced functionality
type Database struct {
	DB *gorm.DB
}

// GetDB returns the raw *gorm.DB instance
func (db *Database) GetDB() *gorm.DB {
	return db.DB
}

// AutoMigrate runs GORM's AutoMigrate on provided models
func (db *Database) AutoMigrate(models ...any) error {
	return db.DB.AutoMigrate(models...)
}

// Transaction executes a function within a database transaction
func (db *Database) Transaction(fn func(*gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

// WithTransaction returns a new Database instance with a transaction
func (db *Database) WithTransaction(tx *gorm.DB) *Database {
	return &Database{DB: tx}
}

// ConfigureConnectionPool configures the database connection pool
func (db *Database) ConfigureConnectionPool(maxOpen, maxIdle int, maxLifetime time.Duration) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	return nil
}

// Ping checks if the database connection is alive
func (db *Database) Ping() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close closes the database connection
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Stats returns database connection statistics
func (db *Database) Stats() (interface{}, error) {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return nil, err
	}
	return sqlDB.Stats(), nil
}

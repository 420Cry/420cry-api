package database

import "gorm.io/gorm"

// Database wraps a gorm.DB instance
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

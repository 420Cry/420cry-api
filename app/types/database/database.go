// Package types provides database types and utilities.
//
// DB TYPE. NO IDEA WHY I HAVE GETDB AND AUTOMICRATE HERE, BUT THIS WILL BE REFACED LATER IN CRY-55
package types

import "gorm.io/gorm"

// Database struct holds the GORM DB connection
type Database struct {
	DB *gorm.DB
}

// GetDB provides access to the database connection
func (db *Database) GetDB() *gorm.DB {
	return db.DB
}

// AutoMigrate will auto-migrate the models
func (db *Database) AutoMigrate(models ...interface{}) error {
	return db.DB.AutoMigrate(models...)
}

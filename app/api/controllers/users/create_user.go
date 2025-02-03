package controllers

import (
	types "cry-api/app/types/users"

	"gorm.io/gorm"
)

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, user types.User) (types.User, error) {
	if err := db.Create(&user).Error; err != nil {
		return types.User{}, err
	}
	return user, nil
}

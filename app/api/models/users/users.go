package models

import (
	types "cry-api/app/types/users"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, user types.User) types.User {
	db.Create(&user)
	return user
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *gorm.DB) []types.User {
	var users []types.User
	db.Find(&users)
	return users
}

// GetUserByID retrieves a user by ID
func GetUserByID(db *gorm.DB, id int) (types.User, error) {
	var user types.User
	result := db.First(&user, id)
	if result.Error != nil {
		return types.User{}, result.Error
	}
	return user, nil
}

// VerifyUser updates the user's verification status
func VerifyUser(db *gorm.DB, uuid string, token string) error {
	var user types.User
	result := db.Where("uuid = ? AND signup_token = ?", uuid, token).First(&user)
	if result.Error != nil {
		return fmt.Errorf("invalid token or user not found")
	}

	user.IsVerified = true
	user.SignupToken = "" // Clear the token after verification
	db.Save(&user)
	return nil
}

// UpdateUser updates an existing user
func UpdateUser(db *gorm.DB, id int, updatedUser types.User) (types.User, error) {
	var user types.User
	result := db.First(&user, id)
	if result.Error != nil {
		return types.User{}, result.Error
	}

	updatedUser.UpdatedAt = time.Now()
	db.Model(&user).Updates(updatedUser)
	return updatedUser, nil
}

// DeleteUser removes a user by ID
func DeleteUser(db *gorm.DB, id int) error {
	result := db.Delete(&types.User{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

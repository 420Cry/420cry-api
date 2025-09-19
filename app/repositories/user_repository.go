// Package repositorie provides methods for interacting with the users.
package repositorie

import (
	"fmt"

	UserModel "cry-api/app/models"

	"gorm.io/gorm"
)

// UserRepository defines the set of methods for interacting with user data storage.
type UserRepository interface {
	// Save persists a user to the database. It creates a new user or updates an existing one.
	Save(user *UserModel.User) error

	// FindByUUID retrieves a user by their UUID.
	FindByUUID(uuid string) (*UserModel.User, error)

	// FindByID retrieves a user by their ID.
	FindByID(id int) (*UserModel.User, error)

	// FindByEmail retrieves a user by their email address.
	FindByEmail(email string) (*UserModel.User, error)

	// FindByUsernameOrEmail retrieves a user by either username or email.
	FindByUsernameOrEmail(username, email string) (*UserModel.User, error)

	// FindByUsername retrieves a user by their username.
	FindByUsername(username string) (*UserModel.User, error)

	// Delete removes a user from the database by their ID.
	Delete(userID int) error
}

// GormUserRepository type
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository returns new GormUserRepository
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Save saves new or updates existing users
func (repo *GormUserRepository) Save(user *UserModel.User) error {
	return repo.db.Save(user).Error
}

// FindByUUID retrieves a user by UUID
func (repo *GormUserRepository) FindByUUID(uuid string) (*UserModel.User, error) {
	var user UserModel.User
	err := repo.db.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID retrieves a user by ID
func (repo *GormUserRepository) FindByID(id int) (*UserModel.User, error) {
	var user UserModel.User
	err := repo.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (repo *GormUserRepository) FindByEmail(email string) (*UserModel.User, error) {
	var user UserModel.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsernameOrEmail retrieves a user by username or email
func (repo *GormUserRepository) FindByUsernameOrEmail(username string, email string) (*UserModel.User, error) {
	var user UserModel.User
	err := repo.db.Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername retrieves a user by username
func (repo *GormUserRepository) FindByUsername(username string) (*UserModel.User, error) {
	var user UserModel.User
	err := repo.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Delete removes a user by ID
func (repo *GormUserRepository) Delete(userID int) error {
	if err := repo.db.Delete(&UserModel.User{}, userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

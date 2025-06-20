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

	// FindByEmail retrieves a user by their email address.
	FindByEmail(email string) (*UserModel.User, error)

	// FindByUsernameOrEmail retrieves a user by either username or email.
	FindByUsernameOrEmail(username, email string) (*UserModel.User, error)

	// FindByUsername retrieves a user by their username.
	FindByUsername(username string) (*UserModel.User, error)

	// Delete removes a user from the database by their ID.
	Delete(userID int) error

	// FindByResetPasswordToken retrieves a user by their reset password token
	FindByResetPasswordToken(token string) (*UserModel.User, error)

	// FindByAccountVerificationToken retrieves a user by their account verification token.
	FindByAccountVerificationToken(token string) (*UserModel.User, error)
}

// GormUserRepository type
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository returns new GormUserRepository
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Save saves new users
func (repo *GormUserRepository) Save(user *UserModel.User) error {
	return repo.db.Save(user).Error
}

// FindByUUID retrieves a user from the database from UUID
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

// FindByEmail retrieves a user from the database by email
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

// FindByUsernameOrEmail retrieves a user from the database whose username or email matches the provided values.
// It returns a pointer to the UserDomain.User if found, or nil if no matching user exists.
// If an error occurs during the query (other than record not found), it returns the error.
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

// FindByUsername retrieves a user from the database by their username.
// Returns a pointer to the UserDomain.User if found, or nil if no user exists with the given username.
// If an error occurs during the query (other than record not found), it returns the error.
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

// FindByAccountVerificationToken retrieves a user by their account verification token.
// Returns (nil, nil) if not found, or an error if a DB error occurs.
func (repo *GormUserRepository) FindByAccountVerificationToken(token string) (*UserModel.User, error) {
	var user UserModel.User
	err := repo.db.Where("account_verification_token = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByResetPasswordToken retrieves a user by their reset password token
// Returns (nil, nil) if not found, or an error if a DB error occurs.
func (repo *GormUserRepository) FindByResetPasswordToken(token string) (*UserModel.User, error) {
	var user UserModel.User
	err := repo.db.Where("reset_password_token = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Delete from DATABASE. USE WITH CAUTION PLEASE
func (repo *GormUserRepository) Delete(userID int) error {
	if err := repo.db.Delete(&UserModel.User{}, userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

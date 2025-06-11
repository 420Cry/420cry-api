// Package core provides core functionalities for users.
package core

import (
	"fmt"
	"time"

	UserDomain "cry-api/app/domain/users"

	"gorm.io/gorm"
)

// UserRepository defines the methods needed for user persistence
type UserRepository interface {
	Save(user *UserDomain.User) error
	FindByUsernameOrEmail(username, email string) (*UserDomain.User, error)
	FindByVerificationToken(token string) (*UserDomain.User, error)
}

// GormUserRepository implements the UserRepository interface for GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository instance
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Save persists the user in the database
func (repo *GormUserRepository) Save(user *UserDomain.User) error {
	return repo.db.Save(user).Error
}

// FindByUsernameOrEmail retrieves a user by their username or email
func (repo *GormUserRepository) FindByUsernameOrEmail(username, email string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("username = ?", username).Or("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByVerificationToken retrieves a user by their verification token
func (repo *GormUserRepository) FindByVerificationToken(token string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("verification_tokens = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	// Check if the created_at is older than 24 hours
	if time.Since(user.CreatedAt) > 24*time.Hour {
		return nil, fmt.Errorf("account verification token is invalid or expired")
	}
	return &user, nil
}

// FindByAccountVerificationToken retrieves a user by their account verification token (THIS IS OPT AND LOGIC IS CURRENTLY INCORRECT)
// THIS WILL BE REFACTORED IN CRY-55.
func (repo *GormUserRepository) FindByAccountVerificationToken(token string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("token = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// Check if the created_at is older than 24 hours
	if time.Since(user.CreatedAt) > 24*time.Hour {
		return nil, fmt.Errorf("account verification token is invalid or expired")
	}

	// Check if user is already verified
	if user.IsVerified {
		return nil, fmt.Errorf("account already verified")
	}

	return &user, nil
}

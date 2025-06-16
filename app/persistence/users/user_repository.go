// Package persistence provides methods for interacting with the users.
package persistence

import (
	"fmt"
	"time"

	UserDomain "cry-api/app/domain/users"

	"gorm.io/gorm"
)

// GormUserRepository type
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository returns new GormUserRepository
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Save saves new users
func (repo *GormUserRepository) Save(user *UserDomain.User) error {
	return repo.db.Save(user).Error
}

// FindByUUID retrieves a user from the database from UUID
func (repo *GormUserRepository) FindByUUID(uuid string) (*UserDomain.User, error) {
	var user UserDomain.User
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
func (repo *GormUserRepository) FindByEmail(email string) (*UserDomain.User, error) {
	var user UserDomain.User
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
func (repo *GormUserRepository) FindByUsernameOrEmail(username string, email string) (*UserDomain.User, error) {
	var user UserDomain.User
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
func (repo *GormUserRepository) FindByUsername(username string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByVerificationToken retrieves a user from the database by their verification token.
// It returns the user if the token exists and has not expired (within 24 hours).
// If the token is not found, it returns (nil, nil).
// If the token is found but expired, it returns an error indicating expiration.
// Any other database errors are returned as errors.
func (repo *GormUserRepository) FindByVerificationToken(token string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("verification_tokens = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if time.Since(user.VerificationTokenCreatedAt) > 24*time.Hour {
		return nil, fmt.Errorf("verification token expired")
	}
	return &user, nil
}

// FindByUserToken retrieves a user from the database by their token.
func (repo *GormUserRepository) FindByUserToken(token string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("verification_tokens = ?", token).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if time.Since(user.VerificationTokenCreatedAt) > 24*time.Hour {
		return nil, fmt.Errorf("user token expired")
	}
	return &user, nil
}

// FindByAccountVerificationToken retrieves a user by their account verification token.
// Returns (nil, nil) if not found, or an error if a DB error occurs.
func (repo *GormUserRepository) FindByAccountVerificationToken(token string) (*UserDomain.User, error) {
	var user UserDomain.User
	err := repo.db.Where("token = ?", token).First(&user).Error
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
	if err := repo.db.Delete(&UserDomain.User{}, userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

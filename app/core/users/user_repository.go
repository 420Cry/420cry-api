package core

import (
	UserDomain "cry-api/app/domain/users"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UserRepository defines the methods needed for user persistence
type UserRepository interface {
	Save(user *UserDomain.User) error
	FindByUsernameOrEmail(username, email string) (*UserDomain.User, error)
	FindByVerificationToken(token string) (*UserDomain.User, error)
	FindByAccountVerificationToken(token string) (*UserDomain.User, error)
	Delete(userID int) error
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
	if time.Since(user.VerificationTokenCreatedAt) > 24*time.Hour {
		return nil, fmt.Errorf("account verification token is invalid or expired")
	}
	return &user, nil
}

// FindByUserToken retrieves a user by their verification token
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
		return nil, fmt.Errorf("account verification token is invalid or expired")
	}
	return &user, nil
}

// Delete removes the user from the database (hard delete)
func (repo *GormUserRepository) Delete(userID int) error {
	if err := repo.db.Delete(&UserDomain.User{}, userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

// FindByAccountVerificationToken retrieves a user by their account verification token
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

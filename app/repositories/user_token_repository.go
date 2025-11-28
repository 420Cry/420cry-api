// Package repositorie provides methods for interacting with the user tokens.
package repositorie

import (
	"fmt"
	"time"

	UserModel "cry-api/app/models"

	"gorm.io/gorm"
)

// UserTokenRepository defines methods for interacting with user tokens.
type UserTokenRepository interface {
	// Save persists a token to the database
	Save(token *UserModel.UserToken) error

	// FindValidToken retrieves a valid (non-expired, non-consumed) token by value and purpose
	FindValidToken(token, purpose string) (*UserModel.UserToken, error)

	// ConsumeToken marks a token as consumed
	ConsumeToken(userID int, token, purpose string) error

	// DeleteExpired removes expired tokens
	DeleteExpired() error

	// FindLatestValidToken retrieves all tokens for a user (optionally by purpose)
	FindLatestValidToken(userID int, purpose string) (*UserModel.UserToken, error)
}

// GormUserTokenRepository implements UserTokenRepository using GORM
type GormUserTokenRepository struct {
	db *gorm.DB
}

// NewGormUserTokenRepository returns a new GormUserTokenRepository
func NewGormUserTokenRepository(db *gorm.DB) *GormUserTokenRepository {
	return &GormUserTokenRepository{db: db}
}

// Save inserts or updates a token
func (repo *GormUserTokenRepository) Save(token *UserModel.UserToken) error {
	return repo.db.Save(token).Error
}

// FindValidToken retrieves a valid (non-expired, non-consumed) token that hasn't been used
func (repo *GormUserTokenRepository) FindValidToken(token, purpose string) (*UserModel.UserToken, error) {
	var userToken UserModel.UserToken
	err := repo.db.
		Where("token = ? AND purpose = ? AND consumed = ? AND used_at IS NULL AND expires_at > ?", token, purpose, false, time.Now()).
		First(&userToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &userToken, nil
}

// ConsumeToken marks a token as consumed
func (repo *GormUserTokenRepository) ConsumeToken(userID int, token, purpose string) error {
	result := repo.db.Model(&UserModel.UserToken{}).
		Where("user_id = ? AND token = ? AND purpose = ? AND consumed = ?", userID, token, purpose, false).
		Updates(map[string]interface{}{
			"consumed": true,
			"used_at":  time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("token not found or already consumed")
	}
	return nil
}

// DeleteExpired removes all expired tokens
func (repo *GormUserTokenRepository) DeleteExpired() error {
	return repo.db.Where("expires_at <= ?", time.Now()).Delete(&UserModel.UserToken{}).Error
}

// FindLatestValidToken retrieves the latest valid (non-expired, non-consumed) token for a user and purpose
func (repo *GormUserTokenRepository) FindLatestValidToken(userID int, purpose string) (*UserModel.UserToken, error) {
	var token UserModel.UserToken
	err := repo.db.
		Where("user_id = ? AND purpose = ? AND consumed = ? AND expires_at > ?", userID, purpose, false, time.Now()).
		Order("created_at DESC"). // pick the latest token
		First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// Package services provides business logic for user tokens
package services

import (
	"time"

	UserModel "cry-api/app/models"
	UserRepository "cry-api/app/repositories"
)

// UserTokenServiceInterface defines the contract for token operations
type UserTokenServiceInterface interface {
	Save(token *UserModel.UserToken) error
	FindValidToken(token, purpose string) (*UserModel.UserToken, error)
	FindLatestValidToken(userID int, purpose string) (*UserModel.UserToken, error)
	ConsumeToken(userID int, token, purpose string) error
	DeleteExpired() error
}

// UserTokenService handles token-related operations
type UserTokenService struct {
	tokenRepo UserRepository.UserTokenRepository
}

// NewUserTokenService creates a new instance of UserTokenService
func NewUserTokenService(tokenRepo UserRepository.UserTokenRepository) *UserTokenService {
	return &UserTokenService{tokenRepo: tokenRepo}
}

// Save persists a token
func (s *UserTokenService) Save(token *UserModel.UserToken) error {
	return s.tokenRepo.Save(token)
}

// FindValidToken retrieves a valid token by value and purpose
func (s *UserTokenService) FindValidToken(token, purpose string) (*UserModel.UserToken, error) {
	t, err := s.tokenRepo.FindValidToken(token, purpose)
	if err != nil {
		return nil, err
	}
	if t == nil || t.ExpiresAt.Before(time.Now()) {
		return nil, nil
	}
	return t, nil
}

// FindLatestValidToken retrieves the latest valid (non-expired, non-consumed) token for a user and purpose
func (s *UserTokenService) FindLatestValidToken(userID int, purpose string) (*UserModel.UserToken, error) {
	return s.tokenRepo.FindLatestValidToken(userID, purpose)
}

// ConsumeToken marks a token as used
func (s *UserTokenService) ConsumeToken(userID int, token, purpose string) error {
	return s.tokenRepo.ConsumeToken(userID, token, purpose)
}

// DeleteExpired removes all expired tokens
func (s *UserTokenService) DeleteExpired() error {
	return s.tokenRepo.DeleteExpired()
}

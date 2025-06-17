package services

import (
	"errors"

	UserModel "cry-api/app/models"
	UserRepository "cry-api/app/repositories"
	PasswordService "cry-api/app/services/password"
)

// AuthService handles user authentication.
type AuthService struct {
	userRepo UserRepository.UserRepository
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(userRepo UserRepository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// AuthenticateUser verifies username and password, and checks if user is verified.
func (s *AuthService) AuthenticateUser(username, password string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if err := PasswordService.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid password")
	}
	if !user.IsVerified {
		return nil, errors.New("user not verified")
	}
	return user, nil
}

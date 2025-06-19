package services

import (
	"errors"

	UserModel "cry-api/app/models"
	UserRepository "cry-api/app/repositories"
	PasswordService "cry-api/app/services/password"
)

// AuthService handles user authentication.
type AuthService struct {
	userRepo        UserRepository.UserRepository
	passwordService PasswordService.PasswordServiceInterface
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(userRepo UserRepository.UserRepository, passwordService PasswordService.PasswordServiceInterface) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		passwordService: passwordService,
	}
}

// AuthServiceInterface defines the contract for user service methods.
type AuthServiceInterface interface {
	AuthenticateUser(username, password string) (*UserModel.User, error)
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
	if err := s.passwordService.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid password")
	}
	if !user.IsVerified {
		return nil, errors.New("user not verified")
	}
	return user, nil
}

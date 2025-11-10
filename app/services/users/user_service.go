// Package services provides business logic for user-related operations.
package services

import (
	"errors"
	"fmt"

	"cry-api/app/factories"
	UserModel "cry-api/app/models"
	UserRepository "cry-api/app/repositories"
	AuthService "cry-api/app/services/auth"
	EmailService "cry-api/app/services/email"
	SignUpError "cry-api/app/types/errors"

	"gorm.io/gorm"
)

// UserService handles user-related business logic such as
// creating users, authenticating, and verifying accounts.
type UserService struct {
	userRepo      UserRepository.UserRepository // User data repository interface
	userTokenRepo UserRepository.UserTokenRepository
	emailService  EmailService.EmailServiceInterface // Email service interface for sending emails
	authService   AuthService.AuthServiceInterface
}

// UserServiceInterface defines the contract for user service methods.
type UserServiceInterface interface {
	CreateUser(fullname, username, email, password string) (*UserModel.User, error)
	GetUserByUUID(uuid string) (*UserModel.User, error)
	UpdateUser(user *UserModel.User) error
	FindUserByEmail(email string) (*UserModel.User, error)
	FindUserByUsername(username string) (*UserModel.User, error)
	FindUserByID(id int) (*UserModel.User, error)
	FindUserTokenByPurpose(userID int, purpose string) (*UserModel.UserToken, error)
	FindUserTokenByValueAndPurpose(tokenValue, purpose string) (*UserModel.UserToken, error)
}

// NewUserService creates a new instance of UserService with provided user repository and email service.
func NewUserService(
	userRepo UserRepository.UserRepository,
	userTokenRepo UserRepository.UserTokenRepository,
	emailService EmailService.EmailServiceInterface,
	authService AuthService.AuthServiceInterface,
) *UserService {
	return &UserService{
		userRepo:      userRepo,
		userTokenRepo: userTokenRepo,
		emailService:  emailService,
		authService:   authService,
	}
}

// GetUserByUUID returns user or nil
func (s *UserService) GetUserByUUID(uuid string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user or refreshes the verification token if user exists but is unverified.
func (s *UserService) CreateUser(fullname, username, email, password string) (*UserModel.User, error) {
	// Check if a user with the same username or email already exists
	existingUser, err := s.userRepo.FindByUsernameOrEmail(username, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		// Return 409 Conflict error if user exists
		return nil, SignUpError.ErrUserConflict
	}

	// Create a new user instance using factory
	newUser, err := factories.NewUser(fullname, username, email, password)
	if err != nil {
		return nil, err
	}

	// Persist the new user to the repository
	if err := s.userRepo.Save(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// UpdateUser updates the user in the repository.
func (s *UserService) UpdateUser(user *UserModel.User) error {
	return s.userRepo.Save(user)
}

/* FindUserByEmail checks the user information by email address and return accordinglyy*/
func (s *UserService) FindUserByEmail(email string) (*UserModel.User, error) {
	foundUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding the user for this email: %w", err)
	}

	return foundUser, nil
}

/* FindUserByUsername returns user by username */
func (s *UserService) FindUserByUsername(username string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding the user for this username: %w", err)
	}
	return user, nil
}

/* FindUserByID returns user by id */
func (s *UserService) FindUserByID(id int) (*UserModel.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// FindUserTokenByPurpose finds a valid token for a given user ID and purpose
func (s *UserService) FindUserTokenByPurpose(userID int, purpose string) (*UserModel.UserToken, error) {
	return s.userTokenRepo.FindLatestValidToken(userID, purpose)
}

// FindUserTokenByValueAndPurpose returns user if tokenValue and purpose are matched
func (s *UserService) FindUserTokenByValueAndPurpose(tokenValue, purpose string) (*UserModel.UserToken, error) {
	return s.userTokenRepo.FindValidToken(tokenValue, purpose)
}

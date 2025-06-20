// Package services provides business logic for user-related operations.
package services

import (
	"fmt"

	"cry-api/app/factories"
	UserModel "cry-api/app/models"
	UserRepository "cry-api/app/repositories"
	EmailService "cry-api/app/services/email"
	SignUpError "cry-api/app/types/errors"
)

// UserService handles user-related business logic such as
// creating users, authenticating, and verifying accounts.
type UserService struct {
	userRepo            UserRepository.UserRepository      // User data repository interface
	emailService        EmailService.EmailServiceInterface // Email service interface for sending emails
	VerificationService VerificationServiceInterface
	AuthService         AuthServiceInterface
}

// UserServiceInterface defines the contract for user service methods.
type UserServiceInterface interface {
	CreateUser(fullname, username, email, password string) (*UserModel.User, error)
	GetUserByUUID(uuid string) (*UserModel.User, error)
	UpdateUser(user *UserModel.User) error
	FindUserByEmail(email string) (*UserModel.User, error)
	FindUserByResetPasswordToken(token string) (*UserModel.User, error)
}

// NewUserService creates a new instance of UserService with provided user repository and email service.
func NewUserService(
	userRepo UserRepository.UserRepository,
	emailService EmailService.EmailServiceInterface,
	verificationService VerificationServiceInterface,
	authService AuthServiceInterface,
) *UserService {
	return &UserService{
		userRepo:            userRepo,
		emailService:        emailService,
		VerificationService: verificationService,
		AuthService:         authService,
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
		return nil, fmt.Errorf("error finding the user for this email")
	}

	if foundUser == nil {
		return nil, fmt.Errorf("no user found using this email: %s", email)
	}

	return foundUser, nil
}

// CheckUserByResetPasswordToken finds users based on reset password token
func (s *UserService) FindUserByResetPasswordToken(token string) (*UserModel.User, error) {
	foundUser, err := s.userRepo.FindByResetPasswordToken(token)
	if err != nil {
		return nil, fmt.Errorf("error finding the user for this token")
	}

	if foundUser == nil {
		return nil, fmt.Errorf("no user found using this email")
	}

	return foundUser, nil
}

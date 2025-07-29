// Package services provides business logic for user-related operations.
package services

import (
	"errors"
	"fmt"

	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"cry-api/app/factories"
	Models "cry-api/app/models"
	Repository "cry-api/app/repositories"
	EmailService "cry-api/app/services/email"
	SignUpError "cry-api/app/types/errors"
	OAuthType "cry-api/app/types/oauth"
)

// UserService handles user-related business logic such as
// creating users, authenticating, and verifying accounts.
type UserService struct {
	userRepo            Repository.UserRepository // User data repository interface
	transactionRepo     Repository.TransactionRepository
	emailService        EmailService.EmailServiceInterface // Email service interface for sending emails
	VerificationService VerificationServiceInterface
	AuthService         AuthServiceInterface
}

// UserServiceInterface defines the contract for user service methods.
type UserServiceInterface interface {
	CreateUser(fullname, username, email, password string, isVerified bool, isProfileCompleted bool) (*Models.User, error)
	CreateUserByGoogleAuth(googleUserInfo *OAuthType.IGoogleUserResponse, token *oauth2.Token) (*Models.User, error)
	GetUserByUUID(uuid string) (*Models.User, error)
	UpdateUser(user *Models.User) error
	FindUserByEmail(email string) (*Models.User, error)
	FindUserByResetPasswordToken(token string) (*Models.User, error)
}

// NewUserService creates a new instance of UserService with provided user repository and email service.
func NewUserService(
	userRepo Repository.UserRepository,
	transactionRepo Repository.TransactionRepository,
	emailService EmailService.EmailServiceInterface,
	verificationService VerificationServiceInterface,
	authService AuthServiceInterface,
) *UserService {
	return &UserService{
		userRepo:            userRepo,
		transactionRepo:     transactionRepo,
		emailService:        emailService,
		VerificationService: verificationService,
		AuthService:         authService,
	}
}

// GetUserByUUID returns user or nil
func (s *UserService) GetUserByUUID(uuid string) (*Models.User, error) {
	user, err := s.userRepo.FindByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user or refreshes the verification token if user exists but is unverified.
func (s *UserService) CreateUser(fullname, username, email, password string, isVerified bool, isProfileCompleted bool) (*Models.User, error) {
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
	newUser, err := factories.NewUser(fullname, username, email, password, isVerified, isProfileCompleted)
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
func (s *UserService) UpdateUser(user *Models.User) error {
	return s.userRepo.Save(user)
}

/* FindUserByEmail checks the user information by email address and return accordinglyy*/
func (s *UserService) FindUserByEmail(email string) (*Models.User, error) {
	foundUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding the user for this email: %w", err)
	}

	return foundUser, nil
}

// FindUserByResetPasswordToken finds users based on reset password token
func (s *UserService) FindUserByResetPasswordToken(token string) (*Models.User, error) {
	foundUser, err := s.userRepo.FindByResetPasswordToken(token)
	if err != nil {
		return nil, fmt.Errorf("error finding the user for this token")
	}

	if foundUser == nil {
		return nil, fmt.Errorf("no user found using this email")
	}

	return foundUser, nil
}

func (s *UserService) CreateUserByGoogleAuth(googleUserInfo *OAuthType.IGoogleUserResponse, token *oauth2.Token) (*Models.User, error) {
	randomPassword, err := factories.Generate32ByteToken()

	if err != nil {
		return nil, err
	}

	isVerified := true
	isProfileCompleted := false

	createdUser, err := factories.NewUser(googleUserInfo.GivenName, googleUserInfo.Email, googleUserInfo.Email, randomPassword, isVerified, isProfileCompleted)

	if err != nil {
		return nil, err
	}

	provider := "Google"
	providerId := googleUserInfo.Sub
	oauthAccount, err := factories.NewOAuthAccount(createdUser, provider, providerId, googleUserInfo.Email, token)

	if err != nil {
		return nil, err
	}

	if err := s.transactionRepo.CreateUserByOAuth(createdUser, oauthAccount); err != nil {
		return nil, err
	}

	return createdUser, nil
}

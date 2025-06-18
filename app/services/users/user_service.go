// Package services provides business logic for user-related operations.
package services

import (
	"errors"
	"fmt"
	"time"

	"cry-api/app/factories"
	UserModel "cry-api/app/models"
	UserRepository "cry-api/app/repositories"
	EmailService "cry-api/app/services/email"
	PasswordService "cry-api/app/services/password"
)

// UserService handles user-related business logic such as
// creating users, authenticating, and verifying accounts.
type UserService struct {
	userRepo     UserRepository.UserRepository      // User data repository interface
	emailService EmailService.EmailServiceInterface // Email service interface for sending emails
}

// UserServiceInterface defines the contract for user service methods.
type UserServiceInterface interface {
	CreateUser(fullname, username, email, password string) (*UserModel.User, error)
	AuthenticateUser(username, password string) (*UserModel.User, error)
	VerifyUserWithTokens(userToken, verifyToken string) (*UserModel.User, error)
	CheckAccountVerificationToken(token string) (*UserModel.User, error)
	CheckEmailVerificationToken(token string) (*UserModel.User, error)
	CheckUserByBothTokens(token string, verificationToken string) (*UserModel.User, error)
}

// NewUserService creates a new instance of UserService with provided user repository and email service.
func NewUserService(userRepo UserRepository.UserRepository, emailService EmailService.EmailServiceInterface) *UserService {
	return &UserService{userRepo: userRepo, emailService: emailService}
}

// CreateUser creates a new user or refreshes the verification token if user exists but is unverified.
// Returns the created user or an error.
func (s *UserService) CreateUser(fullname, username, email, password string) (*UserModel.User, error) {
	// Check if a user with the same username or email already exists
	existingUser, err := s.userRepo.FindByUsernameOrEmail(username, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		// Handle existing user case: refresh token if unverified and expired
		refreshed, err := s.handleExistingUser(existingUser, username, email)
		if err != nil {
			return nil, err
		}
		if refreshed != nil {
			return refreshed, nil
		}
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

// AuthenticateUser verifies the username and password and returns the user if valid and verified.
func (s *UserService) AuthenticateUser(username, password string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check password validity using password service
	if err := PasswordService.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid password")
	}

	// Ensure user is verified before allowing authentication
	if !user.IsVerified {
		return nil, errors.New("user not verified")
	}

	return user, nil
}

// VerifyUserWithTokens verifies user account with matching user token and verification token.
// If successful, marks the user as verified and clears tokens.
func (s *UserService) VerifyUserWithTokens(token, verificationToken string) (*UserModel.User, error) {
	user, err := s.CheckUserByBothTokens(token, verificationToken)
	if err != nil {
		return nil, err
	}

	user.IsVerified = true
	user.VerificationTokens = ""
	user.AccountVerificationToken = nil

	// Save updated user status
	if err := s.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to update user verification status: %w", err)
	}

	return user, nil
}

// CheckAccountVerificationToken returns the user associated with the given account token.
// Returns error if token is invalid or user not found.
func (s *UserService) CheckAccountVerificationToken(token string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByAccountVerificationToken(token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid account token")
	}
	return user, nil
}

// CheckEmailVerificationToken verifies the user's email by token,
// marks user as verified, clears verification tokens, and saves user.
func (s *UserService) CheckEmailVerificationToken(token string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("invalid verification token")
	}

	user.IsVerified = true
	user.VerificationTokens = ""

	err = s.userRepo.Save(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CheckUserByBothTokens checks that the user matches both the user token and verification token.
// Returns user if tokens match; error otherwise.
func (s *UserService) CheckUserByBothTokens(token, verificationToken string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByVerificationToken(verificationToken)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("invalid verification token")
	}

	if user.AccountVerificationToken == nil || *user.AccountVerificationToken != token {
		return nil, fmt.Errorf("token does not match")
	}

	return user, nil
}

// handleExistingUser handles the case when a user with the same username or email already exists.
// If user is unverified and verification token expired, refreshes the token.
// Returns the updated user or nil.
func (s *UserService) handleExistingUser(existingUser *UserModel.User, username, email string) (*UserModel.User, error) {
	if existingUser.Username == username || existingUser.Email == email {
		if !existingUser.IsVerified {
			// If verification token expired (older than 24h), generate a new one
			if time.Since(existingUser.VerificationTokenCreatedAt) > 24*time.Hour {
				newVerificationToken, err := factories.GenerateVerificationToken()
				if err != nil {
					return nil, err
				}

				existingUser.VerificationTokens = newVerificationToken
				existingUser.VerificationTokenCreatedAt = time.Now()

				if err := s.userRepo.Save(existingUser); err != nil {
					return nil, err
				}
				return existingUser, nil
			}
			// Token still valid, return user as-is
			return existingUser, nil
		}
		return nil, fmt.Errorf("user with %s is already verified", username)
	}
	return nil, nil
}

/* CheckIfUserExists checks the user information by email address and return accordinglyy*/
func (s *UserService) CheckIfUserExists(email string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error finding the user for this email")
	}

	if user == nil {
		return nil, fmt.Errorf("no user found using this email: %s", email)
	}

	return user, nil
}

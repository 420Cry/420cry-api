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

// UserService provides methods for managing user-related operations,
// including interactions with the user repository and email services.
type UserService struct {
	userRepo     UserRepository.UserRepository
	emailService EmailService.EmailServiceInterface
}

// NewUserService creates &UserService
func NewUserService(userRepo UserRepository.UserRepository, emailService EmailService.EmailServiceInterface) *UserService {
	return &UserService{userRepo: userRepo, emailService: emailService}
}

// UserServiceInterface provides methods of UserService
type UserServiceInterface interface {
	CreateUser(fullname, username, email, password string) (*UserModel.User, string, error)
	AuthenticateUser(username, password string) (*UserModel.User, error)
	VerifyUserWithTokens(userToken, verifyToken string) (*UserModel.User, error)
	CheckAccountVerificationToken(token string) (*UserModel.User, error)
	CheckEmailVerificationToken(token string) (*UserModel.User, error)
	CheckUserByBothTokens(token string, verificationToken string) (*UserModel.User, error)
}

// CreateUser creates a new user with the provided fullname, username, email, and password.
// If a user with the given username or email already exists, it handles the existing user case,
// potentially refreshing the user if needed. Returns the created or refreshed user, a verification token,
// and an error if any occurred during the process.
func (service *UserService) CreateUser(fullname, username, email, password string) (*UserModel.User, string, error) {
	// Check if the user already exists
	existingUser, err := service.userRepo.FindByUsernameOrEmail(username, email)
	if err != nil {
		return nil, "", err
	}

	// Handle existing user case (unverified or verified)
	if existingUser != nil {
		refreshed, err := service.handleExistingUser(existingUser, username, email)
		if err != nil {
			return nil, "", err
		}
		if refreshed != nil {
			return refreshed, refreshed.VerificationTokens, nil
		}
	}

	newUser, err := factories.NewUser(fullname, username, email, password)
	if err != nil {
		return nil, "", err
	}

	err = service.userRepo.Save(newUser)
	if err != nil {
		return nil, "", err
	}

	var token string
	if newUser.Token != nil {
		token = *newUser.Token
	} else {
		token = ""
	}
	return newUser, token, nil
}

// AuthenticateUser attempts to authenticate a user with the provided username and password.
// It returns the authenticated user if the credentials are valid and the user is verified.
// If the user is not found, the password is invalid, or the user is not verified, an error is returned.
func (service *UserService) AuthenticateUser(username string, password string) (*UserModel.User, error) {
	user, err := service.userRepo.FindByUsername(username)
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

// CheckAccountVerificationToken checks if the provided account token is valid
func (service *UserService) CheckAccountVerificationToken(token string) (*UserModel.User, error) {
	user, err := service.userRepo.FindByAccountVerificationToken(token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid account token")
	}
	return user, nil
}

// VerifyUserWithTokens validates the token + OTP and marks the user as verified
func (service *UserService) VerifyUserWithTokens(token string, verificationToken string) (*UserModel.User, error) {
	user, err := service.CheckUserByBothTokens(token, verificationToken)
	if err != nil {
		return nil, err
	}

	user.IsVerified = true
	user.VerificationTokens = ""
	user.Token = nil

	if err := service.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to update user verification status: %w", err)
	}

	return user, nil
}

// CheckUserByBothTokens checks if the provided verification token is valid
func (service *UserService) CheckUserByBothTokens(token string, verificationToken string) (*UserModel.User, error) {
	// Find user by token
	user, err := service.userRepo.FindByVerificationToken(verificationToken)
	if err != nil {
		return nil, err
	}

	// If no user found
	if user == nil {
		return nil, fmt.Errorf("invalid verification token")
	}

	// Check if the URL token matches
	if user.Token == nil || *user.Token != token {
		return nil, fmt.Errorf("token does not match")
	}

	// Return user and no error if both tokens are valid
	return user, nil
}

// CheckEmailVerificationToken checks if the provided verification token is valid
func (service *UserService) CheckEmailVerificationToken(token string) (*UserModel.User, error) {
	// Find the user associated with the token
	user, err := service.userRepo.FindByVerificationToken(token)
	if err != nil {
		return nil, err
	}

	// If the user is not found, log and return an error
	if user == nil {
		return nil, fmt.Errorf("invalid verification token")
	}
	// Update the user's verification status and remove tokens
	user.IsVerified = true
	user.VerificationTokens = ""

	// Save the updated user to the repository
	err = service.userRepo.Save(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// handleExistingUser checks if the existing user is unverified and handles accordingly
func (service *UserService) handleExistingUser(existingUser *UserModel.User, username, email string) (*UserModel.User, error) {
	if existingUser.Username == username || existingUser.Email == email {
		if !existingUser.IsVerified {
			if time.Since(existingUser.VerificationTokenCreatedAt) > 24*time.Hour {
				newVerificationToken, err := factories.GenerateVerificationToken()
				if err != nil {
					return nil, err
				}

				existingUser.VerificationTokens = newVerificationToken
				existingUser.VerificationTokenCreatedAt = time.Now()

				err = service.userRepo.Save(existingUser)
				if err != nil {
					return nil, err
				}
				return existingUser, nil
			}
			// If token is still valid, return same user
			return existingUser, nil
		}

		return nil, fmt.Errorf("user with %s is already verified", username)
	}

	return nil, nil
}

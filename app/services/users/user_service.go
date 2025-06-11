// Package services provides business logic for user-related operations.
package services

import (
	"errors"
	"fmt"
	"time"

	core "cry-api/app/core/users"
	UserDomain "cry-api/app/domain/users"
	EmailServices "cry-api/app/services/email"
)

// UserRepository defines the methods needed for user persistence
type UserRepository interface {
	Save(user *UserDomain.User) error
	FindByUsernameOrEmail(username, email string) (*UserDomain.User, error)
	FindByVerificationToken(token string) (*UserDomain.User, error)
	FindByAccountVerificationToken(token string) (*UserDomain.User, error)
	Delete(userID int) error
}

// UserService provides operations related to users
type UserService struct {
	userRepo     core.UserRepository
	emailService *EmailServices.EmailService
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo core.UserRepository, emailService *EmailServices.EmailService) *UserService {
	return &UserService{userRepo: userRepo, emailService: emailService}
}

func (service *UserService) CreateUser(fullname, username, email, password string) (*UserDomain.User, string, error) {
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

	newUser, err := UserDomain.NewUser(fullname, username, email, password)
	if err != nil {
		return nil, "", err
	}

	err = service.userRepo.Save(newUser)
	if err != nil {
		return nil, "", err
	}

	return newUser, newUser.Token, nil
}

// handleExistingUser checks if the existing user is unverified and handles accordingly
func (service *UserService) handleExistingUser(existingUser *UserDomain.User, username, email string) (*UserDomain.User, error) {
	if existingUser.Username == username || existingUser.Email == email {
		if !existingUser.IsVerified {
			if time.Since(existingUser.VerificationTokenCreatedAt) > 24*time.Hour {
				newVerificationToken, err := UserDomain.GenerateVerificationToken()
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

// CheckUserByBothTokens checks if the provided verification token is valid
func (service *UserService) CheckUserByBothTokens(token string, verificationToken string) (*UserDomain.User, error) {
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
	if user.Token != token {
		return nil, fmt.Errorf("token does not match")
	}

	// Return user and no error if both tokens are valid
	return user, nil
}

// CheckEmailVerificationToken checks if the provided verification token is valid
func (service *UserService) CheckEmailVerificationToken(token string) (*UserDomain.User, error) {
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

// CheckAccountVerificationToken checks if the provided account token is valid
func (service *UserService) CheckAccountVerificationToken(token string) (*UserDomain.User, error) {
	user, err := service.userRepo.FindByAccountVerificationToken(token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid account token")
	}
	return user, nil
}

func (service *UserService) AuthenticateUser(username string, password string) (*UserDomain.User, error) {
	user, err := service.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid password")
	}

	if !user.IsVerified {
		return nil, errors.New("user not verified")
	}

	return user, nil
}

// VerifyUserWithTokens validates the token + OTP and marks the user as verified
func (service *UserService) VerifyUserWithTokens(token string, verificationToken string) (*UserDomain.User, error) {
	user, err := service.CheckUserByBothTokens(token, verificationToken)
	if err != nil {
		return nil, err
	}

	user.IsVerified = true
	user.VerificationTokens = ""
	user.Token = ""

	if err := service.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to update user verification status: %w", err)
	}

	return user, nil
}

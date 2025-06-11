// Package services provides business logic for user-related operations.
package services

import (
	"fmt"
	"log"

	UserDomain "cry-api/app/domain/users"
	EmailServices "cry-api/app/services/email"
)

// UserRepository defines the methods needed for user persistence
type UserRepository interface {
	Save(user *UserDomain.User) error
	FindByUsernameOrEmail(username, email string) (*UserDomain.User, error)
	FindByVerificationToken(token string) (*UserDomain.User, error)
	FindByAccountVerificationToken(token string) (*UserDomain.User, error)
}

// UserService provides operations related to users
type UserService struct {
	userRepo     UserRepository
	emailService *EmailServices.EmailService
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo UserRepository, emailService *EmailServices.EmailService) *UserService {
	return &UserService{userRepo: userRepo, emailService: emailService}
}

// CreateUser creates a new user and returns the created user and the verification token
func (service *UserService) CreateUser(fullname, username, email, password string) (*UserDomain.User, string, error) {
	// Check if the user already exists
	existingUser, err := service.userRepo.FindByUsernameOrEmail(username, email)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		return nil, "", err
	}

	if existingUser != nil {
		if existingUser.Username == username {
			return nil, "", fmt.Errorf("username is already taken")
		}
		if existingUser.Email == email {
			return nil, "", fmt.Errorf("email is already taken")
		}
	}

	newUser, err := UserDomain.NewUser(fullname, username, email, password)
	if err != nil {
		log.Printf("Error creating new user: %v", err)
		return nil, "", err
	}

	// Save the user to the repository
	err = service.userRepo.Save(newUser)
	if err != nil {
		return nil, "", err
	}

	return newUser, newUser.Token, nil
}

// CheckEmailVerificationToken checks if the provided verification token is valid
func (service *UserService) CheckEmailVerificationToken(token string) (*UserDomain.User, error) {
	// Find the user associated with the token
	user, err := service.userRepo.FindByVerificationToken(token)
	if err != nil {
		log.Printf("Error finding user by verification token: %v", err)
		return nil, err
	}

	// If the user is not found, log and return an error
	if user == nil {
		log.Printf("No user found for token: %s", token)
		return nil, fmt.Errorf("invalid verification token")
	}
	// Update the user's verification status and remove tokens
	user.IsVerified = true
	user.VerificationTokens = ""

	// Save the updated user to the repository
	err = service.userRepo.Save(user)
	if err != nil {
		log.Printf("Error updating user verification status: %v", err)
		return nil, err
	}
	return user, nil
}

// CheckAccountVerificationToken checks if the provided account token is valid
func (service *UserService) CheckAccountVerificationToken(token string) (*UserDomain.User, error) {
	user, err := service.userRepo.FindByAccountVerificationToken(token)
	if err != nil {
		log.Printf("Error finding user by account token: %v", err)
		return nil, err
	}
	if user == nil {
		log.Printf("No user found for token: %s", token)
		return nil, fmt.Errorf("invalid account token")
	}
	return user, nil
}

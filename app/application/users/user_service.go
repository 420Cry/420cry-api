package application

import (
	EmailApplication "cry-api/app/application/email"
	UserDomain "cry-api/app/domain/users"
	"fmt"
	"log"
)

// UserRepository defines the methods needed for user persistence
type UserRepository interface {
	Save(user *UserDomain.User) error
	FindByUsernameOrEmail(username, email string) (*UserDomain.User, error)
}

// UserService provides operations related to users
type UserService struct {
	userRepo     UserRepository
	emailService *EmailApplication.EmailService
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo UserRepository, emailService *EmailApplication.EmailService) *UserService {
	return &UserService{userRepo: userRepo, emailService: emailService}
}

// CreateUser creates a new user and returns the created user and the verification token
func (service *UserService) CreateUser(username, email, password string) (*UserDomain.User, string, error) {
	// Log the incoming request
	log.Printf("Creating user with username: %s, email: %s", username, email)

	// Check if the user already exists
	existingUser, err := service.userRepo.FindByUsernameOrEmail(username, email)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		return nil, "", err
	}

	if existingUser != nil {
		if existingUser.Username == username {
			log.Printf("Username %s is already taken", username)
			return nil, "", fmt.Errorf("username is already taken")
		}
		if existingUser.Email == email {
			log.Printf("Email %s is already taken", email)
			return nil, "", fmt.Errorf("email is already taken")
		}
	}

	// Create a new user in the domain layer
	log.Printf("Creating new user...")
	newUser, err := UserDomain.NewUser(username, email, password)
	if err != nil {
		log.Printf("Error creating new user: %v", err)
		return nil, "", err
	}

	// Save the user to the repository
	err = service.userRepo.Save(newUser)
	if err != nil {
		log.Printf("Error saving new user: %v", err)
		return nil, "", err
	}

	// Log the successful user creation
	log.Printf("User created successfully with ID: %d, Username: %s", newUser.ID, newUser.Username)

	// Return the verification token (this can be generated as part of the user creation)
	return newUser, newUser.Token, nil
}

package userapplication

import (
	UserDomain "cry-api/app/domain/users"
	"fmt"
)

// UserRepository defines the methods needed for user persistence
type UserRepository interface {
	Save(user *UserDomain.User) error
	FindByUsernameOrEmail(username, email string) (*UserDomain.User, error)
}

// UserService provides operations related to users
type UserService struct {
	userRepo UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser creates a new user
func (service *UserService) CreateUser(username, email, password string) (*UserDomain.User, error) {
	// Check if the user already exists
	existingUser, err := service.userRepo.FindByUsernameOrEmail(username, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		if existingUser.Username == username {
			return nil, fmt.Errorf("username is already taken")
		}
		if existingUser.Email == email {
			return nil, fmt.Errorf("email is already taken")
		}
	}

	// Create a new user in the domain layer
	newUser, err := UserDomain.NewUser(username, email, password)
	if err != nil {
		return nil, err
	}

	// Save the user to the repository
	err = service.userRepo.Save(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

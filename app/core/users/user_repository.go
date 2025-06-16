// Package usercore package provide interface of UserRepository
package usercore

import UserDomain "cry-api/app/domain/users"

// UserRepository defines the contract for user persistence and retrieval operations.
// It provides methods to save, find, and delete users using various identifiers such as UUID, email, username, and tokens.
type UserRepository interface {
	// Save persists a user to the repository.
	Save(user *UserDomain.User) error

	// FindByUUID retrieves a user by their UUID.
	FindByUUID(uuid string) (*UserDomain.User, error)

	// FindByEmail retrieves a user by their email address.
	FindByEmail(email string) (*UserDomain.User, error)

	// FindByUsernameOrEmail retrieves a user by their username or email.
	FindByUsernameOrEmail(username, email string) (*UserDomain.User, error)

	// FindByVerificationToken retrieves a user by their verification token.
	FindByVerificationToken(token string) (*UserDomain.User, error)

	// FindByAccountVerificationToken retrieves a user by their account verification token.
	FindByAccountVerificationToken(token string) (*UserDomain.User, error)

	// FindByUsername retrieves a user by their username.
	FindByUsername(username string) (*UserDomain.User, error)

	// FindByUserToken retrieves a user by their user token.
	FindByUserToken(token string) (*UserDomain.User, error)

	// Delete removes a user by their ID.
	Delete(userID int) error
}

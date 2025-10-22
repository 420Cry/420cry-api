// Package errors defines error msgs
package errors

var (
	// ErrUserNotFound returns "user not found" as error
	ErrUserNotFound = NewNotFoundError("User", "user not found")
	// ErrInvalidPassword returns "invalid password" as error
	ErrInvalidPassword = NewUnauthorizedError("invalid password")
	// ErrUserNotVerified returns "user not verified" as error
	ErrUserNotVerified = NewUnauthorizedError("user not verified")
)

// Package errors defines error msgs
package errors

import "errors"

var (
	// ErrUserNotFound returns "user not found" as error
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidPassword returns "nvalid password" as error
	ErrInvalidPassword = errors.New("invalid password")
	// ErrUserNotVerified returns "user not verified" as error
	ErrUserNotVerified = errors.New("user not verified")
)

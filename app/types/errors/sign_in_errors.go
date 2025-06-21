// Package errors defines error msgs
package errors

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotVerified = errors.New("user not verified")
)

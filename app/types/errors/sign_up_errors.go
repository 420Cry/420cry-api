// Package errors defines error msgs
package errors

// ErrUserConflict returns "user already exists" as error
var ErrUserConflict = NewConflictError("User", "User already exists")

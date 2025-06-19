// Package errors defines error msgs
package errors

import "errors"

// ErrUserConflict returns "user already exists" as error
var ErrUserConflict = errors.New("user already exists")

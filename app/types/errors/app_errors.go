// Package errors defines application-specific error types
package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ValidationError represents validation errors
type ValidationError struct {
	*AppError
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
}

// NewValidationError creates a new validation error
func NewValidationError(field, value, message string) *ValidationError {
	return &ValidationError{
		AppError: NewAppError(http.StatusBadRequest, message, ""),
		Field:    field,
		Value:    value,
	}
}

// NotFoundError represents resource not found errors
type NotFoundError struct {
	*AppError
	Resource string `json:"resource,omitempty"`
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, message string) *NotFoundError {
	if message == "" {
		message = fmt.Sprintf("%s not found", resource)
	}
	return &NotFoundError{
		AppError: NewAppError(http.StatusNotFound, message, ""),
		Resource: resource,
	}
}

// UnauthorizedError represents unauthorized access errors
type UnauthorizedError struct {
	*AppError
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *UnauthorizedError {
	if message == "" {
		message = "Unauthorized access"
	}
	return &UnauthorizedError{
		AppError: NewAppError(http.StatusUnauthorized, message, ""),
	}
}

// ConflictError represents resource conflict errors
type ConflictError struct {
	*AppError
	Resource string `json:"resource,omitempty"`
}

// NewConflictError creates a new conflict error
func NewConflictError(resource, message string) *ConflictError {
	if message == "" {
		message = fmt.Sprintf("%s already exists", resource)
	}
	return &ConflictError{
		AppError: NewAppError(http.StatusConflict, message, ""),
		Resource: resource,
	}
}

// InternalServerError represents internal server errors
type InternalServerError struct {
	*AppError
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string) *InternalServerError {
	if message == "" {
		message = "Internal server error"
	}
	return &InternalServerError{
		AppError: NewAppError(http.StatusInternalServerError, message, ""),
	}
}

// Predefined common errors
var (
	ErrInvalidJSON     = NewAppError(http.StatusBadRequest, "Invalid JSON format", "")
	ErrInvalidInput    = NewAppError(http.StatusBadRequest, "Invalid input", "")
	ErrDatabaseError   = NewInternalServerError("Database operation failed")
	ErrEmailSendFailed = NewInternalServerError("Failed to send email")
	ErrTokenGeneration = NewInternalServerError("Failed to generate token")
	ErrTokenValidation = NewUnauthorizedError("Invalid or expired token")
)

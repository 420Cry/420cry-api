// Package validators provides request validation functionality
package validators

import (
	app_errors "cry-api/app/types/errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

// UserSignupValidator validates user signup requests
type UserSignupValidator struct {
	Fullname string `json:"fullname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserSigninValidator validates user signin requests
type UserSigninValidator struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResetPasswordValidator validates password reset requests
type UserResetPasswordValidator struct {
	Email string `json:"email" binding:"required,email"`
}

// UserVerifyResetPasswordValidator validates reset password verification requests
type UserVerifyResetPasswordValidator struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ValidateUserSignup validates user signup data
func ValidateUserSignup(c *gin.Context) (*UserSignupValidator, error) {
	var input UserSignupValidator
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, app_errors.ErrInvalidJSON
	}

	// Validate fullname
	if err := validateFullname(input.Fullname); err != nil {
		return nil, err
	}

	// Validate username
	if err := validateUsername(input.Username); err != nil {
		return nil, err
	}

	// Validate email
	if err := validateEmail(input.Email); err != nil {
		return nil, err
	}

	// Validate password
	if err := validatePassword(input.Password); err != nil {
		return nil, err
	}

	return &input, nil
}

// ValidateUserSignin validates user signin data
func ValidateUserSignin(c *gin.Context) (*UserSigninValidator, error) {
	var input UserSigninValidator
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, app_errors.ErrInvalidJSON
	}

	// Validate username
	if err := validateUsername(input.Username); err != nil {
		return nil, err
	}

	// Validate password
	if err := validatePassword(input.Password); err != nil {
		return nil, err
	}

	return &input, nil
}

// ValidateResetPassword validates password reset request
func ValidateResetPassword(c *gin.Context) (*UserResetPasswordValidator, error) {
	var input UserResetPasswordValidator
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, app_errors.ErrInvalidJSON
	}

	// Validate email
	if err := validateEmail(input.Email); err != nil {
		return nil, err
	}

	return &input, nil
}

// ValidateVerifyResetPassword validates reset password verification
func ValidateVerifyResetPassword(c *gin.Context) (*UserVerifyResetPasswordValidator, error) {
	var input UserVerifyResetPasswordValidator
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, app_errors.ErrInvalidJSON
	}

	// Validate token
	if err := validateToken(input.Token); err != nil {
		return nil, err
	}

	// Validate password
	if err := validatePassword(input.Password); err != nil {
		return nil, err
	}

	return &input, nil
}

// validateFullname validates the fullname field
func validateFullname(fullname string) error {
	fullname = strings.TrimSpace(fullname)

	if len(fullname) < 2 {
		return app_errors.NewValidationError("fullname", fullname, "Fullname must be at least 2 characters long")
	}

	if len(fullname) > 100 {
		return app_errors.NewValidationError("fullname", fullname, "Fullname must not exceed 100 characters")
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	validNameRegex := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
	if !validNameRegex.MatchString(fullname) {
		return app_errors.NewValidationError("fullname", fullname, "Fullname contains invalid characters")
	}

	return nil
}

// validateUsername validates the username field
func validateUsername(username string) error {
	username = strings.TrimSpace(username)

	if len(username) < 3 {
		return app_errors.NewValidationError("username", username, "Username must be at least 3 characters long")
	}

	if len(username) > 50 {
		return app_errors.NewValidationError("username", username, "Username must not exceed 50 characters")
	}

	// Check for alphanumeric characters and underscores only
	validUsernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsernameRegex.MatchString(username) {
		return app_errors.NewValidationError("username", username, "Username can only contain letters, numbers, and underscores")
	}

	return nil
}

// validateEmail validates the email field
func validateEmail(email string) error {
	email = strings.TrimSpace(email)

	if len(email) == 0 {
		return app_errors.NewValidationError("email", email, "Email is required")
	}

	if len(email) > 254 {
		return app_errors.NewValidationError("email", email, "Email is too long")
	}

	// Basic email regex validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return app_errors.NewValidationError("email", email, "Invalid email format")
	}

	return nil
}

// validatePassword validates the password field
func validatePassword(password string) error {
	if len(password) < 8 {
		return app_errors.NewValidationError("password", "", "Password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return app_errors.NewValidationError("password", "", "Password must not exceed 128 characters")
	}

	// Check for at least one uppercase letter
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return app_errors.NewValidationError("password", "", "Password must contain at least one uppercase letter")
	}

	if !hasLower {
		return app_errors.NewValidationError("password", "", "Password must contain at least one lowercase letter")
	}

	if !hasDigit {
		return app_errors.NewValidationError("password", "", "Password must contain at least one number")
	}

	if !hasSpecial {
		return app_errors.NewValidationError("password", "", "Password must contain at least one special character")
	}

	return nil
}

// validateToken validates the token field
func validateToken(token string) error {
	token = strings.TrimSpace(token)

	if len(token) == 0 {
		return app_errors.NewValidationError("token", "", "Token is required")
	}

	if len(token) < 6 {
		return app_errors.NewValidationError("token", "", "Token is too short")
	}

	return nil
}

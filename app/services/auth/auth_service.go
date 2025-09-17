package services

import (
	"fmt"

	UserModel "cry-api/app/models"
	UserRepository "cry-api/app/repositories"
	TwoFactorService "cry-api/app/services/2fa"
	SignInError "cry-api/app/types/errors"
)

// AuthService handles user authentication.
type AuthService struct {
	userRepo        UserRepository.UserRepository
	passwordService PasswordServiceInterface
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(userRepo UserRepository.UserRepository, passwordService PasswordServiceInterface) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		passwordService: passwordService,
	}
}

// AuthServiceInterface defines the contract for user service methods.
type AuthServiceInterface interface {
	AuthenticateUser(username, password string) (*UserModel.User, error)
	SaveTOTPSecret(userUUID, secret string) error
	VerifyOTP(secret string, otp string) (bool, error)
}

// AuthenticateUser verifies username and password, and checks if user is verified.
func (s *AuthService) AuthenticateUser(username, password string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, SignInError.ErrUserNotFound
	}
	if err := s.passwordService.CheckPassword(user.Password, password); err != nil {
		return nil, SignInError.ErrInvalidPassword
	}
	if !user.IsVerified {
		return nil, SignInError.ErrUserNotVerified
	}
	return user, nil
}

// SaveTOTPSecret saves the TOTP secret for the user identified by UUID.
func (s *AuthService) SaveTOTPSecret(userUUID, secret string) error {
	// Get user by UUID
	user, err := s.userRepo.FindByUUID(userUUID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Set the TOTP secret
	user.TwoFASecret = &secret

	// Persist the change
	err = s.userRepo.Save(user)
	if err != nil {
		return fmt.Errorf("failed to save TOTP secret: %w", err)
	}

	return nil
}

// VerifyOTP verifies the OTP.
func (s *AuthService) VerifyOTP(secret string, otp string) (bool, error) {
	isValid := TwoFactorService.VerifyTOTP(secret, otp)
	if !isValid {
		return false, fmt.Errorf("invalid OTP token")
	}

	return true, nil
}

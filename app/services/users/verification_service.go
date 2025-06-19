package services

import (
	"fmt"

	UserModel "cry-api/app/models"
	UserRepository "cry-api/app/repositories"
)

// VerificationService handles verification token checks and user verification.
type VerificationService struct {
	userRepo UserRepository.UserRepository
}

// NewVerificationService creates a new VerificationService instance.
func NewVerificationService(userRepo UserRepository.UserRepository) *VerificationService {
	return &VerificationService{userRepo: userRepo}
}

// VerificationServiceInterface defines the contract for VerificationServiceInterface methods.
type VerificationServiceInterface interface {
	VerifyUserWithTokens(userToken, verifyToken string) (*UserModel.User, error)
	CheckAccountVerificationToken(token string) (*UserModel.User, error)
	CheckEmailVerificationToken(token string) (*UserModel.User, error)
	CheckUserByBothTokens(token string, verificationToken string) (*UserModel.User, error)
}

// VerifyUserWithTokens validates both tokens and marks user as verified.
func (s *VerificationService) VerifyUserWithTokens(token, verificationToken string) (*UserModel.User, error) {
	user, err := s.CheckUserByBothTokens(token, verificationToken)
	if err != nil {
		return nil, err
	}

	user.IsVerified = true
	user.VerificationTokens = ""
	user.AccountVerificationToken = nil

	if err := s.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to update user verification status: %w", err)
	}

	return user, nil
}

// CheckAccountVerificationToken validates an account token and returns the user.
func (s *VerificationService) CheckAccountVerificationToken(token string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByAccountVerificationToken(token)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid account token")
	}
	return user, nil
}

// CheckUserByBothTokens verifies both the URL token and verification token.
func (s *VerificationService) CheckUserByBothTokens(emailVerificationToken, verificationToken string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByVerificationToken(verificationToken)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid verification token")
	}
	if user.AccountVerificationToken == nil || *user.AccountVerificationToken != emailVerificationToken {
		return nil, fmt.Errorf("token does not match")
	}
	return user, nil
}

// CheckEmailVerificationToken verifies the email token and marks user as verified.
func (s *VerificationService) CheckEmailVerificationToken(emailVerificationToken string) (*UserModel.User, error) {
	user, err := s.userRepo.FindByVerificationToken(emailVerificationToken)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid verification token")
	}

	user.IsVerified = true
	user.VerificationTokens = ""

	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}
	return user, nil
}

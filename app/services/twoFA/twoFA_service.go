package service

import (
	"context"
	core "cry-api/app/core/twoFA"
	domain "cry-api/app/domain/users"
	"errors"
	"fmt"

	"github.com/pquerna/otp/totp"
)

// UserRepo defines DB methods your service needs (already have this somewhere)
type UserRepo interface {
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update2FA(ctx context.Context, id int, secret string, enabled bool) error
}

type TwoFAService struct {
	userRepo UserRepo
	issuer   string
}

func NewTwoFAService(userRepo UserRepo, issuer string) core.TwoFAService {
	return &TwoFAService{
		userRepo: userRepo,
		issuer:   issuer,
	}
}

func (s *TwoFAService) GenerateSecret(ctx context.Context, userID int, username, issuer string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: username,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate 2fa key: %w", err)
	}

	secret := key.Secret()
	otpURL := key.URL()

	return secret, otpURL, nil
}

func (s *TwoFAService) VerifyCode(ctx context.Context, userID int, code string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}
	if user.TwoFASecret == nil || !user.TwoFAEnabled {
		return false, errors.New("2fa not enabled")
	}

	return totp.Validate(code, *user.TwoFASecret), nil
}

func (s *TwoFAService) Enable2FA(ctx context.Context, userID int, secret string) error {
	return s.userRepo.Update2FA(ctx, userID, secret, true)
}

func (s *TwoFAService) Disable2FA(ctx context.Context, userID int) error {
	return s.userRepo.Update2FA(ctx, userID, "", false)
}

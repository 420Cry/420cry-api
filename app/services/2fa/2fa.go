// Package services provides 2fa logic
package services

import (
	"encoding/base64"
	"fmt"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// TwoFactorServiceInterface provides methods for 2fa
type TwoFactorServiceInterface interface {
	GenerateOtpauthURL(userEmail, secret string) string
	GenerateQRCodeBase64(url string) (string, error)
	GenerateTOTP(userEmail string) (string, string, error)
}

// TwoFactorService implements TwoFactorServiceInterface by wrapping package functions.
type TwoFactorService struct{}

// NewTwoFactorService creates a new instance of TwoFactorService.
func NewTwoFactorService() *TwoFactorService {
	return &TwoFactorService{}
}

// GenerateTOTP creates a new TOTP key and returns the secret and otpauth URL.
func (s *TwoFactorService) GenerateTOTP(userEmail string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "420CRY",
		AccountName: userEmail,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// GenerateOtpauthURL builds the otpauth URL from an existing secret and user email.
func (s *TwoFactorService) GenerateOtpauthURL(userEmail, secret string) string {
	issuer := "420CRY"
	return fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer,
		userEmail,
		secret,
		issuer,
	)
}

// GenerateQRCodeBase64 encodes the otpauth URL into a base64 QR code image.
func (s *TwoFactorService) GenerateQRCodeBase64(url string) (string, error) {
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	base64Image := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	return base64Image, nil
}

// VerifyTOTP verifies the user-provided token against the stored secret.
// This can remain a package function or be added as a method if you want.
func VerifyTOTP(secret string, token string) bool {
	return totp.Validate(token, secret)
}

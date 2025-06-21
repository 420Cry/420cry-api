// Package services provides business logic for handling 2fa operations.
package services

import (
	"encoding/base64"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// VerifyTOTP verifies the user-provided token against the stored secret.
// Returns true if valid, false otherwise.
func VerifyTOTP(secret string, token string) bool {
	return totp.Validate(token, secret)
}

// GenerateTOTP returns OTP for 2fa app
func GenerateTOTP(userEmail string) (secret, otpAuthURL string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "420CRY",
		AccountName: userEmail,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// GenerateQRCodeBase64 returns QRcode string
func GenerateQRCodeBase64(url string) (string, error) {
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	base64Image := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	return base64Image, nil
}

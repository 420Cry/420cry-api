package services

import (
	"encoding/base64"
	"fmt"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// VerifyTOTP verifies the user-provided token against the stored secret.
// Returns true if valid, false otherwise.
func VerifyTOTP(secret string, token string) bool {
	return totp.Validate(token, secret)
}

// GenerateTOTP creates a new TOTP key and returns the secret and otpauth URL.
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

// GenerateOtpauthURL builds the otpauth URL from an existing secret and user email.
func GenerateOtpauthURL(userEmail, secret string) string {
	issuer := "420CRY"
	// Construct the URL manually following the otpauth URI scheme
	return fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer,
		userEmail,
		secret,
		issuer,
	)
}

// GenerateQRCodeBase64 encodes the otpauth URL into a base64 QR code image.
func GenerateQRCodeBase64(url string) (string, error) {
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	base64Image := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	return base64Image, nil
}

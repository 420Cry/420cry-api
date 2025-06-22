package tests

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	services "cry-api/app/services/2fa"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
)

func TestVerifyTOTP(t *testing.T) {
	twoFactorService := services.NewTwoFactorService()

	secret, _, err := twoFactorService.GenerateTOTP("testuser@example.com")
	assert.NoError(t, err)

	token, err := totp.GenerateCode(secret, time.Now())
	assert.NoError(t, err)

	// Valid token should return true
	valid := services.VerifyTOTP(secret, token) // This remains package function
	assert.True(t, valid)

	// Invalid token should return false
	invalid := services.VerifyTOTP(secret, "123456")
	assert.False(t, invalid)
}

func TestGenerateTOTP(t *testing.T) {
	twoFactorService := services.NewTwoFactorService()

	secret, url, err := twoFactorService.GenerateTOTP("user@example.com")

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
	assert.NotEmpty(t, url)
	assert.True(t, strings.HasPrefix(url, "otpauth://totp/420CRY:user@example.com"))

	// The secret should be base32 encoded (all uppercase letters/numbers)
	for _, c := range secret {
		assert.Contains(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567", string(c))
	}
}

func TestGenerateOtpauthURL(t *testing.T) {
	twoFactorService := services.NewTwoFactorService()

	userEmail := "user@example.com"
	secret := "S3CR3T123"

	url := twoFactorService.GenerateOtpauthURL(userEmail, secret)
	expectedPrefix := "otpauth://totp/420CRY:user@example.com?secret=S3CR3T123&issuer=420CRY"

	assert.Equal(t, expectedPrefix, url)
}

func TestGenerateQRCodeBase64(t *testing.T) {
	twoFactorService := services.NewTwoFactorService()

	url := "otpauth://totp/420CRY:user@example.com?secret=S3CR3T123&issuer=420CRY"

	qrCode, err := twoFactorService.GenerateQRCodeBase64(url)
	assert.NoError(t, err)
	assert.NotEmpty(t, qrCode)

	assert.True(t, strings.HasPrefix(qrCode, "data:image/png;base64,"))

	b64data := strings.TrimPrefix(qrCode, "data:image/png;base64,")
	decoded, err := base64.StdEncoding.DecodeString(b64data)
	assert.NoError(t, err)
	assert.Greater(t, len(decoded), 0)
}

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
	secret, _, err := services.GenerateTOTP("testuser@example.com")
	assert.NoError(t, err)
	token, err := totp.GenerateCode(secret, time.Now())
	assert.NoError(t, err)

	// Valid token should return true
	valid := services.VerifyTOTP(secret, token)
	assert.True(t, valid)

	// Invalid token should return false
	invalid := services.VerifyTOTP(secret, "123456")
	assert.False(t, invalid)
}

func TestGenerateTOTP(t *testing.T) {
	secret, url, err := services.GenerateTOTP("user@example.com")

	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
	assert.NotEmpty(t, url)
	assert.True(t, strings.HasPrefix(url, "otpauth://totp/420CRY:user@example.com"))

	// The secret should be base32 encoded (all uppercase letters/numbers)
	// Just check it contains only valid base32 chars A-Z and 2-7
	for _, c := range secret {
		assert.Contains(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567", string(c))
	}
}

func TestGenerateOtpauthURL(t *testing.T) {
	userEmail := "user@example.com"
	secret := "S3CR3T123"

	url := services.GenerateOtpauthURL(userEmail, secret)
	expectedPrefix := "otpauth://totp/420CRY:user@example.com?secret=S3CR3T123&issuer=420CRY"

	assert.Equal(t, expectedPrefix, url)
}

func TestGenerateQRCodeBase64(t *testing.T) {
	url := "otpauth://totp/420CRY:user@example.com?secret=S3CR3T123&issuer=420CRY"

	qrCode, err := services.GenerateQRCodeBase64(url)
	assert.NoError(t, err)
	assert.NotEmpty(t, qrCode)

	// Should start with correct prefix
	assert.True(t, strings.HasPrefix(qrCode, "data:image/png;base64,"))

	// Check that base64 decoding works without error
	b64data := strings.TrimPrefix(qrCode, "data:image/png;base64,")
	decoded, err := base64.StdEncoding.DecodeString(b64data)
	assert.NoError(t, err)
	assert.Greater(t, len(decoded), 0)
}

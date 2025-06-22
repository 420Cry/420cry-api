// Package tests provides tests for 2fa
package tests

import (
	"strings"
	"testing"

	TwoFactorService "cry-api/app/services/2fa"
)

func TestGenerateTOTP(t *testing.T) {
	secret, otpAuthURL, err := TwoFactorService.GenerateTOTP("user@example.com")
	if err != nil {
		t.Fatalf("GenerateTOTP returned error: %v", err)
	}

	if secret == "" {
		t.Error("GenerateTOTP returned empty secret")
	}
	if otpAuthURL == "" {
		t.Error("GenerateTOTP returned empty otpAuthURL")
	}
	if !strings.HasPrefix(otpAuthURL, "otpauth://") {
		t.Errorf("otpAuthURL does not start with otpauth://, got: %s", otpAuthURL)
	}
}

func TestGenerateQRCodeBase64(t *testing.T) {
	// First generate a valid otpauth URL for testing
	_, otpAuthURL, err := TwoFactorService.GenerateTOTP("user@example.com")
	if err != nil {
		t.Fatalf("GenerateTOTP returned error: %v", err)
	}

	qrCodeBase64, err := TwoFactorService.GenerateQRCodeBase64(otpAuthURL)
	if err != nil {
		t.Fatalf("GenerateQRCodeBase64 returned error: %v", err)
	}

	if qrCodeBase64 == "" {
		t.Error("GenerateQRCodeBase64 returned empty string")
	}

	if !strings.HasPrefix(qrCodeBase64, "data:image/png;base64,") {
		t.Errorf("QR code base64 string does not have expected prefix, got: %s", qrCodeBase64[:30])
	}
}

func TestGenerateQRCodeBase64_AnyInput(t *testing.T) {
	inputs := []string{
		"validstring",
		"\x00\x01\x02",
		"https://example.com",
	}

	for _, input := range inputs {
		qrCodeBase64, err := TwoFactorService.GenerateQRCodeBase64(input)
		if err != nil {
			t.Errorf("GenerateQRCodeBase64 returned error for input %q: %v", input, err)
		}
		if !strings.HasPrefix(qrCodeBase64, "data:image/png;base64,") {
			t.Errorf("QR code base64 string does not have expected prefix for input %q", input)
		}
		if len(qrCodeBase64) < 30 {
			t.Errorf("QR code base64 string seems too short for input %q", input)
		}
	}

	_, err := TwoFactorService.GenerateQRCodeBase64("")
	if err == nil {
		t.Error("GenerateQRCodeBase64 did not return error for empty input")
	}
}

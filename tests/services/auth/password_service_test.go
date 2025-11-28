package tests

import (
	"testing"

	PasswordService "cry-api/app/services/auth/password"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
	ps := PasswordService.NewPasswordService()

	password := "mysecretpassword"
	hashed, err := ps.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed, "hashed password should not match plain password")
}

func TestCheckPassword_Success(t *testing.T) {
	ps := PasswordService.NewPasswordService()

	password := "mypassword"
	hashed, err := ps.HashPassword(password)
	assert.NoError(t, err)

	err = ps.CheckPassword(hashed, password)
	assert.NoError(t, err, "passwords should match")
}

func TestCheckPassword_Failure(t *testing.T) {
	ps := PasswordService.NewPasswordService()

	password := "mypassword"
	hashed, err := ps.HashPassword(password)
	assert.NoError(t, err)

	// Provide wrong password to check failure
	err = ps.CheckPassword(hashed, "wrongpassword")
	assert.Error(t, err, "passwords should not match")
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	ps := PasswordService.NewPasswordService()

	// Hashing empty password should still work, bcrypt allows it
	hashed, err := ps.HashPassword("")
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)

	// Check password empty vs hashed
	err = ps.CheckPassword(hashed, "")
	assert.NoError(t, err)
}

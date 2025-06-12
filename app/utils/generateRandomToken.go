package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

/* GenerateRandomToken (Currently using for generate reset password link params*/
func GenerateRandomToken() (string, error) {
	tokenLength := 16

	b := make([]byte, tokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("cannot generate random token")
	}

	token := hex.EncodeToString(b)
	return token, nil
}

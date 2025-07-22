package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

type EncryptService struct {
}

type EncryptServiceInterface interface {
}

func NewEncryptService() *EncryptService {
	return &EncryptService{}
}

// CreateHashKey hashes the secret passphrase with salt, returns hashkey and salt
func (s *EncryptService) CreateHashKey(passphrase string) ([]byte, []byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, nil, err
	}

	hashKey, err := scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32)

	if err != nil {
		return nil, nil, err
	}

	return hashKey, salt, nil
}

func (s *EncryptService) EncryptToken(token []byte, secretKey string) (string, error) {
	hashKey, salt, err := s.CreateHashKey(secretKey)

	if err != nil {
		return "", err
	}

	aesBlock, err := aes.NewCipher(hashKey)

	if err != nil {
		fmt.Println(err)
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)

	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, gcmInstance.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipheredText := gcmInstance.Seal(nonce, nonce, token, nil)

	final := append(salt, cipheredText...)

	return hex.EncodeToString(final), nil
}

func (s *EncryptService) DecryptToken(encryptedHex string, passphrase string) (string, error) {

	encryptedData, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	if len(encryptedData) < 16 {
		return "", fmt.Errorf("invalid encrypted data")
	}
	salt := encryptedData[:16]
	ciphertext := encryptedData[16:]

	key, err := scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipherTextOnly := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, cipherTextOnly, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

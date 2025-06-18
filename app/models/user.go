// Package models contains the user models
package models

import (
	"time"
)

// User represents a user entity in the system
type User struct {
	ID                         int       `json:"id"`
	UUID                       string    `json:"uuid" gorm:"unique;not null"`
	Username                   string    `json:"username" gorm:"unique;not null"`
	Email                      string    `json:"email" gorm:"unique;not null"`
	Fullname                   string    `json:"fullname"`
	Password                   string    `json:"-" gorm:"not null"`
	AccountVerificationToken   *string   `json:"account_verification_token,omitempty" gorm:"unique"`
	VerificationTokens         string    `json:"verification_tokens,omitempty" gorm:"size:6"`
	VerificationTokenCreatedAt time.Time `json:"verification_token_created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	IsVerified                 bool      `json:"is_verified" gorm:"not null;default:false"`
	CreatedAt                  time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt                  time.Time `json:"updated_at" gorm:"type:timestamp;default:NULL;autoUpdateTime"`
	TwoFASecret                *string   `db:"two_fa_secret"`
	TwoFAEnabled               bool      `db:"two_fa_enabled"`
}

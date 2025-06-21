// Package models contains the user models
package models

import (
	"time"
)

// User represents a user entity in the system
type User struct {
	ID                          int        `json:"id"`
	UUID                        string     `json:"uuid" gorm:"unique;not null"`
	Username                    string     `json:"username" gorm:"unique;not null"`
	Email                       string     `json:"email" gorm:"unique;not null"`
	Fullname                    string     `json:"fullname"`
	Password                    string     `json:"-" gorm:"not null"`
	AccountVerificationToken    *string    `json:"account_verification_token,omitempty" gorm:"unique"`
	ResetPasswordToken          string     `json:"reset_password_token,omitempty" gorm:"unique"`
	ResetPasswordTokenCreatedAt *time.Time `json:"reset_password_token_created_at" gorm:"type:timestamp"`
	VerificationTokens          string     `json:"verification_tokens,omitempty" gorm:"size:6"`
	VerificationTokenCreatedAt  time.Time  `json:"verification_token_created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	IsVerified                  bool       `json:"is_verified" gorm:"not null;default:false"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt                   time.Time  `json:"updated_at" gorm:"type:timestamp;default:NULL;autoUpdateTime"`
	TwoFASecret                 *string    `json:"two_fa_secret,omitempty" gorm:"column:two_fa_secret"`
	TwoFAEnabled                bool       `json:"two_fa_enabled" gorm:"column:two_fa_enabled;not null;default:false"`
}

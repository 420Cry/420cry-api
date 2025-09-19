// Package models contains the user models
package models

import (
	"time"
)

// User represents a user entity in the system
type User struct {
	ID           int       `json:"id"`
	UUID         string    `json:"uuid" gorm:"unique;not null"`
	Username     string    `json:"username" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Fullname     string    `json:"fullname"`
	Password     string    `json:"-" gorm:"not null"`
	IsVerified   bool      `json:"is_verified" gorm:"not null;default:false"`
	TwoFASecret  *string   `json:"two_fa_secret,omitempty" gorm:"column:two_fa_secret"`
	TwoFAEnabled bool      `json:"two_fa_enabled" gorm:"not null;default:false"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp;default:NULL;autoUpdateTime"`

	// Relations
	Tokens []UserToken `json:"tokens" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

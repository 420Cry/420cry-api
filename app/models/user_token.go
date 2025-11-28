package models

import (
	"time"
)

// UserToken represents a user_token entity in the system
// (used for OTPs, verification links, password resets, etc.)
type UserToken struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id" gorm:"not null;index"`
	Token     string     `json:"token" gorm:"not null"`
	Purpose   string     `json:"purpose" gorm:"not null"` // e.g. "account_verification", "password_reset", "login_otp", "transaction"
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"type:timestamp;not null"`
	Consumed  bool       `json:"consumed" gorm:"not null;default:false"`
	UsedAt    *time.Time `json:"used_at"`
}

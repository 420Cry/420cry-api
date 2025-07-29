package models

import (
	"time"

	"gorm.io/gorm"
)

type Oauth_Accounts struct {
	gorm.Model
	UserId       int       `json:"user_id"`
	User         User      `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:CASCADE"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Provider     string    `json:"provider" gorm:"type:VARCHAR(50);not null;uniqueIndex:idx_provider_providerId"`
	ProviderId   string    `json:"providerId" gorm:"type:VARCHAR(100);not null;uniqueIndex:idx_provider_providerId"`
	AccessToken  string    `json:"access_token" gorm:"type:TEXT;not null"`
	RefreshToken string    `json:"refresh_token" gorm:"type:TEXT;not null"`
	TokenExpiry  time.Time `json:"token_expiry" gorm:"type:timestamp"`
}

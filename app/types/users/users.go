package types

import "time"

type User struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID        string    `json:"uuid" gorm:"unique;not null"`
	Username    string    `json:"username" gorm:"unique;not null"`
	Fullname    string    `json:"fullname"`
	Password    string    `json:"-" gorm:"not null"`
	SignupToken string    `json:"signup_token,omitempty" gorm:"unique"`
	IsVerified  bool      `json:"is_verified" gorm:"not null;default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;default:NULL;autoUpdateTime"`
}

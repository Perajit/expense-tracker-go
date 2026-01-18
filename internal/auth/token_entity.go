package auth

import "time"

type TokenEntity struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	TokenID   string    `gorm:"unique;not null"`
	IsRevoked bool      `gorm:"default:false"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func (TokenEntity) TableName() string {
	return "auth_tokens"
}

package user

import (
	"gorm.io/gorm"
)

type UserEntity struct {
	gorm.Model
	Username string `gorm:"not null;uniqueIndex:idx_users_username"`
	Password string `gorm:"not null"`
	Email    string `gorm:"not null"`
}

func (UserEntity) TableName() string {
	return "users"
}

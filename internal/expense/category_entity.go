package expense

import (
	"gorm.io/gorm"
)

type CategoryEntity struct {
	gorm.Model
	UserID    uint   `gorm:"not null;uniqueIndex:idx_categories_user_name;default:0"`
	Name      string `gorm:"not null;uniqueIndex:idx_categories_user_name"`
	IsDefault bool   `gorm:"index;default:false"`
}

func (CategoryEntity) TableName() string {
	return "expense_categories"
}

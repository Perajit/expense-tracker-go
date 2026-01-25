package expense

import "gorm.io/gorm"

type TagEntity struct {
	gorm.Model
	UserID uint   `gorm:"not null;uniqueIndex:idx_tags_user_name"`
	Name   string `gorm:"not null;uniqueIndex:idx_tags_user_name"`
}

func (TagEntity) TableName() string {
	return "expense_tags"
}

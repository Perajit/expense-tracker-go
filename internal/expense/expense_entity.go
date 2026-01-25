package expense

import (
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ExpenseEntity struct {
	gorm.Model
	UserID     uint            `gorm:"not null;index:idx_expenses_user_date"`
	Date       int64           `gorm:"not null;index:idx_expenses_user_date"`
	Amount     decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	User       user.UserEntity `gorm:"foreignKey:UserID"`
	Note       string          `gorm:"type:text"`
	CategoryID uint            `gorm:"not null;index:idx_expenses_category"`
	Category   CategoryEntity  `gorm:"foreignKey:CategoryID"`
	Tags       []TagEntity     `gorm:"many2many:expenses_tags;"`
}

func (ExpenseEntity) TableName() string {
	return "expenses"
}

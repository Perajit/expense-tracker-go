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
	Note       string          `gorm:"type:text"`
	CategoryID uint            `gorm:"not null;index:idx_expenses_category"`
	User       user.UserEntity `gorm:"foreignKey:UserID"`
	Category   CategoryEntity  `gorm:"foreignKey:CategoryID"`
	Tags       []TagEntity     `gorm:"many2many:expense_tag_links;constraint:OnDelete:CASCADE"`
}

func (ExpenseEntity) TableName() string {
	return "expenses"
}

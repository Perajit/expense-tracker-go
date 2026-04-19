package testutil

import (
	"log"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDB() *gorm.DB {
	dialector := sqlite.Open("file::memory:")
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = gormDB.AutoMigrate(
		&user.UserEntity{},
		&expense.CategoryEntity{},
		&expense.TagEntity{},
		&expense.ExpenseEntity{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return gormDB
}

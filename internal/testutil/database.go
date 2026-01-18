package testutil

import (
	"github.com/DATA-DOG/go-sqlmock"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	db, mockSql, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{Conn: db})
	gormDB, _ := gorm.Open(dialector, &gorm.Config{})

	mockSql.ExpectBegin()
	mockSql.ExpectCommit()

	return gormDB
}

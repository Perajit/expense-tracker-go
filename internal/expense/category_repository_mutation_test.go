package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCategoryRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11

		createdEntity := expense.CategoryEntity{UserID: userID, Name: "cat1"}

		db := testutil.SetupDB()

		repo := expense.NewCategoryRepository(db)

		err := repo.Create(&createdEntity)

		var entity expense.CategoryEntity
		db.Model(&expense.CategoryEntity{}).First(&entity)
		assert.Equal(t, createdEntity.Name, entity.Name)
		assert.Equal(t, createdEntity.UserID, entity.UserID)
		assert.NoError(t, err)
	})
}

func TestCategoryRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11

		insertedCategories := []expense.CategoryEntity{
			{UserID: userID, Name: "cat1"},
			{UserID: userID, Name: "cat2"},
		}
		updatedEntity := expense.CategoryEntity{
			Model:  gorm.Model{ID: id},
			Name:   "new",
			UserID: userID + uint(1),
		}

		db := testutil.SetupDB()
		db.Create(&insertedCategories)

		repo := expense.NewCategoryRepository(db)

		err := repo.Update(&updatedEntity)

		var entity expense.CategoryEntity
		db.Model(&expense.CategoryEntity{}).Where("id = ?", id).First(&entity)
		assert.Equal(t, updatedEntity.Name, entity.Name)
		assert.Equal(t, updatedEntity.UserID, entity.UserID)
		assert.NoError(t, err)
	})
}

func TestCategoryRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11

		insertedCategories := []expense.CategoryEntity{
			{UserID: userID, Name: "cat1"},
			{UserID: userID, Name: "cat2"},
		}

		db := testutil.SetupDB()
		db.Create(&insertedCategories)

		repo := expense.NewCategoryRepository(db)

		err := repo.Delete(id)

		var count int64
		db.Model(&expense.CategoryEntity{}).Where("id = ?", id).Count(&count)
		assert.Equal(t, int64(0), count)
		assert.NoError(t, err)
	})
}

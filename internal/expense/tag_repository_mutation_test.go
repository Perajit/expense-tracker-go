package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTagRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11

		createdEntity := expense.TagEntity{UserID: userID, Name: "tag1"}

		db := testutil.SetupDB()

		repo := expense.NewTagRepository(db)

		err := repo.Create(&createdEntity)

		var entity expense.TagEntity
		db.Model(&expense.TagEntity{}).First(&entity)
		assert.Equal(t, createdEntity.Name, entity.Name)
		assert.Equal(t, createdEntity.UserID, entity.UserID)
		assert.NoError(t, err)
	})
}

func TestTagRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11

		insertedTags := []expense.TagEntity{
			{UserID: userID, Name: "tag1"},
			{UserID: userID, Name: "tag2"},
		}
		updatedEntity := expense.TagEntity{
			Model:  gorm.Model{ID: 1},
			Name:   "new",
			UserID: userID + uint(1),
		}

		db := testutil.SetupDB()
		db.Create(&insertedTags)

		repo := expense.NewTagRepository(db)

		err := repo.Update(&updatedEntity)

		var entity expense.TagEntity
		db.Model(&expense.TagEntity{}).Where("id = ?", id).First(&entity)
		assert.Equal(t, updatedEntity.Name, entity.Name)
		assert.Equal(t, updatedEntity.UserID, entity.UserID)
		assert.NoError(t, err)
	})
}

func TestTagRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11

		insertedTags := []expense.TagEntity{
			{UserID: userID, Name: "tag1"},
			{UserID: userID, Name: "tag2"},
		}

		db := testutil.SetupDB()
		db.Create(&insertedTags)

		repo := expense.NewTagRepository(db)

		err := repo.Delete(id)

		var count int64
		db.Model(&expense.TagEntity{}).Where("id = ?", id).Count(&count)
		assert.Equal(t, int64(0), count)
		assert.NoError(t, err)
	})
}

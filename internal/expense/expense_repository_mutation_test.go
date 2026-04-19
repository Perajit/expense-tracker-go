package expense_test

import (
	"testing"
	"time"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestExpenseRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11
		var categoryID uint = 2

		createdEntity := expense.ExpenseEntity{
			UserID:     userID,
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(100),
			Note:       "expense1",
			CategoryID: categoryID,
			Tags: []expense.TagEntity{
				{Model: gorm.Model{ID: 3}, UserID: userID, Name: "tag3"},
				{Model: gorm.Model{ID: 4}, UserID: userID, Name: "tag4"},
			},
		}

		db := testutil.SetupDB()

		repo := expense.NewExpenseRepository(db)

		err := repo.Create(&createdEntity)

		var entity expense.ExpenseEntity
		db.Model(&expense.ExpenseEntity{}).Preload("Tags").First(&entity)
		assert.Equal(t, createdEntity.UserID, entity.UserID)
		assert.Equal(t, createdEntity.Date, entity.Date)
		assert.Equal(t, createdEntity.Amount, entity.Amount)
		assert.Equal(t, createdEntity.Note, entity.Note)
		assert.Equal(t, createdEntity.CategoryID, entity.CategoryID)
		assert.Equal(t, 2, len(entity.Tags))
		assert.Equal(t, createdEntity.Tags[0].ID, entity.Tags[0].ID)
		assert.Equal(t, createdEntity.Tags[1].ID, entity.Tags[1].ID)
		assert.NoError(t, err)
	})
}

func TestExpenseRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11
		var categoryID uint = 2
		var newCategoryID uint = 3

		insertedExpenses := []expense.ExpenseEntity{
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(100),
				Note:       "expense1",
				CategoryID: categoryID,
				Tags:       []expense.TagEntity{},
			},
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(200),
				Note:       "expense2",
				CategoryID: categoryID,
				Tags:       []expense.TagEntity{},
			},
		}
		updatedEntity := expense.ExpenseEntity{
			Model:      gorm.Model{ID: id},
			UserID:     userID + uint(1),
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(110),
			Note:       "new",
			CategoryID: newCategoryID,
		}

		db := testutil.SetupDB()
		db.Create(insertedExpenses)

		repo := expense.NewExpenseRepository(db)

		err := repo.Update(&updatedEntity)

		var entity expense.ExpenseEntity
		db.Model(&expense.ExpenseEntity{}).Where("id = ?", id).First(&entity)
		assert.Equal(t, updatedEntity.UserID, entity.UserID)
		assert.Equal(t, updatedEntity.Date, entity.Date)
		assert.Equal(t, updatedEntity.Amount, entity.Amount)
		assert.Equal(t, updatedEntity.Note, entity.Note)
		assert.Equal(t, updatedEntity.CategoryID, entity.CategoryID)
		assert.NoError(t, err)
	})
}

func TestExpenseRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11
		var categoryID uint = 2

		insertedEntities := []expense.ExpenseEntity{
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(100),
				Note:       "expense1",
				CategoryID: categoryID,
				Tags:       []expense.TagEntity{},
			},
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(200),
				Note:       "expense2",
				CategoryID: categoryID,
				Tags:       []expense.TagEntity{},
			},
		}

		db := testutil.SetupDB()
		db.Create(insertedEntities)

		repo := expense.NewExpenseRepository(db)

		err := repo.Delete(id)

		var count int64
		db.Model(&expense.ExpenseEntity{}).Where("id = ?", id).Count(&count)
		assert.Equal(t, int64(0), count)
		assert.NoError(t, err)
	})
}

func TestExpenseRepository_UpdateTags(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11
		var categoryID uint = 2

		insertedExpenses := []expense.ExpenseEntity{
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(100),
				Note:       "expense1",
				CategoryID: categoryID,
				Tags:       []expense.TagEntity{},
			},
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(200),
				Note:       "expense2",
				CategoryID: categoryID,
				Tags:       []expense.TagEntity{},
			},
		}
		updatedEntity := expense.ExpenseEntity{
			Model: gorm.Model{ID: id},
		}
		updatedTags := []expense.TagEntity{
			{Model: gorm.Model{ID: 3}, UserID: userID, Name: "tag3"},
			{Model: gorm.Model{ID: 4}, UserID: userID, Name: "tag4"},
		}

		db := testutil.SetupDB()
		db.Create(&insertedExpenses)

		repo := expense.NewExpenseRepository(db)

		err := repo.UpdateTags(&updatedEntity, updatedTags)

		var entity expense.ExpenseEntity
		db.Model(&expense.ExpenseEntity{}).Preload("Tags").Where("id = ?", id).First(&entity)
		assert.Equal(t, len(updatedTags), len(entity.Tags))
		assert.Equal(t, updatedTags[0].ID, entity.Tags[0].ID)
		assert.Equal(t, updatedTags[1].ID, entity.Tags[1].ID)
		assert.NoError(t, err)
	})
}

func TestExpenseRepository_ReplaceCategoryID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11
		var categoryID uint = 2
		var newCategoryID uint = 3

		insertedExpenses := []expense.ExpenseEntity{
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(100),
				Note:       "expense1",
				CategoryID: categoryID,
				Tags:       []expense.TagEntity{},
			},
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(200),
				Note:       "expense2",
				CategoryID: categoryID,
				Tags:       []expense.TagEntity{},
			},
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(200),
				Note:       "expense3",
				CategoryID: categoryID + uint(1),
				Tags:       []expense.TagEntity{},
			},
		}

		db := testutil.SetupDB()
		db.Create(&insertedExpenses)

		repo := expense.NewExpenseRepository(db)

		err := repo.ReplaceCategoryID(categoryID, newCategoryID)

		var list []expense.ExpenseEntity
		db.Model(&expense.ExpenseEntity{}).Find(&list)
		assert.Equal(t, newCategoryID, list[0].CategoryID)
		assert.Equal(t, newCategoryID, list[1].CategoryID)
		assert.Equal(t, insertedExpenses[2].CategoryID, list[2].CategoryID)
		assert.NoError(t, err)
	})
}

func TestExpenseRepository_UnlinkTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11
		var categoryID uint = 2
		var tagID uint = 3

		insertedExpenses := []expense.ExpenseEntity{
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(100),
				Note:       "expense1",
				CategoryID: categoryID,
				Tags: []expense.TagEntity{
					{Model: gorm.Model{ID: tagID}, UserID: userID, Name: "tag3"},
					{Model: gorm.Model{ID: tagID + uint(1)}, UserID: userID, Name: "tag4"},
				},
			},
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(200),
				Note:       "expense2",
				CategoryID: categoryID,
				Tags: []expense.TagEntity{
					{Model: gorm.Model{ID: tagID}, UserID: userID, Name: "tag3"},
					{Model: gorm.Model{ID: tagID}, UserID: userID, Name: "tag3"},
				},
			},
			{
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(300),
				Note:       "expense3",
				CategoryID: categoryID + uint(1),
				Tags: []expense.TagEntity{
					{Model: gorm.Model{ID: tagID + uint(1)}, UserID: userID, Name: "tag5"},
					{Model: gorm.Model{ID: tagID + uint(2)}, UserID: userID, Name: "tag6"},
				},
			},
		}

		db := testutil.SetupDB()
		db.Create(&insertedExpenses)

		repo := expense.NewExpenseRepository(db)

		err := repo.UnlinkTag(tagID)

		var list []expense.ExpenseEntity
		db.Model(&expense.ExpenseEntity{}).Preload("Tags").Find(&list)

		// 1st item
		assert.Equal(t, 1, len(list[0].Tags))
		assert.Equal(t, insertedExpenses[0].Tags[1].ID, list[0].Tags[0].ID)

		// 2nd item
		assert.Equal(t, 0, len(list[1].Tags))

		// 3rd item
		assert.Equal(t, 2, len(list[2].Tags))
		assert.Equal(t, insertedExpenses[2].Tags[0].ID, list[2].Tags[0].ID)
		assert.Equal(t, insertedExpenses[2].Tags[1].ID, list[2].Tags[1].ID)

		assert.NoError(t, err)
	})
}

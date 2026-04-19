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

func TestExpenseRepository_GetByIDAndUser(t *testing.T) {
	var id uint = 1
	var userID uint = 11

	insertedCategories := []expense.CategoryEntity{
		{UserID: userID, Name: "cat1"},
		{UserID: userID, Name: "cat2"},
		{UserID: userID, Name: "cat3"},
	}
	insertedExpenses := []expense.ExpenseEntity{
		{
			UserID:     userID,
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(100),
			Note:       "expense1",
			CategoryID: 1,
			Tags: []expense.TagEntity{
				{Model: gorm.Model{ID: 3}, UserID: userID, Name: "tag3"},
				{Model: gorm.Model{ID: 4}, UserID: userID, Name: "tag4"},
			},
		},
		{
			UserID:     userID,
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(200),
			Note:       "expense2",
			CategoryID: 2,
			Tags: []expense.TagEntity{
				{Model: gorm.Model{ID: 5}, UserID: userID, Name: "tag5"},
			},
		},
	}

	db := testutil.SetupDB()
	db.Create(&insertedCategories)
	db.Create(&insertedExpenses)

	t.Run("exist", func(t *testing.T) {
		repo := expense.NewExpenseRepository(db)

		entity, err := repo.GetByIDAndUser(id, userID)

		assert.Equal(t, id, entity.ID)
		assert.Equal(t, insertedExpenses[0].UserID, entity.UserID)
		assert.Equal(t, insertedExpenses[0].Date, entity.Date)
		assert.Equal(t, insertedExpenses[0].Amount, entity.Amount)
		assert.Equal(t, insertedExpenses[0].Note, entity.Note)
		assert.Equal(t, insertedExpenses[0].CategoryID, entity.CategoryID)
		assert.Equal(t, insertedCategories[0].ID, entity.Category.ID)
		assert.Equal(t, insertedCategories[0].Name, entity.Category.Name)
		assert.Equal(t, 2, len(entity.Tags))
		assert.Equal(t, insertedExpenses[0].Tags[0].ID, entity.Tags[0].ID)
		assert.Equal(t, insertedExpenses[0].Tags[0].Name, entity.Tags[0].Name)
		assert.Equal(t, insertedExpenses[0].Tags[1].ID, entity.Tags[1].ID)
		assert.Equal(t, insertedExpenses[0].Tags[1].Name, entity.Tags[1].Name)
		assert.NoError(t, err)
	})

	t.Run("not exist", func(t *testing.T) {
		repo := expense.NewExpenseRepository(db)

		entity, err := repo.GetByIDAndUser(id, userID+uint(1))

		assert.Nil(t, entity)
		assert.Error(t, err)
	})
}

func TestExpenseRepository_GetByIDAndUserNoAssociation(t *testing.T) {
	var id uint = 1
	var userID uint = 11

	insertedExpenses := []expense.ExpenseEntity{
		{
			UserID:     userID,
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(100),
			Note:       "expense1",
			CategoryID: 1,
			Tags: []expense.TagEntity{
				{Model: gorm.Model{ID: 3}, UserID: userID, Name: "tag3"},
				{Model: gorm.Model{ID: 4}, UserID: userID, Name: "tag4"},
			},
		},
		{
			UserID:     userID,
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(200),
			Note:       "expense2",
			CategoryID: 2,
			Tags: []expense.TagEntity{
				{Model: gorm.Model{ID: 5}, UserID: userID, Name: "tag5"},
			},
		},
	}

	db := testutil.SetupDB()
	db.Create(&insertedExpenses)

	t.Run("exist", func(t *testing.T) {
		repo := expense.NewExpenseRepository(db)

		entity, err := repo.GetByIDAndUser(id, userID)

		assert.Equal(t, id, entity.ID)
		assert.Equal(t, insertedExpenses[0].UserID, entity.UserID)
		assert.Equal(t, insertedExpenses[0].Date, entity.Date)
		assert.Equal(t, insertedExpenses[0].Amount, entity.Amount)
		assert.Equal(t, insertedExpenses[0].Note, entity.Note)
		assert.Equal(t, insertedExpenses[0].CategoryID, entity.CategoryID)
		assert.NoError(t, err)
	})

	t.Run("not exist", func(t *testing.T) {
		repo := expense.NewExpenseRepository(db)

		entity, err := repo.GetByIDAndUser(id, userID+uint(1))

		assert.Nil(t, entity)
		assert.Error(t, err)
	})
}

func TestExpenseRepository_GetByUser(t *testing.T) {
	var userID uint = 11

	insertedExpenses := []expense.ExpenseEntity{
		{
			UserID:     userID,
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(100),
			Note:       "expense1",
			CategoryID: 1,
			Tags: []expense.TagEntity{
				{Model: gorm.Model{ID: 3}, UserID: userID, Name: "tag3"},
				{Model: gorm.Model{ID: 4}, UserID: userID, Name: "tag4"},
			},
		},
		{
			UserID:     userID,
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(200),
			Note:       "expense2",
			CategoryID: 2,
			Tags: []expense.TagEntity{
				{Model: gorm.Model{ID: 5}, UserID: userID, Name: "tag5"},
			},
		},
		{
			UserID:     userID + uint(1),
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(300),
			Note:       "expense2",
			CategoryID: 3,
			Tags:       []expense.TagEntity{},
		},
	}

	db := testutil.SetupDB()
	db.Create(&insertedExpenses)

	t.Run("exist", func(t *testing.T) {
		repo := expense.NewExpenseRepository(db)

		list, err := repo.GetByUser(userID)
		assert.Equal(t, 2, len(list))

		// 1st item
		assert.Equal(t, uint(1), list[0].ID)
		assert.Equal(t, insertedExpenses[0].UserID, list[0].UserID)
		assert.Equal(t, insertedExpenses[0].Date, list[0].Date)
		assert.Equal(t, insertedExpenses[0].Amount, list[0].Amount)
		assert.Equal(t, insertedExpenses[0].Note, list[0].Note)
		assert.Equal(t, insertedExpenses[0].CategoryID, list[0].CategoryID)
		assert.Equal(t, 2, len(list[0].Tags))
		assert.Equal(t, insertedExpenses[0].Tags[0].ID, list[0].Tags[0].ID)
		assert.Equal(t, insertedExpenses[0].Tags[0].Name, list[0].Tags[0].Name)
		assert.Equal(t, insertedExpenses[0].Tags[1].ID, list[0].Tags[1].ID)
		assert.Equal(t, insertedExpenses[0].Tags[1].Name, list[0].Tags[1].Name)

		// 2nd item
		assert.Equal(t, uint(1), list[0].ID)
		assert.Equal(t, insertedExpenses[1].UserID, list[1].UserID)
		assert.Equal(t, insertedExpenses[1].Date, list[1].Date)
		assert.Equal(t, insertedExpenses[1].Amount, list[1].Amount)
		assert.Equal(t, insertedExpenses[1].Note, list[1].Note)
		assert.Equal(t, insertedExpenses[1].CategoryID, list[1].CategoryID)
		assert.Equal(t, 1, len(list[1].Tags))
		assert.Equal(t, insertedExpenses[1].Tags[0].ID, list[1].Tags[0].ID)
		assert.Equal(t, insertedExpenses[1].Tags[0].Name, list[1].Tags[0].Name)

		assert.NoError(t, err)
	})
}

func TestExpenseRepository_IsOwner(t *testing.T) {
	var id uint = 1
	var userID uint = 11

	insertedExpenses := []expense.ExpenseEntity{
		{
			UserID:     userID,
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(100),
			Note:       "expense1",
			CategoryID: 1,
			Tags:       []expense.TagEntity{},
		},
	}

	db := testutil.SetupDB()
	db.Create(&insertedExpenses)

	t.Run("owner", func(t *testing.T) {
		repo := expense.NewExpenseRepository(db)

		// should be true
		result, err := repo.IsOwner(id, userID)

		assert.True(t, result)
		assert.NoError(t, err)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := expense.NewExpenseRepository(db)

		// should be true
		result, err := repo.IsOwner(id, userID+uint(1))

		assert.False(t, result)
		assert.NoError(t, err)
	})
}

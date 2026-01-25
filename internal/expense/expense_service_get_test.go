package expense_test

import (
	"testing"
	"time"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetExpenseByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11
		matchedEntity := &expense.ExpenseEntity{
			Model:      gorm.Model{ID: id},
			UserID:     userID,
			Date:       time.Now().Unix(),
			Amount:     decimal.NewFromInt(100),
			Note:       "expense1",
			CategoryID: 2,
			Category: expense.CategoryEntity{
				Model:  gorm.Model{ID: 2},
				UserID: userID,
				Name:   "cat2",
			},
			Tags: []expense.TagEntity{
				{Model: gorm.Model{ID: 3}, UserID: userID, Name: "tag3"},
				{Model: gorm.Model{ID: 4}, UserID: userID, Name: "tag4"},
			},
		}

		mockExpenseRepo := new(mocks.MockExpenseRepository)
		mockExpenseRepo.On("GetByIDAndUser", id, userID).Return(matchedEntity, nil).Once()

		mockCategoryService := new(mocks.MockCategoryService)

		mockTagService := new(mocks.MockTagService)

		db := testutil.SetupDB()

		service := expense.NewExpenseService(db, mockExpenseRepo, mockCategoryService, mockTagService)
		entity, err := service.GetExpenseByID(id, userID)

		assert.Equal(t, matchedEntity, entity)
		assert.NoError(t, err)
	})
}

func TestGetExpenses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11
		matchedList := []expense.ExpenseEntity{
			{
				Model:      gorm.Model{ID: 1},
				UserID:     userID,
				Date:       time.Now().Unix(),
				Amount:     decimal.NewFromInt(100),
				Note:       "expense1",
				CategoryID: 2,
				Category: expense.CategoryEntity{
					Model:  gorm.Model{ID: 2},
					UserID: userID,
					Name:   "cat2",
				},
				Tags: []expense.TagEntity{
					{Model: gorm.Model{ID: 3}, UserID: userID, Name: "tag3"},
					{Model: gorm.Model{ID: 4}, UserID: userID, Name: "tag4"},
				},
			},
			{
				Model:      gorm.Model{ID: 3},
				UserID:     userID,
				Amount:     decimal.NewFromInt(100),
				Note:       "expense3",
				CategoryID: 5,
				Category: expense.CategoryEntity{
					Model:  gorm.Model{ID: 5},
					UserID: userID,
					Name:   "cat1",
				},
				Tags: []expense.TagEntity{
					{Model: gorm.Model{ID: 6}, UserID: userID, Name: "tag6"},
				},
			},
		}

		db := testutil.SetupDB()

		mockExpenseRepo := new(mocks.MockExpenseRepository)
		mockExpenseRepo.On("GetByUser", userID).Return(matchedList, nil).Once()

		mockCategoryService := new(mocks.MockCategoryService)

		mockTagService := new(mocks.MockTagService)

		service := expense.NewExpenseService(db, mockExpenseRepo, mockCategoryService, mockTagService)
		list, err := service.GetExpenses(userID)

		assert.Equal(t, matchedList, list)
		assert.NoError(t, err)
		mockExpenseRepo.AssertExpectations(t)
	})
}

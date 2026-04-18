package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCategoryService_GetCategoryByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 1

		matchedCategory := &expense.CategoryEntity{
			Model:  gorm.Model{ID: id},
			Name:   "cat1",
			UserID: userID,
		}

		db := testutil.SetupDB()

		mockCategoryRepo := new(mocks.MockCategoryRepository)
		mockCategoryRepo.On("GetByIDAndUser", id, userID).Return(matchedCategory, nil).Once()

		mockExpenseRepo := new(mocks.MockExpenseRepository)

		service := expense.NewCategoryService(db, mockCategoryRepo, mockExpenseRepo)
		entity, err := service.GetCategoryByID(id, userID)

		assert.Equal(t, matchedCategory, entity)
		assert.Nil(t, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}

func TestCategoryService_GetCategories(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11

		matchedList := []expense.CategoryEntity{
			{Model: gorm.Model{ID: 1}, UserID: userID, Name: "cat1"},
			{Model: gorm.Model{ID: 2}, UserID: userID, Name: "cat2"},
		}

		db := testutil.SetupDB()

		mockCategoryRepo := new(mocks.MockCategoryRepository)
		mockCategoryRepo.On("GetAvailableByUser", userID).Return(matchedList, nil).Once()

		mockExpenseRepo := new(mocks.MockExpenseRepository)

		service := expense.NewCategoryService(db, mockCategoryRepo, mockExpenseRepo)
		list, err := service.GetCategories(userID)

		assert.Equal(t, matchedList, list)
		assert.NoError(t, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}

package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetCategoryByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 1
		matchedCategory := &expense.CategoryEntity{
			Model:  gorm.Model{ID: id},
			Name:   "cat1",
			UserID: userID,
		}

		mockCategoryRepo := new(mocks.MockCategoryRepository)
		mockCategoryRepo.On("GetByIDAndUser", id, &userID).Return(matchedCategory, nil).Once()

		service := expense.NewCategoryService(mockCategoryRepo)
		entity, err := service.GetCategoryByID(id, &userID)

		assert.Equal(t, matchedCategory, entity)
		assert.Nil(t, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}

func TestGetCategories(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11
		matchedList := []expense.CategoryEntity{
			{Model: gorm.Model{ID: 1}, UserID: userID, Name: "cat1"},
			{Model: gorm.Model{ID: 2}, UserID: userID, Name: "cat2"},
		}

		mockCategoryRepo := new(mocks.MockCategoryRepository)
		mockCategoryRepo.On("GetByUser", userID).Return(matchedList, nil).Once()

		service := expense.NewCategoryService(mockCategoryRepo)
		list, err := service.GetCategories(userID)

		assert.Equal(t, matchedList, list)
		assert.NoError(t, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}

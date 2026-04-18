package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryService_DeleteCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11
		var defaultID uint = 2

		db := testutil.SetupDB()

		mockCategoryRepo := new(mocks.MockCategoryRepository)
		mockCategoryRepo.On("WithTx", mock.Anything).Return(mockCategoryRepo).Once()
		mockCategoryRepo.On("IsOwner", id, userID).Return(true, nil).Once()
		mockCategoryRepo.On("GetDefaultCategoryID").Return(defaultID, nil).Once()
		mockCategoryRepo.On("Delete", id).Return(nil).Once()

		mockExpenseRepo := new(mocks.MockExpenseRepository)
		mockExpenseRepo.On("WithTx", mock.Anything).Return(mockExpenseRepo).Once()
		mockExpenseRepo.On("ReplaceCategoryID", id, defaultID).Return(nil).Once()

		service := expense.NewCategoryService(db, mockCategoryRepo, mockExpenseRepo)
		err := service.DeleteCategory(id, userID)

		assert.Nil(t, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}

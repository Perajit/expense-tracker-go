package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryService_DeleteTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11

		db := testutil.SetupDB()

		mockTagRepo := new(mocks.MockTagRepository)
		mockTagRepo.On("WithTx", mock.Anything).Return(mockTagRepo)
		mockTagRepo.On("IsOwner", id, userID).Return(true, nil).Once()
		mockTagRepo.On("Delete", id).Return(nil).Once()

		mockExpenseRepo := new(mocks.MockExpenseRepository)
		mockExpenseRepo.On("WithTx", mock.Anything).Return(mockExpenseRepo)
		mockExpenseRepo.On("UnlinkTag", id).Return(nil).Once()

		service := expense.NewTagService(db, mockTagRepo, mockExpenseRepo)
		err := service.DeleteTag(id, userID)

		assert.NoError(t, err)
		mockTagRepo.AssertExpectations(t)
	})
}

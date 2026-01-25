package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	expenseMocks "github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11

		mockCategoryRepo := new(expenseMocks.MockCategoryRepository)
		mockCategoryRepo.On("IsOwner", id, userID).Return(true, nil).Once()
		mockCategoryRepo.On("Delete", id).Return(nil).Once()

		service := expense.NewCategoryService(mockCategoryRepo)
		err := service.DeleteCategory(id, userID)

		assert.Nil(t, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}

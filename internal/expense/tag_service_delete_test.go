package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11

		mockTagRepo := new(mocks.MockTagRepository)
		mockTagRepo.On("IsOwner", id, userID).Return(true, nil).Once()
		mockTagRepo.On("Delete", id).Return(nil).Once()

		service := expense.NewTagService(mockTagRepo)
		err := service.DeleteTag(id, userID)

		assert.NoError(t, err)
		mockTagRepo.AssertExpectations(t)
	})
}

package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryService_CreateTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 1
		var newEntity *expense.TagEntity

		dto := expense.CreateTagRequest{
			Name: "tag1",
		}

		db := testutil.SetupDB()

		mockTagRepo := new(mocks.MockTagRepository)
		mockTagRepo.On("Create", mock.MatchedBy(func(e *expense.TagEntity) bool {
			if e.UserID != userID || e.Name != dto.Name {
				return false
			}
			newEntity = e
			return true
		})).Return(nil).Once()

		mockExpenseRepo := new(mocks.MockExpenseRepository)

		service := expense.NewTagService(db, mockTagRepo, mockExpenseRepo)
		entity, err := service.CreateTag(userID, dto)

		assert.Equal(t, newEntity, entity)
		assert.NoError(t, err)
		mockTagRepo.AssertExpectations(t)
	})
}

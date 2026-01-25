package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 1
		dto := expense.CreateTagRequest{
			Name: "tag1",
		}
		var newEntity *expense.TagEntity

		mockTagRepo := new(mocks.MockTagRepository)
		mockTagRepo.On("Create", mock.MatchedBy(func(e *expense.TagEntity) bool {
			if e.UserID != userID || e.Name != dto.Name {
				return false
			}
			newEntity = e
			return true
		})).Return(nil).Once()

		service := expense.NewTagService(mockTagRepo)
		entity, err := service.CreateTag(userID, dto)

		assert.Equal(t, newEntity, entity)
		assert.NoError(t, err)
		mockTagRepo.AssertExpectations(t)
	})
}

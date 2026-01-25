package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUpdateTag(t *testing.T) {
	newName := "new"

	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11
		dto := expense.UpdateTagRequest{Name: &newName}
		existingEntity := &expense.TagEntity{
			Model:  gorm.Model{ID: id},
			UserID: userID,
			Name:   *dto.Name,
		}

		mockTagRepo := new(mocks.MockTagRepository)
		mockTagRepo.On("GetByIDAndUser", id, userID).Return(existingEntity, nil).Once()
		mockTagRepo.On("Update", mock.MatchedBy(func(e *expense.TagEntity) bool {
			if e.ID != existingEntity.ID || e.UserID != existingEntity.UserID {
				return false
			}
			if e.Name != *dto.Name {
				return false
			}
			return true
		})).Return(nil).Once()

		service := expense.NewTagService(mockTagRepo)
		err := service.UpdateTag(id, userID, dto)

		assert.NoError(t, err)
		mockTagRepo.AssertExpectations(t)
	})
}

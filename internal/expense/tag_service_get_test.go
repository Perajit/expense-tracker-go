package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetTagByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11
		matchedEntity := &expense.TagEntity{
			Model:  gorm.Model{ID: id},
			UserID: userID,
			Name:   "tag1",
		}

		mockTagRepo := new(mocks.MockTagRepository)
		mockTagRepo.On("GetByIDAndUser", id, userID).Return(matchedEntity, nil).Once()

		service := expense.NewTagService(mockTagRepo)
		entity, err := service.GetTagByID(id, userID)

		assert.Equal(t, matchedEntity, entity)
		assert.NoError(t, err)
		mockTagRepo.AssertExpectations(t)
	})
}

func TestGetTagsByIDs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11
		tagIDs := []uint{1, 2}
		matchedList := []expense.TagEntity{
			{Model: gorm.Model{ID: 1}, UserID: userID, Name: "tag1"},
			{Model: gorm.Model{ID: 2}, UserID: userID, Name: "tag2"},
		}

		mockTagRepo := new(mocks.MockTagRepository)
		mockTagRepo.On("GetByIDsAndUser", tagIDs, userID).Return(matchedList, nil).Once()

		service := expense.NewTagService(mockTagRepo)
		list, err := service.GetTagsByIDs(tagIDs, userID)

		assert.Equal(t, matchedList, list)
		assert.NoError(t, err)
		mockTagRepo.AssertExpectations(t)
	})
}

func TestGetTagsByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11
		matchedList := []expense.TagEntity{
			{Model: gorm.Model{ID: 1}, UserID: userID, Name: "tag1"},
			{Model: gorm.Model{ID: 2}, UserID: userID, Name: "tag2"},
		}

		mockTagRepo := new(mocks.MockTagRepository)
		mockTagRepo.On("GetByUser", userID).Return(matchedList, nil).Once()

		service := expense.NewTagService(mockTagRepo)
		list, err := service.GetTags(userID)

		assert.Equal(t, matchedList, list)
		assert.NoError(t, err)
		mockTagRepo.AssertExpectations(t)
	})
}

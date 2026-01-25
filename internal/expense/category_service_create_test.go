package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 1
		dto := expense.CreateCategoryRequest{
			Name: "cat1",
		}
		var newEntity *expense.CategoryEntity

		mockCategoryRepo := new(mocks.MockCategoryRepository)
		mockCategoryRepo.On("ExistsByName", userID, dto.Name).Return(false, nil).Once()
		mockCategoryRepo.On("Create", mock.MatchedBy(func(e *expense.CategoryEntity) bool {
			if e.UserID != userID || e.Name != dto.Name {
				return false
			}
			newEntity = e
			return true
		})).Return(nil).Once()

		service := expense.NewCategoryService(mockCategoryRepo)
		entity, err := service.CreateCategory(userID, dto)

		assert.Equal(t, newEntity, entity)
		assert.NoError(t, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}

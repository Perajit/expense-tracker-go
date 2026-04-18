package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoryService_CreateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 1
		var newEntity *expense.CategoryEntity

		dto := expense.CreateCategoryRequest{
			Name: "cat1",
		}

		db := testutil.SetupDB()

		mockCategoryRepo := new(mocks.MockCategoryRepository)
		mockCategoryRepo.On("ExistsByName", dto.Name, userID).Return(false, nil).Once()
		mockCategoryRepo.On("Create", mock.MatchedBy(func(e *expense.CategoryEntity) bool {
			if e.UserID != userID || e.Name != dto.Name {
				return false
			}
			newEntity = e
			return true
		})).Return(nil).Once()

		mockExpenseRepo := new(mocks.MockExpenseRepository)

		service := expense.NewCategoryService(db, mockCategoryRepo, mockExpenseRepo)
		entity, err := service.CreateCategory(userID, dto)

		assert.Equal(t, newEntity, entity)
		assert.NoError(t, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}

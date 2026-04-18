package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCategoryService_UpdateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11
		newName := "new"

		existingEntity := &expense.CategoryEntity{
			Model:  gorm.Model{ID: id},
			UserID: userID,
			Name:   "cat1",
		}

		dto := expense.UpdateTagRequest{Name: &newName}

		db := testutil.SetupDB()

		mockCategoryRepo := new(mocks.MockCategoryRepository)
		mockCategoryRepo.On("GetByIDAndUser", id, userID).Return(existingEntity, nil).Once()
		mockCategoryRepo.On("ExistsByName", newName, userID).Return(false, nil).Once()
		mockCategoryRepo.On("Update", mock.MatchedBy(func(e *expense.CategoryEntity) bool {
			if e.ID != id || e.UserID != userID {
				return false
			}
			if e.Name != *dto.Name {
				return false
			}
			return true
		})).Return(nil).Once()

		mockExpenseRepo := new(mocks.MockExpenseRepository)

		service := expense.NewCategoryService(db, mockCategoryRepo, mockExpenseRepo)
		err := service.UpdateCategory(id, userID, expense.UpdateCategoryRequest(dto))

		assert.Nil(t, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}

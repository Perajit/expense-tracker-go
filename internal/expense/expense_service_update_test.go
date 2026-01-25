package expense_test

import (
	"slices"
	"testing"
	"time"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/expense/mocks"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUpdateExpense(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		var userID uint = 11
		newDate := time.Now().Add(-time.Minute)
		newAmount := decimal.NewFromInt(100)
		newNote := "new"
		var newCategoryID uint = 5
		tags := []expense.TagEntity{
			{Model: gorm.Model{ID: 6}, UserID: userID, Name: "tag6"},
			{Model: gorm.Model{ID: 7}, UserID: userID, Name: "tag7"},
			{Model: gorm.Model{ID: 8}, UserID: userID, Name: "tag8"},
		}
		dto := expense.UpdateExpenseRequest{
			Date:       &newDate,
			Amount:     &newAmount,
			Note:       &newNote,
			CategoryID: &newCategoryID,
			TagIDs:     &[]uint{tags[0].ID, tags[1].ID, tags[2].ID},
		}
		existingEntity := &expense.ExpenseEntity{
			Model:      gorm.Model{ID: id},
			UserID:     userID,
			Date:       time.Now().Add(-time.Hour).Unix(),
			Amount:     decimal.NewFromInt(100),
			Note:       "expense1",
			CategoryID: 2,
		}

		db := testutil.SetupDB()

		mockExpenseRepo := new(mocks.MockExpenseRepository)
		mockExpenseRepo.On("GetByIDAndUserNoAssociation", id, userID).Return(existingEntity, nil)
		mockExpenseRepo.On("WithTx", mock.Anything).Return(mockExpenseRepo)
		mockExpenseRepo.On("Update", mock.MatchedBy(func(e *expense.ExpenseEntity) bool {
			if e.ID != existingEntity.ID || e.UserID != existingEntity.UserID {
				return false
			}
			if e.Date != dto.Date.Unix() || e.Amount != *dto.Amount || e.Note != *dto.Note || e.CategoryID != *dto.CategoryID {
				return false
			}
			if !slices.Equal(e.Tags, tags) {
				return false
			}
			return true
		})).Return(nil).Once()

		mockCategoryService := new(mocks.MockCategoryService)
		mockCategoryService.On("IsCategoryOwner", newCategoryID, userID).Return(true, nil).Once()

		mockTagService := new(mocks.MockTagService)
		mockTagService.On("GetTagsByIDs", mock.MatchedBy(func(l []uint) bool {
			return slices.Equal(l, *dto.TagIDs)
		}), userID).Return(tags, nil)

		service := expense.NewExpenseService(db, mockExpenseRepo, mockCategoryService, mockTagService)
		err := service.UpdateExpense(id, userID, dto)

		assert.NoError(t, err)
		mockExpenseRepo.AssertExpectations(t)
	})
}

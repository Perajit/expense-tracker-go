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

func TestCreateExpense(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var userID uint = 11
		tags := []expense.TagEntity{
			{Model: gorm.Model{ID: 3}, UserID: userID, Name: "tag3"},
			{Model: gorm.Model{ID: 4}, UserID: userID, Name: "tag4"},
		}
		dto := expense.CreateExpenseRequest{
			Date:       time.Now(),
			Amount:     decimal.NewFromInt(100),
			Note:       "expense1",
			CategoryID: 2,
			TagIDs:     []uint{tags[0].ID, tags[1].ID},
		}
		var createdEntity *expense.ExpenseEntity

		db := testutil.SetupDB()

		mockExpenseRepo := new(mocks.MockExpenseRepository)
		mockExpenseRepo.On("Create", mock.MatchedBy(func(e *expense.ExpenseEntity) bool {
			if e.UserID != userID || e.Amount != dto.Amount || e.Note != dto.Note || e.CategoryID != dto.CategoryID {
				return false
			}
			if !slices.Equal(e.Tags, tags) {
				return false
			}
			createdEntity = e
			return true
		})).Return(nil).Once()

		mockCategoryService := new(mocks.MockCategoryService)
		mockCategoryService.On("IsCategoryOwner", dto.CategoryID, userID).Return(true, nil).Once()

		mockTagService := new(mocks.MockTagService)
		mockTagService.On("GetTagsByIDs", dto.TagIDs, userID).Return(tags, nil).Once()

		service := expense.NewExpenseService(db, mockExpenseRepo, mockCategoryService, mockTagService)
		entity, err := service.CreateExpense(userID, dto)

		assert.Equal(t, createdEntity, entity)
		assert.NoError(t, err)
		mockExpenseRepo.AssertExpectations(t)
	})
}

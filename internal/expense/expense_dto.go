package expense

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateExpenseRequest struct {
	Date       time.Time       `json:"date" validate:"required"`
	Amount     decimal.Decimal `json:"amount" validate:"required"`
	Note       string          `json:"note"`
	CategoryID uint            `json:"categoyId"`
	TagIDs     []uint          `json:"tagIds"`
}

type UpdateExpenseRequest struct {
	Date       *time.Time       `json:"date"`
	Amount     *decimal.Decimal `json:"amount"`
	Note       *string          `json:"note"`
	CategoryID *uint            `json:"categoyId"`
	TagIDs     *[]uint          `json:"tagIds"`
}

type ExpenseResponse struct {
	ID       uint             `json:"id"`
	Date     time.Time        `json:"date"`
	Amount   decimal.Decimal  `json:"amount"`
	Note     string           `json:"note"`
	Category CategoryResponse `json:"categoy"`
	Tags     []TagResponse    `json:"tags"`
}

func (ExpenseResponse) FromEntity(expense ExpenseEntity) ExpenseResponse {
	tagResponses := []TagResponse{}
	for _, tag := range expense.Tags {
		tagResponses = append(tagResponses, TagResponse{}.FromEntity(tag))
	}

	return ExpenseResponse{
		ID:       expense.ID,
		Date:     time.Unix(expense.Date, 0),
		Amount:   expense.Amount,
		Note:     expense.Note,
		Category: CategoryResponse{}.FromEntity(expense.Category),
		Tags:     tagResponses,
	}
}

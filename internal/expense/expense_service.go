package expense

import (
	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"gorm.io/gorm"
)

type ExpenseService interface {
	GetExpenses(authUserID uint) ([]ExpenseEntity, error)
	GetExpenseByID(id uint, authUserID uint) (*ExpenseEntity, error)
	CreateExpense(authUserID uint, dto CreateExpenseRequest) (*ExpenseEntity, error)
	UpdateExpense(id uint, authUserID uint, dto UpdateExpenseRequest) error
	DeleteExpense(id uint, authUserID uint) error
}

type expenseService struct {
	db              *gorm.DB
	expenseRepo     ExpenseRepository
	categoryService CategoryService
	tagService      TagService
}

func NewExpenseService(db *gorm.DB, expenseRepo ExpenseRepository, categoryService CategoryService, tagService TagService) ExpenseService {
	return &expenseService{
		db:              db,
		expenseRepo:     expenseRepo,
		categoryService: categoryService,
		tagService:      tagService,
	}
}

func (s *expenseService) GetExpenses(authUserID uint) ([]ExpenseEntity, error) {
	return s.expenseRepo.GetByUser(authUserID)
}

func (s *expenseService) GetExpenseByID(id uint, authUserID uint) (*ExpenseEntity, error) {
	return s.expenseRepo.GetByIDAndUser(id, authUserID)
}

func (s *expenseService) CreateExpense(authUserID uint, dto CreateExpenseRequest) (*ExpenseEntity, error) {
	isOwner, err := s.categoryService.IsCategoryOwner(dto.CategoryID, authUserID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, apperror.ErrUnauthorized
	}

	tags, err := s.tagService.GetTagsByIDs(dto.TagIDs, authUserID)
	if err != nil {
		return nil, err
	}

	expense := &ExpenseEntity{
		UserID:     authUserID,
		Date:       dto.Date.Unix(),
		Amount:     dto.Amount,
		Note:       dto.Note,
		CategoryID: dto.CategoryID,
		Tags:       tags,
	}
	if err := s.expenseRepo.Create(expense); err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *expenseService) UpdateExpense(id uint, authUserID uint, dto UpdateExpenseRequest) error {
	expense, err := s.expenseRepo.GetByIDAndUserNoAssociation(id, authUserID)
	if err != nil {
		return apperror.ErrNotFound
	}

	if dto.Date != nil {
		expense.Date = dto.Date.Unix()
	}

	if dto.Amount != nil {
		expense.Amount = *dto.Amount
	}

	if dto.Note != nil {
		expense.Note = *dto.Note
	}

	if dto.CategoryID != nil {
		isOwner, err := s.categoryService.IsCategoryOwner(*dto.CategoryID, authUserID)
		if err != nil {
			return err
		}
		if !isOwner {
			return apperror.ErrUnauthorized
		}

		expense.CategoryID = *dto.CategoryID
	}

	if dto.TagIDs != nil {
		tags, err := s.tagService.GetTagsByIDs(*dto.TagIDs, authUserID)
		if err != nil {
			return err
		}

		expense.Tags = tags
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		expenseRepo := s.expenseRepo.WithTx(tx)

		if err := expenseRepo.Update(expense); err != nil {
			return err
		}

		if err := expenseRepo.UpdateTags(expense, expense.Tags); err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (s *expenseService) DeleteExpense(id uint, authUserID uint) error {
	isOwner, err := s.expenseRepo.IsOwner(id, authUserID)
	if err != nil {
		return err
	}
	if !isOwner {
		return apperror.ErrUnauthorized
	}

	return s.expenseRepo.Delete(id)
}

package expense

import (
	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"gorm.io/gorm"
)

type CategoryService interface {
	GetCategoryByID(id uint, authUserID uint) (*CategoryEntity, error)
	GetCategories(authUserID uint) ([]CategoryEntity, error)
	IsCategoryOwner(id uint, authUserID uint) (bool, error)
	CreateCategory(authUserID uint, dto CreateCategoryRequest) (*CategoryEntity, error)
	UpdateCategory(id uint, authUserID uint, dto UpdateCategoryRequest) error
	DeleteCategory(id uint, userId uint) error
}

type categoryService struct {
	db           *gorm.DB
	categoryRepo CategoryRepository
	expenseRepo  ExpenseRepository
}

func NewCategoryService(db *gorm.DB, categoryRepo CategoryRepository, expenseRepo ExpenseRepository) CategoryService {
	return &categoryService{
		db:           db,
		categoryRepo: categoryRepo,
		expenseRepo:  expenseRepo,
	}
}

func (s *categoryService) GetCategoryByID(id uint, authUserID uint) (*CategoryEntity, error) {
	return s.categoryRepo.GetByIDAndUser(id, authUserID)
}

func (s *categoryService) GetCategories(authUserID uint) ([]CategoryEntity, error) {
	return s.categoryRepo.GetAvailableByUser(authUserID)
}

func (s *categoryService) IsCategoryOwner(id uint, authUserID uint) (bool, error) {
	return s.categoryRepo.IsOwner(id, authUserID)
}

func (s *categoryService) CreateCategory(authUserID uint, dto CreateCategoryRequest) (*CategoryEntity, error) {
	duplicated, err := s.categoryRepo.ExistsByName(dto.Name, authUserID)
	if err != nil {
		return nil, err
	}
	if duplicated {
		return nil, apperror.ErrRecordDuplication
	}

	category := &CategoryEntity{
		UserID: authUserID,
		Name:   dto.Name,
	}
	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) UpdateCategory(id uint, authUserID uint, dto UpdateCategoryRequest) error {
	category, err := s.categoryRepo.GetByIDAndUser(id, authUserID)
	if err != nil {
		return err
	}

	if dto.Name != nil {
		duplicated, err := s.categoryRepo.ExistsByName(*dto.Name, authUserID)
		if err != nil {
			return err
		}
		if duplicated {
			return apperror.ErrRecordDuplication
		}

		category.Name = *dto.Name
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return err
	}

	return nil
}

func (s *categoryService) DeleteCategory(id uint, authUserID uint) error {
	isOwner, err := s.categoryRepo.IsOwner(id, authUserID)
	if err != nil {
		return err
	}
	if !isOwner {
		return apperror.ErrUnauthorized
	}

	defaultID, err := s.categoryRepo.GetDefaultCategoryID()
	if err != nil {
		return err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		categoryRepo := s.categoryRepo.WithTx(tx)
		expenseRepo := s.expenseRepo.WithTx(tx)

		if err := expenseRepo.ReplaceCategoryID(id, defaultID); err != nil {
			return err
		}

		if err := categoryRepo.Delete(id); err != nil {
			return err
		}

		return nil
	})

	return nil
}

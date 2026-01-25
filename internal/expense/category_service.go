package expense

import "github.com/Perajit/expense-tracker-go/internal/apperror"

type CategoryService interface {
	GetCategoryByID(id uint, authUserID *uint) (*CategoryEntity, error)
	GetCategories(authUserID uint) ([]CategoryEntity, error)
	IsCategoryOwner(id uint, authUserID uint) (bool, error)
	CreateCategory(authUserID uint, dto CreateCategoryRequest) (*CategoryEntity, error)
	UpdateCategory(id uint, authUserID uint, dto UpdateCategoryRequest) error
	DeleteCategory(id uint, userId uint) error
}

type categoryService struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(categoryRepo CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) GetCategoryByID(id uint, authUserID *uint) (*CategoryEntity, error) {
	return s.categoryRepo.GetByIDAndUser(id, authUserID)
}

func (s *categoryService) GetCategories(authUserID uint) ([]CategoryEntity, error) {
	return s.categoryRepo.GetByUser(authUserID)
}

func (s *categoryService) IsCategoryOwner(id uint, authUserID uint) (bool, error) {
	return s.categoryRepo.IsOwner(id, authUserID)
}

func (s *categoryService) CreateCategory(authUserID uint, dto CreateCategoryRequest) (*CategoryEntity, error) {
	duplicated, err := s.categoryRepo.ExistsByName(authUserID, dto.Name)
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
	category, err := s.categoryRepo.GetByIDAndUser(id, &authUserID)
	if err != nil {
		return err
	}

	if dto.Name != nil {
		duplicated, err := s.categoryRepo.ExistsByName(authUserID, *dto.Name)
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

	return s.categoryRepo.Delete(id)
}

package expense

import (
	"gorm.io/gorm"
)

type CategoryRepository interface {
	WithTx(tx *gorm.DB) CategoryRepository
	GetByIDAndUser(id uint, userID uint) (*CategoryEntity, error)
	GetAvailableByUser(userID uint) ([]CategoryEntity, error)
	GetDefaultCategoryID() (uint, error)
	IsOwner(id uint, userID uint) (bool, error)
	ExistsByName(name string, userID uint) (bool, error)
	Create(category *CategoryEntity) error
	Update(category *CategoryEntity) error
	Delete(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) WithTx(tx *gorm.DB) CategoryRepository {
	if tx == nil {
		return r
	}

	return &categoryRepository{db: tx}
}

func (r *categoryRepository) GetByIDAndUser(id uint, userID uint) (*CategoryEntity, error) {
	var category CategoryEntity
	err := r.db.Model(&CategoryEntity{}).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		First(&category).
		Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) GetAvailableByUser(userID uint) ([]CategoryEntity, error) {
	var categories []CategoryEntity
	err := r.db.Model(&CategoryEntity{}).
		Where("user_id = ?", userID).
		Or("user_id = ?", 0).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) GetDefaultCategoryID() (uint, error) {
	var category CategoryEntity
	err := r.db.Model(&CategoryEntity{}).
		Where("user_id = 0").
		Where("is_default = ?", true).
		First(&category).
		Error
	if err != nil {
		return 0, err
	}

	return category.ID, nil
}

func (r *categoryRepository) IsOwner(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&CategoryEntity{}).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Count(&count).
		Error

	return count > 0, err
}

func (r *categoryRepository) ExistsByName(name string, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&CategoryEntity{}).
		Where("name = ?", name).
		Where(r.db.Where("user_id = ?", userID).Or("user_id = ?", 0)).
		Count(&count).Error

	return count > 0, err
}

func (r *categoryRepository) Create(category *CategoryEntity) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *CategoryEntity) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&CategoryEntity{}, id).Error
}

package expense

import "gorm.io/gorm"

type CategoryRepository interface {
	GetByIDAndUser(id uint, userID *uint) (*CategoryEntity, error)
	GetByUser(userID uint) ([]CategoryEntity, error)
	IsOwner(id uint, userID uint) (bool, error)
	ExistsByName(userID uint, name string) (bool, error)
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

func (r *categoryRepository) GetByIDAndUser(id uint, userID *uint) (*CategoryEntity, error) {
	var category CategoryEntity
	if err := r.db.Where("id = ?", id).Where("user_id = ?", userID).First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) GetByUser(userID uint) ([]CategoryEntity, error) {
	var categories []CategoryEntity
	if err := r.db.Where("user_id = ?", userID).Or("is_default = ?", false).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) IsOwner(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Where("id = ?", id).Where("user_id = ?", userID).Count(&count).Error

	return count > 0, err
}

func (r *categoryRepository) ExistsByName(userID uint, name string) (bool, error) {
	var count int64
	err := r.db.Where("name = ?", name).
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
	return r.db.Delete(id).Error
}

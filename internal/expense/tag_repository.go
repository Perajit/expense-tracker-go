package expense

import (
	"gorm.io/gorm"
)

type TagRepository interface {
	GetByIDAndUser(id uint, userID uint) (*TagEntity, error)
	GetByIDsAndUser(ids []uint, userID uint) ([]TagEntity, error)
	GetByUser(userID uint) ([]TagEntity, error)
	IsOwner(id uint, userID uint) (bool, error)
	Create(tag *TagEntity) error
	Update(tag *TagEntity) error
	Delete(id uint) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) GetByIDAndUser(id uint, userID uint) (*TagEntity, error) {
	var tag TagEntity
	if err := r.db.Where("id = ?", id).Where("user_id = ?", userID).First(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *tagRepository) GetByIDsAndUser(ids []uint, userID uint) ([]TagEntity, error) {
	var tags []TagEntity
	if err := r.db.Where("id IN ?", ids).Where("user_id = ?", userID).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagRepository) GetByUser(userID uint) ([]TagEntity, error) {
	var tags []TagEntity
	if err := r.db.Where("user_id = ?", userID).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagRepository) IsOwner(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&TagEntity{}).Where("id = ?", id).Where("user_id = ?", userID).Count(&count).Error

	return count > 0, err
}

func (r *tagRepository) Create(tag *TagEntity) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) Update(tag *TagEntity) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	return r.db.Delete(id).Error
}

package user

import "gorm.io/gorm"

type UserRepository interface {
	GetByID(id uint) (*UserEntity, error)
	GetByUsername(email string) (*UserEntity, error)
	ExistsByUsername(email string) (bool, error)
	Create(user *UserEntity) error
	Update(user *UserEntity) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(id uint) (*UserEntity, error) {
	var user UserEntity
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*UserEntity, error) {
	var user UserEntity
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&UserEntity{}).Where("username = ?", username).Count(&count).Error

	return count > 0, err
}

func (r *userRepository) Create(user *UserEntity) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *UserEntity) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&UserEntity{}, id).Error
}

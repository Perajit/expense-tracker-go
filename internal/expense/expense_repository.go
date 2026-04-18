package expense

import "gorm.io/gorm"

type ExpenseRepository interface {
	WithTx(tx *gorm.DB) ExpenseRepository
	GetByIDAndUser(id uint, userID uint) (*ExpenseEntity, error)
	GetByIDAndUserNoAssociation(id uint, userID uint) (*ExpenseEntity, error)
	GetByUser(userID uint) ([]ExpenseEntity, error)
	IsOwner(id uint, userID uint) (bool, error)
	Create(expense *ExpenseEntity) error
	Update(expense *ExpenseEntity) error
	Delete(id uint) error
	UpdateTags(expense *ExpenseEntity, tags []TagEntity) error
	ReplaceCategoryID(categoryID uint, newCategoryID uint) error
	UnlinkTag(tagID uint) error
}

type expenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) WithTx(tx *gorm.DB) ExpenseRepository {
	if tx == nil {
		return r
	}

	return &expenseRepository{db: tx}
}

func (r *expenseRepository) GetByIDAndUser(id uint, userID uint) (*ExpenseEntity, error) {
	var expense ExpenseEntity
	err := r.db.Model(&ExpenseEntity{}).
		Preload("Category").
		Preload("Tags").
		Where("id = ?", id).
		Where("user_id = ?", userID).
		First(&expense).
		Error
	if err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *expenseRepository) GetByIDAndUserNoAssociation(id uint, userID uint) (*ExpenseEntity, error) {
	var expense ExpenseEntity
	err := r.db.Model(&ExpenseEntity{}).
		First(&expense, id).
		Error
	if err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *expenseRepository) GetByUser(userID uint) ([]ExpenseEntity, error) {
	var expenses []ExpenseEntity
	err := r.db.Preload("Category").
		Preload("Tags").
		Where("user_id = ?", userID).
		Find(&expenses).
		Error
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *expenseRepository) IsOwner(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&ExpenseEntity{}).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Count(&count).
		Error

	return count > 0, err
}

func (r *expenseRepository) Create(expense *ExpenseEntity) error {
	return r.db.Create(expense).Error
}

func (r *expenseRepository) Update(expense *ExpenseEntity) error {
	return r.db.Save(expense).Error
}

func (r *expenseRepository) Delete(id uint) error {
	return r.db.Delete(&ExpenseEntity{}, id).Error
}

func (r *expenseRepository) UpdateTags(expense *ExpenseEntity, tags []TagEntity) error {
	return r.db.Model(expense).Association("Tags").Replace(tags)
}

func (r *expenseRepository) ReplaceCategoryID(categoryID uint, newCategoryID uint) error {
	return r.db.Model(&ExpenseEntity{}).
		Where("category_id = ?", categoryID).
		Update("category_id", newCategoryID).
		Error
}

func (r *expenseRepository) UnlinkTag(tagID uint) error {
	return r.db.Table("expense_tag_links").
		Where("tag_entity_id = ?", tagID).
		Delete(map[string]interface{}{}).
		Error
}

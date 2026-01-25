package expense

import "gorm.io/gorm"

type ExpenseRepository interface {
	WithTx(tx *gorm.DB) ExpenseRepository
	GetByUser(userID uint) ([]ExpenseEntity, error)
	GetByIDAndUser(id uint, userID uint) (*ExpenseEntity, error)
	GetByIDAndUserNoAssociation(id uint, userID uint) (*ExpenseEntity, error)
	IsOwner(id uint, userID uint) (bool, error)
	Create(expense *ExpenseEntity) error
	Update(expense *ExpenseEntity) error
	Delete(id uint) error
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

func (r *expenseRepository) GetByUser(userID uint) ([]ExpenseEntity, error) {
	var expenses []ExpenseEntity
	if err := r.db.Preload("Category").
		Preload("Tags").
		Where("user_id = ?", userID).
		Find(&expenses).
		Error; err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *expenseRepository) GetByIDAndUser(id uint, userID uint) (*ExpenseEntity, error) {
	var expense ExpenseEntity
	if err := r.db.Preload("Category").
		Preload("Tags").
		Where("id = ?", id).
		Where("user_id = ?", userID).
		First(&expense).
		Error; err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *expenseRepository) GetByIDAndUserNoAssociation(id uint, userID uint) (*ExpenseEntity, error) {
	var expense ExpenseEntity
	if err := r.db.First(&expense, id).Error; err != nil {
		return nil, err
	}

	return &expense, nil
}

func (r *expenseRepository) IsOwner(id uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&ExpenseEntity{}).Where("id = ?", id).Where("user_id = ?", userID).Count(&count).Error

	return count > 0, err
}

func (r *expenseRepository) Create(expense *ExpenseEntity) error {
	return r.db.Create(expense).Error
}

func (r *expenseRepository) Update(expense *ExpenseEntity) error {
	if err := r.db.Save(expense).Error; err != nil {
		return err
	}

	tags := expense.Tags
	if err := r.db.Model(expense).Association("Tags").Replace(tags); err != nil {
		return err
	}

	return nil
}

func (r *expenseRepository) Delete(id uint) error {
	return r.db.Delete(id).Error
}

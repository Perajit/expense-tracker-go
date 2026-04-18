package expense

import (
	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"gorm.io/gorm"
)

type TagService interface {
	GetTags(authUserID uint) ([]TagEntity, error)
	GetTagByID(id uint, authUserID uint) (*TagEntity, error)
	GetTagsByIDs(ids []uint, authUserID uint) ([]TagEntity, error)
	CreateTag(authUserID uint, dto CreateTagRequest) (*TagEntity, error)
	UpdateTag(id uint, authUserID uint, dto UpdateTagRequest) error
	DeleteTag(id uint, authUserID uint) error
}

type tagService struct {
	db          *gorm.DB
	tagRepo     TagRepository
	expenseRepo ExpenseRepository
}

func (s *tagService) GetTags(authUserID uint) ([]TagEntity, error) {
	return s.tagRepo.GetByUser(authUserID)
}

func NewTagService(db *gorm.DB, tagRepo TagRepository, expenseRepo ExpenseRepository) TagService {
	return &tagService{
		db:          db,
		tagRepo:     tagRepo,
		expenseRepo: expenseRepo,
	}
}

func (s *tagService) GetTagByID(id uint, authUserID uint) (*TagEntity, error) {
	return s.tagRepo.GetByIDAndUser(id, authUserID)
}

func (s *tagService) GetTagsByIDs(ids []uint, authUserID uint) ([]TagEntity, error) {
	return s.tagRepo.GetByIDsAndUser(ids, authUserID)
}

func (s *tagService) CreateTag(authUserID uint, dto CreateTagRequest) (*TagEntity, error) {
	tag := &TagEntity{
		UserID: authUserID,
		Name:   dto.Name,
	}
	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *tagService) UpdateTag(id uint, authUserID uint, dto UpdateTagRequest) error {
	tag, err := s.tagRepo.GetByIDAndUser(id, authUserID)
	if err != nil {
		return err
	}

	if dto.Name != nil {
		tag.Name = *dto.Name
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return err
	}

	return nil
}

func (s *tagService) DeleteTag(id uint, authUserID uint) error {
	isOwner, err := s.tagRepo.IsOwner(id, authUserID)
	if err != nil {
		return err
	}
	if !isOwner {
		return apperror.ErrUnauthorized
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		tagRepo := s.tagRepo.WithTx(tx)
		expenseRepo := s.expenseRepo.WithTx(tx)

		if err := expenseRepo.UnlinkTag(id); err != nil {
			return err
		}

		if err := tagRepo.Delete(id); err != nil {
			return err
		}

		return nil
	})

	return nil
}

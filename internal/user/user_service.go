package user

import (
	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/util"
)

type UserService interface {
	GetUserByID(id uint, authUserID uint) (*UserEntity, error)
	GetUserByUsername(email string) (*UserEntity, error)
	CreateUser(dto CreateUserRequest) (*UserEntity, error)
	UpdateUser(id uint, authUserID uint, dto UpdateUserRequest) error
	DeleteUser(id uint, authUserID uint) error
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id uint, authUserID uint) (*UserEntity, error) {
	if id != authUserID {
		return nil, apperror.ErrUnauthorized
	}

	return s.userRepo.GetByID(id)
}

func (s *userService) GetUserByUsername(username string) (*UserEntity, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *userService) CreateUser(dto CreateUserRequest) (*UserEntity, error) {
	exists, err := s.userRepo.ExistsByUsername(dto.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.ErrUserDuplication
	}

	hashedPassword, err := util.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user := &UserEntity{
		Username: dto.Username,
		Password: hashedPassword,
		Email:    dto.Email,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(id uint, authUserID uint, dto UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if user.ID != authUserID {
		return apperror.ErrUnauthorized
	}

	if dto.Password != nil {
		hashedPassword, err := util.HashPassword(*dto.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	if dto.Email != nil {
		user.Email = *dto.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteUser(id uint, authUserID uint) error {
	if id != authUserID {
		return apperror.ErrUnauthorized
	}

	return s.userRepo.Delete(id)
}

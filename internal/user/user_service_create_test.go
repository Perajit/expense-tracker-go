package user_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/Perajit/expense-tracker-go/internal/user/mocks"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dto := user.CreateUserRequest{
			Username: "user1",
			Password: "pwd123",
			Email:    "test@example.com",
		}
		var newEntity *user.UserEntity

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("ExistsByUsername", dto.Username).Return(false, nil).Once()
		mockUserRepo.On("Create", mock.MatchedBy(func(e *user.UserEntity) bool {
			if e.Username != dto.Username || e.Email != dto.Email {
				return false
			}
			if err := util.VerifyPassword(e.Password, dto.Password); err != nil {
				return false
			}
			newEntity = e
			return true
		})).Return(nil).Once()

		service := user.NewUserService(mockUserRepo)
		entity, err := service.CreateUser(dto)

		assert.Equal(t, newEntity, entity)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error_duplicated_email", func(t *testing.T) {
		dto := user.CreateUserRequest{
			Username: "user1",
			Password: "pwd123",
			Email:    "duplicated@example.com",
		}

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("ExistsByUsername", dto.Username).Return(true, nil).Once()

		service := user.NewUserService(mockUserRepo)
		entity, err := service.CreateUser(dto)

		assert.Nil(t, entity)
		assert.Equal(t, apperror.ErrUserDuplication, err)
		mockUserRepo.AssertExpectations(t)
		mockUserRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("error_create", func(t *testing.T) {
		dto := user.CreateUserRequest{
			Username: "user1",
			Password: "pwd123",
			Email:    "test@example.com",
		}

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("ExistsByUsername", dto.Username).Return(false, nil).Once()
		mockUserRepo.On("Create", mock.Anything).Return(apperror.ErrDefault).Once()

		service := user.NewUserService(mockUserRepo)
		entity, err := service.CreateUser(dto)

		assert.Nil(t, entity)
		assert.Equal(t, apperror.ErrDefault, err)
		mockUserRepo.AssertExpectations(t)
	})
}

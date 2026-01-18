package user_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/Perajit/expense-tracker-go/internal/user/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gorm.io/gorm"
)

func TestGetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1
		matchedEntity := &user.UserEntity{
			Model:    gorm.Model{ID: id},
			Username: "user1",
			Password: "hash123",
			Email:    "test@example.com",
		}

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("GetByID", id).Return(matchedEntity, nil).Once()

		service := user.NewUserService(mockUserRepo)
		entity, err := service.GetUserByID(id, id)

		assert.Equal(t, matchedEntity, entity)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error_permission", func(t *testing.T) {
		var id uint = 1

		mockUserRepo := new(mocks.MockUserRepository)

		service := user.NewUserService(mockUserRepo)
		entity, err := service.GetUserByID(id, 11)

		assert.Nil(t, entity)
		assert.Equal(t, apperror.ErrUnauthorized, err)
		mockUserRepo.AssertNotCalled(t, "GetByID", mock.Anything)
	})

	t.Run("error_get", func(t *testing.T) {
		var id uint = 1

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("GetByID", id).Return(nil, apperror.ErrNotFound).Once()

		service := user.NewUserService(mockUserRepo)
		entity, err := service.GetUserByID(id, id)

		assert.Nil(t, entity)
		assert.Equal(t, apperror.ErrNotFound, err)
	})
}

func TestGetUserByUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		matchedEntity := &user.UserEntity{
			Model:    gorm.Model{ID: 1},
			Username: "user1",
			Password: "hash123",
			Email:    "test@example.com",
		}

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("GetByUsername", matchedEntity.Username).Return(matchedEntity, nil).Once()

		service := user.NewUserService(mockUserRepo)
		entity, err := service.GetUserByUsername(matchedEntity.Username)

		assert.Equal(t, matchedEntity, entity)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		username := "notfound"

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("GetByUsername", username).Return(nil, apperror.ErrNotFound).Once()

		service := user.NewUserService(mockUserRepo)
		entity, err := service.GetUserByUsername(username)

		assert.Nil(t, entity)
		assert.Equal(t, apperror.ErrNotFound, err)
	})
}

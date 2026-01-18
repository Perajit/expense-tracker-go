package user_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/user"
	userMocks "github.com/Perajit/expense-tracker-go/internal/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1

		mockUserRepo := new(userMocks.MockUserRepository)
		mockUserRepo.On("Delete", id).Return(nil).Once()

		service := user.NewUserService(mockUserRepo)
		err := service.DeleteUser(id, id)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error_permission", func(t *testing.T) {
		var id uint = 1

		mockUserRepo := new(userMocks.MockUserRepository)

		service := user.NewUserService(mockUserRepo)
		err := service.DeleteUser(id, 11)

		assert.Equal(t, apperror.ErrUnauthorized, err)
		mockUserRepo.AssertNotCalled(t, "Delete", mock.Anything)
	})

	t.Run("error_delete", func(t *testing.T) {
		var id uint = 1

		mockUserRepo := new(userMocks.MockUserRepository)
		mockUserRepo.On("Delete", id).Return(apperror.ErrDefault).Once()

		service := user.NewUserService(mockUserRepo)
		err := service.DeleteUser(id, id)

		assert.Equal(t, apperror.ErrDefault, err)
		mockUserRepo.AssertExpectations(t)
	})
}

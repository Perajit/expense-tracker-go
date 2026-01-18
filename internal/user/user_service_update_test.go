package user_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/Perajit/expense-tracker-go/internal/user/mocks"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUpdateUser(t *testing.T) {
	var id uint = 1
	newPassword := "pwd456"
	newEmail := "new@example.com"

	tests_success := []struct {
		name       string
		dto        user.UpdateUserRequest
		mockUpdate func(mockUserRepo *mocks.MockUserRepository, matchedEntity *user.UserEntity, dto user.UpdateUserRequest)
	}{
		{
			"success_update_password",
			user.UpdateUserRequest{Password: &newPassword},
			func(mockUserRepo *mocks.MockUserRepository, matchedEntity *user.UserEntity, dto user.UpdateUserRequest) {
				mockUserRepo.On("Update", mock.MatchedBy(func(e *user.UserEntity) bool {
					if e.ID != id || e.Username != matchedEntity.Username || e.Email != matchedEntity.Email {
						return false
					}
					if err := util.VerifyPassword(e.Password, *dto.Password); err != nil {
						return false
					}
					return true
				})).Return(nil).Once()
			},
		},
		{
			"success_update_email",
			user.UpdateUserRequest{Email: &newEmail},
			func(mockUserRepo *mocks.MockUserRepository, matchedEntity *user.UserEntity, dto user.UpdateUserRequest) {
				mockUserRepo.On("Update", mock.MatchedBy(func(e *user.UserEntity) bool {
					if e.ID != id || e.Username != matchedEntity.Username || e.Password != matchedEntity.Password {
						return false
					}
					if e.Email != *dto.Email {
						return false
					}
					return true
				})).Return(nil).Once()
			},
		},
		{
			"success_update_all",
			user.UpdateUserRequest{Password: &newPassword, Email: &newEmail},
			func(mockUserRepo *mocks.MockUserRepository, matchedEntity *user.UserEntity, dto user.UpdateUserRequest) {
				mockUserRepo.On("Update", mock.MatchedBy(func(e *user.UserEntity) bool {
					if e.ID != id || e.Username != matchedEntity.Username {
						return false
					}
					if e.Email != *dto.Email {
						return false
					}
					if err := util.VerifyPassword(e.Password, *dto.Password); err != nil {
						return false
					}
					return true
				})).Return(nil).Once()
			},
		},
	}

	for _, tc := range tests_success {
		t.Run(tc.name, func(t *testing.T) {
			existingEntity := &user.UserEntity{
				Model:    gorm.Model{ID: id},
				Username: "user1",
				Password: "pwd123",
				Email:    "test@example.com",
			}

			mockUserRepo := new(mocks.MockUserRepository)
			mockUserRepo.On("GetByID", id).Return(existingEntity, nil).Once()
			tc.mockUpdate(mockUserRepo, existingEntity, tc.dto)

			service := user.NewUserService(mockUserRepo)
			err := service.UpdateUser(id, id, tc.dto)

			assert.NoError(t, err)
			mockUserRepo.AssertExpectations(t)
		})
	}

	t.Run("error_notfound", func(t *testing.T) {
		dto := user.UpdateUserRequest{
			Password: &newPassword,
			Email:    &newEmail,
		}

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("GetByID", id).Return(nil, apperror.ErrNotFound).Once()

		service := user.NewUserService(mockUserRepo)
		err := service.UpdateUser(id, id, dto)

		assert.Equal(t, apperror.ErrNotFound, err)
		mockUserRepo.AssertExpectations(t)
		mockUserRepo.AssertNotCalled(t, "Update", mock.Anything)
	})

	t.Run("error_permission", func(t *testing.T) {
		dto := user.UpdateUserRequest{
			Password: &newPassword,
			Email:    &newEmail,
		}
		existingEntity := &user.UserEntity{
			Model:    gorm.Model{ID: id},
			Password: "pwd123",
			Email:    "test@example.com",
		}

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("GetByID", id).Return(existingEntity, nil).Once()

		service := user.NewUserService(mockUserRepo)
		err := service.UpdateUser(id, 2, dto)

		assert.Equal(t, apperror.ErrUnauthorized, err)
		mockUserRepo.AssertExpectations(t)
		mockUserRepo.AssertNotCalled(t, "Update", mock.Anything)
	})

	t.Run("error_update", func(t *testing.T) {
		dto := user.UpdateUserRequest{
			Password: &newPassword,
			Email:    &newEmail,
		}
		existingEntity := &user.UserEntity{
			Model:    gorm.Model{ID: id},
			Password: "pwd123",
			Email:    "test@example.com",
		}

		mockUserRepo := new(mocks.MockUserRepository)
		mockUserRepo.On("GetByID", id).Return(existingEntity, nil).Once()
		mockUserRepo.On("Update", mock.Anything).Return(apperror.ErrDefault).Once()

		service := user.NewUserService(mockUserRepo)
		err := service.UpdateUser(id, id, dto)

		assert.Equal(t, apperror.ErrDefault, err)
		mockUserRepo.AssertExpectations(t)
	})
}

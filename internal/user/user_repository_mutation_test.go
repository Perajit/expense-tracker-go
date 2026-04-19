package user_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		createdEntity := user.UserEntity{
			Username: "user1",
			Password: "pwd1",
			Email:    "user1@example.com",
		}

		db := testutil.SetupDB()

		repo := user.NewUserRepository(db)

		err := repo.Create(&createdEntity)

		var entity user.UserEntity
		db.Model(&user.UserEntity{}).First(&entity)
		assert.Equal(t, createdEntity.Username, entity.Username)
		assert.Equal(t, createdEntity.Password, entity.Password)
		assert.Equal(t, createdEntity.Email, entity.Email)
		assert.NoError(t, err)
	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1

		insertedUsers := []user.UserEntity{
			{
				Username: "user1",
				Password: "pwd1",
				Email:    "user1@example.com",
			},
			{
				Username: "user2",
				Password: "pwd2",
				Email:    "user2@example.com",
			},
		}
		updatedEntity := user.UserEntity{
			Model:    gorm.Model{ID: id},
			Username: "new",
			Password: "pwd-new",
			Email:    "new@example.com",
		}

		db := testutil.SetupDB()
		db.Create(&insertedUsers)

		repo := user.NewUserRepository(db)

		err := repo.Update(&updatedEntity)

		var entity user.UserEntity
		db.Model(&user.UserEntity{}).Where("id = ?", id).First(&entity)
		assert.Equal(t, updatedEntity.Username, entity.Username)
		assert.Equal(t, updatedEntity.Password, entity.Password)
		assert.Equal(t, updatedEntity.Email, entity.Email)
		assert.NoError(t, err)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var id uint = 1

		insertedUsers := []user.UserEntity{
			{
				Username: "user1",
				Password: "pwd1",
				Email:    "user1@example.com",
			},
			{
				Username: "user2",
				Password: "pwd2",
				Email:    "user2@example.com",
			},
		}

		db := testutil.SetupDB()
		db.Create(&insertedUsers)

		repo := user.NewUserRepository(db)

		err := repo.Delete(id)

		var count int64
		db.Model(&user.UserEntity{}).Find("id = ?", id).Count(&count)
		assert.Equal(t, int64(0), count)
		assert.NoError(t, err)
	})
}

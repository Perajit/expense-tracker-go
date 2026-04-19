package user_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetByID(t *testing.T) {
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

	t.Run("exist", func(t *testing.T) {
		repo := user.NewUserRepository(db)

		entity, err := repo.GetByID(1)

		assert.Equal(t, uint(1), entity.ID)
		assert.Equal(t, insertedUsers[0].Username, entity.Username)
		assert.Equal(t, insertedUsers[0].Password, entity.Password)
		assert.Equal(t, insertedUsers[0].Email, entity.Email)
		assert.NoError(t, err)
	})

	t.Run("no exist", func(t *testing.T) {
		repo := user.NewUserRepository(db)

		entity, err := repo.GetByID(3)

		assert.Nil(t, entity)
		assert.Error(t, err)
	})
}

func TestUserRepository_GetByUsername(t *testing.T) {
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

	t.Run("exist", func(t *testing.T) {
		repo := user.NewUserRepository(db)

		entity, err := repo.GetByUsername("user1")

		assert.Equal(t, uint(1), entity.ID)
		assert.Equal(t, insertedUsers[0].Username, entity.Username)
		assert.Equal(t, insertedUsers[0].Password, entity.Password)
		assert.Equal(t, insertedUsers[0].Email, entity.Email)
		assert.NoError(t, err)
	})

	t.Run("not exist", func(t *testing.T) {
		repo := user.NewUserRepository(db)

		entity, err := repo.GetByUsername("user3")

		assert.Nil(t, entity)
		assert.Error(t, err)
	})
}

func TestUserRepository_ExistsByUsername(t *testing.T) {
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

	t.Run("exist", func(t *testing.T) {
		repo := user.NewUserRepository(db)

		result, err := repo.ExistsByUsername("user1")

		assert.True(t, result)
		assert.NoError(t, err)
	})

	t.Run("not exist", func(t *testing.T) {
		repo := user.NewUserRepository(db)

		result, err := repo.ExistsByUsername("user3")

		assert.False(t, result)
		assert.NoError(t, err)
	})
}

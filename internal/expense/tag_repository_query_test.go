package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTagRepository_GetByIDAndUser(t *testing.T) {
	var id uint = 1
	var userID uint = 11

	insertedTags := []expense.TagEntity{
		{UserID: userID, Name: "tag1"},
		{UserID: userID, Name: "tag2"},
	}

	db := testutil.SetupDB()
	db.Create(&insertedTags)

	t.Run("exist", func(t *testing.T) {
		repo := expense.NewTagRepository(db)

		entity, err := repo.GetByIDAndUser(id, userID)

		assert.Equal(t, id, entity.ID)
		assert.Equal(t, insertedTags[0].Name, entity.Name)
		assert.Equal(t, userID, entity.UserID)
		assert.NoError(t, err)
	})

	t.Run("not exist", func(t *testing.T) {
		repo := expense.NewTagRepository(db)

		entity, err := repo.GetByIDAndUser(id, userID+uint(1))

		assert.Nil(t, entity)
		assert.Error(t, err)
	})
}

func TestTagRepository_GetByIDsAndUser(t *testing.T) {
	var userID uint = 11

	insertedTags := []expense.TagEntity{
		{UserID: userID, Name: "tag1"},
		{UserID: userID, Name: "tag2"},
		{UserID: userID + uint(1), Name: "tag3"},
	}

	db := testutil.SetupDB()
	db.Create(&insertedTags)

	t.Run("all matched", func(t *testing.T) {
		repo := expense.NewTagRepository(db)

		list, err := repo.GetByIDsAndUser([]uint{1, 2}, userID)

		assert.Equal(t, 2, len(list))

		// 1st item
		assert.Equal(t, uint(1), list[0].ID)
		assert.Equal(t, insertedTags[0].Name, list[0].Name)
		assert.Equal(t, userID, list[0].UserID)

		// 2nd item
		assert.Equal(t, uint(2), list[1].ID)
		assert.Equal(t, insertedTags[1].Name, list[1].Name)
		assert.Equal(t, userID, list[1].UserID)

		assert.NoError(t, err)
	})

	t.Run("partially matched", func(t *testing.T) {
		repo := expense.NewTagRepository(db)

		list, err := repo.GetByIDsAndUser([]uint{1, 2}, userID+uint(1))

		assert.Equal(t, 0, len(list))
		assert.NoError(t, err)
	})

	t.Run("unmatched", func(t *testing.T) {
		repo := expense.NewTagRepository(db)

		list, err := repo.GetByIDsAndUser([]uint{1, 2, 3}, userID)

		assert.Equal(t, 2, len(list))

		// 1st item
		assert.Equal(t, uint(1), list[0].ID)
		assert.Equal(t, insertedTags[0].Name, list[0].Name)
		assert.Equal(t, userID, list[0].UserID)

		// 2nd item
		assert.Equal(t, uint(2), list[1].ID)
		assert.Equal(t, insertedTags[1].Name, list[1].Name)
		assert.Equal(t, userID, list[1].UserID)

		assert.NoError(t, err)
	})
}

func TestTagRepository_GetByUser(t *testing.T) {
	var userID uint = 11

	insertedTags := []expense.TagEntity{
		{UserID: userID, Name: "tag1"},
		{UserID: userID, Name: "tag2"},
		{UserID: userID + uint(1), Name: "tag3"},
	}

	db := testutil.SetupDB()
	db.Create(&insertedTags)

	t.Run("exist", func(t *testing.T) {
		repo := expense.NewTagRepository(db)

		list, err := repo.GetByUser(userID)

		assert.Equal(t, 2, len(list))

		// 1st item
		assert.Equal(t, uint(1), list[0].ID)
		assert.Equal(t, insertedTags[0].Name, list[0].Name)
		assert.Equal(t, userID, list[0].UserID)

		// 2nd item
		assert.Equal(t, uint(2), list[1].ID)
		assert.Equal(t, insertedTags[1].Name, list[1].Name)
		assert.Equal(t, userID, list[1].UserID)

		assert.NoError(t, err)
	})

	t.Run("user-category not exist", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		list, err := repo.GetAvailableByUser(userID + uint(2))

		assert.Equal(t, 0, len(list))
		assert.NoError(t, err)
	})
}

func TestTagRepository_IsOwner(t *testing.T) {
	var id uint = 1
	var userID uint = 11

	insertedTags := []expense.TagEntity{
		{UserID: userID, Name: "tag1"},
	}

	db := testutil.SetupDB()
	db.Create(&insertedTags)

	t.Run("owner", func(t *testing.T) {
		repo := expense.NewTagRepository(db)

		// should be true
		result, err := repo.IsOwner(id, userID)

		assert.True(t, result)
		assert.NoError(t, err)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := expense.NewTagRepository(db)

		// should be true
		result, err := repo.IsOwner(id, userID+uint(1))

		assert.False(t, result)
		assert.NoError(t, err)
	})
}

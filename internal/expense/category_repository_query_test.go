package expense_test

import (
	"testing"

	"github.com/Perajit/expense-tracker-go/internal/expense"
	"github.com/Perajit/expense-tracker-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCategoryRepository_GetByIDAndUser(t *testing.T) {
	var id uint = 1
	var userID uint = 11

	insertedCategories := []expense.CategoryEntity{
		{UserID: userID, Name: "cat1"},
		{UserID: userID, Name: "cat2"},
	}

	db := testutil.SetupDB()
	db.Create(&insertedCategories)

	t.Run("exist", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		entity, err := repo.GetByIDAndUser(id, userID)

		assert.Equal(t, id, entity.ID)
		assert.Equal(t, insertedCategories[0].Name, entity.Name)
		assert.Equal(t, userID, entity.UserID)
		assert.NoError(t, err)
	})

	t.Run("not exist", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		entity, err := repo.GetByIDAndUser(2, userID+uint(1))

		assert.Nil(t, entity)
		assert.Error(t, err)
	})
}

func TestCategoryRepository_GetAvailableByUser(t *testing.T) {
	var userID uint = 11

	insertedCategories := []expense.CategoryEntity{
		{UserID: 0, Name: "cat1"},
		{UserID: userID, Name: "cat2"},
		{UserID: userID, Name: "cat3"},
		{UserID: userID + uint(1), Name: "cat4"},
	}

	db := testutil.SetupDB()
	db.Create(&insertedCategories)

	t.Run("user-category exist", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		list, err := repo.GetAvailableByUser(userID)

		assert.Equal(t, 3, len(list))

		// 1st item
		assert.Equal(t, uint(1), list[0].ID)
		assert.Equal(t, insertedCategories[0].Name, list[0].Name)
		assert.Equal(t, uint(0), list[0].UserID)

		// 2nd item
		assert.Equal(t, uint(2), list[1].ID)
		assert.Equal(t, insertedCategories[1].Name, list[1].Name)
		assert.Equal(t, userID, list[1].UserID)

		// 3rd item
		assert.Equal(t, uint(3), list[2].ID)
		assert.Equal(t, insertedCategories[2].Name, list[2].Name)
		assert.Equal(t, userID, list[2].UserID)

		assert.NoError(t, err)
	})

	t.Run("user-category not exist", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		list, err := repo.GetAvailableByUser(userID + uint(2))

		assert.Equal(t, 1, len(list))

		// 1st item
		assert.Equal(t, uint(1), list[0].ID)
		assert.Equal(t, insertedCategories[0].Name, list[0].Name)
		assert.Equal(t, uint(0), list[0].UserID)

		assert.NoError(t, err)
	})
}

func TestCategoryRepository_GetDefaultCategoryID(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		insertedCategories := []expense.CategoryEntity{
			{UserID: 0, Name: "cat1"},
			{UserID: 11, Name: "cat2"},
			{UserID: 0, Name: "cat3", IsDefault: true},
		}

		db := testutil.SetupDB()
		db.Create(&insertedCategories)

		repo := expense.NewCategoryRepository(db)

		id, err := repo.GetDefaultCategoryID()

		assert.Equal(t, uint(3), id)
		assert.NoError(t, err)
	})

	t.Run("not exist", func(t *testing.T) {
		insertedCategories := []expense.CategoryEntity{
			{UserID: 0, Name: "cat1"},
			{UserID: 11, Name: "cat2"},
		}

		db := testutil.SetupDB()
		db.Create(&insertedCategories)

		repo := expense.NewCategoryRepository(db)

		id, err := repo.GetDefaultCategoryID()

		assert.Equal(t, uint(0), id)
		assert.Error(t, err)
	})
}

func TestCategoryRepository_IsOwner(t *testing.T) {
	var id uint = 1
	var userID uint = 11

	insertedCategories := []expense.CategoryEntity{
		{UserID: userID, Name: "cat1"},
	}

	db := testutil.SetupDB()
	db.Create(&insertedCategories)

	t.Run("owner", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		result, err := repo.IsOwner(id, userID)

		assert.True(t, result)
		assert.NoError(t, err)
	})

	t.Run("not owner", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		result, err := repo.IsOwner(id, userID+uint(1))

		assert.False(t, result)
		assert.NoError(t, err)
	})
}

func TestCategoryRepository_ExistsByName(t *testing.T) {
	var userID uint = 11

	insertedCategories := []expense.CategoryEntity{
		{UserID: 0, Name: "cat1"},
		{UserID: userID, Name: "cat2"},
		{UserID: userID + uint(1), Name: "cat3"},
	}

	db := testutil.SetupDB()
	db.Create(&insertedCategories)

	t.Run("owner", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		// owner: should be true
		result, err := repo.ExistsByName(insertedCategories[1].Name, userID)

		assert.True(t, result)
		assert.NoError(t, err)
	})

	t.Run("system", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		// owner: should be true
		result, err := repo.ExistsByName(insertedCategories[0].Name, uint(0))

		assert.True(t, result)
		assert.NoError(t, err)
	})

	t.Run("another user", func(t *testing.T) {
		repo := expense.NewCategoryRepository(db)

		// owner: should be true
		result, err := repo.ExistsByName(insertedCategories[1].Name, userID+uint(1))

		assert.False(t, result)
		assert.NoError(t, err)
	})
}

package database_test

import (
	"testing"
	"time"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase_Count(t *testing.T) {
	// Create test data
	testUsers := createTestUsers(t)

	t.Run("count all records", func(t *testing.T) {
		var count int64
		err := database.Database[*TestUser]().Count(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(len(testUsers)))
	})

	t.Run("count with query condition", func(t *testing.T) {
		var count int64
		err := database.Database[*TestUser]().WithQuery(&TestUser{IsActive: true}).Count(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(2)) // At least 2 active users from test data
	})

	t.Run("count with raw query", func(t *testing.T) {
		var count int64
		err := database.Database[*TestUser]().WithQueryRaw("age > ?", 25).Count(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1)) // At least 1 user with age > 25
	})

	t.Run("count empty result", func(t *testing.T) {
		var count int64
		err := database.Database[*TestUser]().WithQuery(&TestUser{Name: "NonExistentUser"}).Count(&count)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("test Product model Count", func(t *testing.T) {
		// Create test product
		now := time.Now()
		product := &TestProduct{
			Base: model.Base{
				ID:        "count-product-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Count Test Product",
			Description: "This is a count test product",
			Price:       299.99,
			CategoryID:  "electronics",
		}

		err := database.Database[*TestProduct]().Create(product)
		require.NoError(t, err)

		var count int64
		err = database.Database[*TestProduct]().Count(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))
	})

	t.Run("test Category model Count", func(t *testing.T) {
		// Create test category
		now := time.Now()
		category := &TestCategory{
			Base: model.Base{
				ID:        "count-category-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Count Test Category",
			ParentID: "",
		}

		err := database.Database[*TestCategory]().Create(category)
		require.NoError(t, err)

		var count int64
		err = database.Database[*TestCategory]().Count(&count)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))
	})
}

func TestDatabase_First(t *testing.T) {
	// Create test data
	createTestUsers(t)

	t.Run("get first record", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().First(&user)
		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Name)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
	})

	t.Run("first with query condition", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().WithQuery(&TestUser{IsActive: true}).First(&user)
		require.NoError(t, err)
		assert.True(t, user.IsActive)
	})

	t.Run("first with order", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().WithOrder("name").First(&user)
		require.NoError(t, err)
		assert.NotEmpty(t, user.Name)
	})

	t.Run("first with no results", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().WithQuery(&TestUser{Name: "NonExistentUser"}).First(&user)
		assert.Error(t, err) // Should return error when no records found
		assert.Empty(t, user.ID)
	})

	t.Run("first with nil destination", func(t *testing.T) {
		err := database.Database[*TestUser]().First(nil)
		assert.Error(t, err) // Should return error for nil destination
	})

	t.Run("test Product model First", func(t *testing.T) {
		// Create test product
		now := time.Now()
		product := &TestProduct{
			Base: model.Base{
				ID:        "first-product-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "First Test Product",
			Description: "This is a first test product",
			Price:       399.99,
			CategoryID:  "electronics",
		}

		err := database.Database[*TestProduct]().Create(product)
		require.NoError(t, err)

		var retrievedProduct TestProduct
		err = database.Database[*TestProduct]().First(&retrievedProduct)
		require.NoError(t, err)
		assert.NotEmpty(t, retrievedProduct.ID)
		assert.NotEmpty(t, retrievedProduct.Name)
	})

	t.Run("test Category model First", func(t *testing.T) {
		// Create test category
		now := time.Now()
		category := &TestCategory{
			Base: model.Base{
				ID:        "first-category-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "First Test Category",
			ParentID: "",
		}

		err := database.Database[*TestCategory]().Create(category)
		require.NoError(t, err)

		var retrievedCategory TestCategory
		err = database.Database[*TestCategory]().First(&retrievedCategory)
		require.NoError(t, err)
		assert.NotEmpty(t, retrievedCategory.ID)
		assert.NotEmpty(t, retrievedCategory.Name)
	})
}

func TestDatabase_Last(t *testing.T) {
	// Create test data
	createTestUsers(t)

	t.Run("get last record", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().Last(&user)
		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Name)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
	})

	t.Run("last with query condition", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().WithQuery(&TestUser{IsActive: true}).Last(&user)
		require.NoError(t, err)
		assert.True(t, user.IsActive)
	})

	t.Run("last with order", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().WithOrder("name desc").Last(&user)
		require.NoError(t, err)
		assert.NotEmpty(t, user.Name)
	})

	t.Run("last with no results", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().WithQuery(&TestUser{Name: "NonExistentUser"}).Last(&user)
		assert.Error(t, err) // Should return error when no records found
		assert.Empty(t, user.ID)
	})

	t.Run("last with nil destination", func(t *testing.T) {
		err := database.Database[*TestUser]().Last(nil)
		assert.Error(t, err) // Should return error for nil destination
	})

	t.Run("test Product model Last", func(t *testing.T) {
		// Create test product
		now := time.Now()
		product := &TestProduct{
			Base: model.Base{
				ID:        "last-product-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Last Test Product",
			Description: "This is a last test product",
			Price:       499.99,
			CategoryID:  "electronics",
		}

		err := database.Database[*TestProduct]().Create(product)
		require.NoError(t, err)

		var retrievedProduct TestProduct
		err = database.Database[*TestProduct]().Last(&retrievedProduct)
		require.NoError(t, err)
		assert.NotEmpty(t, retrievedProduct.ID)
		assert.NotEmpty(t, retrievedProduct.Name)
	})

	t.Run("test Category model Last", func(t *testing.T) {
		// Create test category
		now := time.Now()
		category := &TestCategory{
			Base: model.Base{
				ID:        "last-category-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Last Test Category",
			ParentID: "",
		}

		err := database.Database[*TestCategory]().Create(category)
		require.NoError(t, err)

		var retrievedCategory TestCategory
		err = database.Database[*TestCategory]().Last(&retrievedCategory)
		require.NoError(t, err)
		assert.NotEmpty(t, retrievedCategory.ID)
		assert.NotEmpty(t, retrievedCategory.Name)
	})
}

func TestDatabase_Take(t *testing.T) {
	// Create test data
	testUsers := createTestUsers(t)

	t.Run("take first record found", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().Take(&user)
		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Name)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
	})

	t.Run("take with query condition", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().WithQuery(&TestUser{IsActive: true}).Take(&user)
		require.NoError(t, err)
		assert.True(t, user.IsActive)
	})

	t.Run("take with specific condition", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().WithQueryRaw("name = ?", testUsers[0].Name).Take(&user)
		require.NoError(t, err)
		assert.Equal(t, testUsers[0].Name, user.Name)
		assert.Equal(t, testUsers[0].Email, user.Email)
	})

	t.Run("take with no results", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().WithQuery(&TestUser{Name: "NonExistentUser"}).Take(&user)
		assert.Error(t, err) // Should return error when no records found
		assert.Empty(t, user.ID)
	})

	t.Run("take with nil destination", func(t *testing.T) {
		err := database.Database[*TestUser]().Take(nil)
		assert.Error(t, err) // Should return error for nil destination
	})

	t.Run("test Product model Take", func(t *testing.T) {
		// Create test product
		now := time.Now()
		product := &TestProduct{
			Base: model.Base{
				ID:        "take-product-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Take Test Product",
			Description: "This is a take test product",
			Price:       599.99,
			CategoryID:  "electronics",
		}

		err := database.Database[*TestProduct]().Create(product)
		require.NoError(t, err)

		var retrievedProduct TestProduct
		err = database.Database[*TestProduct]().WithQueryRaw("name = ?", "Take Test Product").Take(&retrievedProduct)
		require.NoError(t, err)
		assert.Equal(t, "Take Test Product", retrievedProduct.Name)
		assert.Equal(t, 599.99, retrievedProduct.Price)
	})

	t.Run("test Category model Take", func(t *testing.T) {
		// Create test category
		now := time.Now()
		category := &TestCategory{
			Base: model.Base{
				ID:        "take-category-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Take Test Category",
			ParentID: "",
		}

		err := database.Database[*TestCategory]().Create(category)
		require.NoError(t, err)

		var retrievedCategory TestCategory
		err = database.Database[*TestCategory]().WithQueryRaw("name = ?", "Take Test Category").Take(&retrievedCategory)
		require.NoError(t, err)
		assert.Equal(t, "Take Test Category", retrievedCategory.Name)
		assert.Equal(t, "", retrievedCategory.ParentID)
	})
}

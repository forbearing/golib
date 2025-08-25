package database_test

import (
	"testing"
	"time"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase_Get(t *testing.T) {
	// Create test data
	testUsers := createTestUsers(t)

	t.Run("successfully get existing record", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().Get(&user, "user-001")

		assert.NoError(t, err)
		assert.Equal(t, "user-001", user.ID)
		assert.Equal(t, "Zhang San", user.Name)
		assert.Equal(t, "zhangsan@example.com", user.Email)
		assert.Equal(t, 25, user.Age)
		assert.True(t, user.IsActive)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
	})

	t.Run("get non-existent record", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().Get(&user, "non-existent-id")

		// Should have no error, but user object should be empty
		assert.NoError(t, err)
		assert.Empty(t, user.ID)
		assert.Empty(t, user.Name)
		assert.Empty(t, user.Email)
	})

	t.Run("get record with empty ID", func(t *testing.T) {
		var user TestUser
		err := database.Database[*TestUser]().Get(&user, "")

		assert.NoError(t, err)
		assert.Empty(t, user.ID)
	})

	t.Run("pass nil pointer", func(t *testing.T) {
		err := database.Database[*TestUser]().Get(nil, "user-001")

		// Should return error because nil pointer was passed
		assert.Error(t, err)
	})

	t.Run("verify all test users can be retrieved correctly", func(t *testing.T) {
		for _, expectedUser := range testUsers {
			var user TestUser
			err := database.Database[*TestUser]().Get(&user, expectedUser.ID)

			assert.NoError(t, err)
			assert.Equal(t, expectedUser.ID, user.ID)
			assert.Equal(t, expectedUser.Name, user.Name)
			assert.Equal(t, expectedUser.Email, user.Email)
			assert.Equal(t, expectedUser.Age, user.Age)
			assert.Equal(t, expectedUser.IsActive, user.IsActive)
		}
	})

	t.Run("test Product model", func(t *testing.T) {
		// Create test product
		now := time.Now()
		product := &TestProduct{
			Base: model.Base{
				ID:        "product-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       99.99,
			CategoryID:  "category-001",
		}

		err := database.Database[*TestProduct]().Create(product)
		require.NoError(t, err)

		// Get product
		var retrievedProduct TestProduct
		err = database.Database[*TestProduct]().Get(&retrievedProduct, "product-001")

		assert.NoError(t, err)
		assert.Equal(t, "product-001", retrievedProduct.ID)
		assert.Equal(t, "Test Product", retrievedProduct.Name)
		assert.Equal(t, "This is a test product", retrievedProduct.Description)
		assert.Equal(t, 99.99, retrievedProduct.Price)
		assert.Equal(t, "category-001", retrievedProduct.CategoryID)
	})

	t.Run("test Category model", func(t *testing.T) {
		// Create test category
		now := time.Now()
		category := &TestCategory{
			Base: model.Base{
				ID:        "category-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Electronics",
			ParentID: "",
		}

		err := database.Database[*TestCategory]().Create(category)
		require.NoError(t, err)

		// Get category
		var retrievedCategory TestCategory
		err = database.Database[*TestCategory]().Get(&retrievedCategory, "category-001")

		assert.NoError(t, err)
		assert.Equal(t, "category-001", retrievedCategory.ID)
		assert.Equal(t, "Electronics", retrievedCategory.Name)
		assert.Empty(t, retrievedCategory.ParentID)
	})
}

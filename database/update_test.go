package database_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase_Update(t *testing.T) {
	// Create test data
	testUsers := createTestUsers(t)

	t.Run("successfully update existing record", func(t *testing.T) {
		// Modify first user
		testUsers[0].Name = "Zhang San Updated"
		testUsers[0].Email = "zhangsan_updated@example.com"
		testUsers[0].Age = 26
		// Note: Due to GORM Save method not updating zero value fields (like false),
		// we update other fields first, then update boolean field separately

		err := database.Database[*TestUser](nil).Update(testUsers[0])
		assert.NoError(t, err)

		// Update boolean field to false separately
		err = database.Database[*TestUser](nil).UpdateById(testUsers[0].ID, "is_active", false)
		assert.NoError(t, err)

		// Verify update was successful
		var updatedUser TestUser
		err = database.Database[*TestUser](nil).Get(&updatedUser, testUsers[0].ID)
		assert.NoError(t, err)
		assert.Equal(t, "Zhang San Updated", updatedUser.Name)
		assert.Equal(t, "zhangsan_updated@example.com", updatedUser.Email)
		assert.Equal(t, 26, updatedUser.Age)
		assert.Equal(t, false, updatedUser.IsActive)
		assert.NotNil(t, updatedUser.UpdatedAt)
		// UpdatedAt should be updated
		assert.True(t, updatedUser.UpdatedAt.After(*testUsers[0].CreatedAt))
	})

	t.Run("batch update multiple records", func(t *testing.T) {
		// Create additional test users for batch update
		now := time.Now()
		batchUsers := []*TestUser{
			{
				Base: model.Base{
					ID:        "batch-update-user-001",
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Batch Update User 1",
				Email:    "batchupdate1@example.com",
				Age:      25,
				IsActive: true,
			},
			{
				Base: model.Base{
					ID:        "batch-update-user-002",
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Batch Update User 2",
				Email:    "batchupdate2@example.com",
				Age:      30,
				IsActive: true,
			},
		}

		err := database.Database[*TestUser](nil).Create(batchUsers...)
		require.NoError(t, err)

		// Modify user information
		batchUsers[0].Name = "Batch Update User 1 - Updated"
		batchUsers[0].Age = 27
		batchUsers[1].Name = "Batch Update User 2 - Updated"
		batchUsers[1].Age = 32

		// Batch update
		err = database.Database[*TestUser](nil).Update(batchUsers...)
		assert.NoError(t, err)

		// Verify all records were updated
		for _, user := range batchUsers {
			var retrievedUser TestUser
			err = database.Database[*TestUser](nil).Get(&retrievedUser, user.ID)
			assert.NoError(t, err)
			assert.Equal(t, user.Name, retrievedUser.Name)
			assert.Equal(t, user.Age, retrievedUser.Age)
		}
	})

	t.Run("update non-existent record", func(t *testing.T) {
		now := time.Now()
		nonExistentUser := &TestUser{
			Base: model.Base{
				ID:        "non-existent-update-user",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Non-existent User",
			Email:    "nonexistent@example.com",
			Age:      25,
			IsActive: true,
		}

		// Update method uses GORM's Save, which performs upsert operation
		// If record doesn't exist, it will create a new record
		err := database.Database[*TestUser](nil).Update(nonExistentUser)
		assert.NoError(t, err)

		// Verify record was created
		var user TestUser
		err = database.Database[*TestUser](nil).Get(&user, "non-existent-update-user")
		assert.NoError(t, err)
		assert.Equal(t, "Non-existent User", user.Name)
		assert.Equal(t, "nonexistent@example.com", user.Email)
		assert.Equal(t, 25, user.Age)
		assert.True(t, user.IsActive)
	})

	t.Run("update empty list", func(t *testing.T) {
		// Passing empty list should not error
		err := database.Database[*TestUser](nil).Update()
		assert.NoError(t, err)
	})

	t.Run("control batch update size with WithBatchSize", func(t *testing.T) {
		// Create multiple test users
		now := time.Now()
		var batchSizeUsers []*TestUser
		for i := 0; i < 5; i++ {
			user := &TestUser{
				Base: model.Base{
					ID:        fmt.Sprintf("batch-size-update-user-%03d", i),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     fmt.Sprintf("Batch Size Update User %d", i),
				Email:    fmt.Sprintf("batchsizeupdate%d@example.com", i),
				Age:      20 + i,
				IsActive: true,
			}
			batchSizeUsers = append(batchSizeUsers, user)
		}

		err := database.Database[*TestUser](nil).Create(batchSizeUsers...)
		require.NoError(t, err)

		// Modify all users' age
		for _, user := range batchSizeUsers {
			user.Age += 10
			user.Name += " - Updated"
		}

		// Update with small batch size
		err = database.Database[*TestUser](nil).WithBatchSize(2).Update(batchSizeUsers...)
		assert.NoError(t, err)

		// Verify all records were updated
		for _, user := range batchSizeUsers {
			var retrievedUser TestUser
			err = database.Database[*TestUser](nil).Get(&retrievedUser, user.ID)
			assert.NoError(t, err)
			assert.Equal(t, user.Name, retrievedUser.Name)
			assert.Equal(t, user.Age, retrievedUser.Age)
		}
	})

	t.Run("test Product model update", func(t *testing.T) {
		// Create test product
		now := time.Now()
		product := &TestProduct{
			Base: model.Base{
				ID:        "update-product-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Product to Update",
			Description: "This is a test product to be updated",
			Price:       199.99,
			CategoryID:  "category-001",
		}

		err := database.Database[*TestProduct](nil).Create(product)
		require.NoError(t, err)

		// Update product information
		product.Name = "Updated Product"
		product.Description = "This is an updated test product"
		product.Price = 299.99

		err = database.Database[*TestProduct](nil).Update(product)
		assert.NoError(t, err)

		// Verify product was updated
		var retrievedProduct TestProduct
		err = database.Database[*TestProduct](nil).Get(&retrievedProduct, "update-product-001")
		assert.NoError(t, err)
		assert.Equal(t, "Updated Product", retrievedProduct.Name)
		assert.Equal(t, "This is an updated test product", retrievedProduct.Description)
		assert.Equal(t, 299.99, retrievedProduct.Price)
	})

	t.Run("test Category model update", func(t *testing.T) {
		// Create test category
		now := time.Now()
		category := &TestCategory{
			Base: model.Base{
				ID:        "update-category-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Category to Update",
			ParentID: "",
		}

		err := database.Database[*TestCategory](nil).Create(category)
		require.NoError(t, err)

		// Update category information
		category.Name = "Updated Category"
		category.ParentID = "parent-category-001"

		err = database.Database[*TestCategory](nil).Update(category)
		assert.NoError(t, err)

		// Verify category was updated
		var retrievedCategory TestCategory
		err = database.Database[*TestCategory](nil).Get(&retrievedCategory, "update-category-001")
		assert.NoError(t, err)
		assert.Equal(t, "Updated Category", retrievedCategory.Name)
		assert.Equal(t, "parent-category-001", retrievedCategory.ParentID)
	})

	t.Run("test partial field update", func(t *testing.T) {
		// Create test user
		now := time.Now()
		user := &TestUser{
			Base: model.Base{
				ID:        "partial-update-user",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Partial Update User",
			Email:    "partialupdate@example.com",
			Age:      25,
			IsActive: true,
		}

		err := database.Database[*TestUser](nil).Create(user)
		require.NoError(t, err)

		// Update only some fields
		originalEmail := user.Email
		originalIsActive := user.IsActive
		user.Name = "Partial Update User - Updated"
		user.Age = 30
		// Email and IsActive remain unchanged

		err = database.Database[*TestUser](nil).Update(user)
		assert.NoError(t, err)

		// Verify update results
		var retrievedUser TestUser
		err = database.Database[*TestUser](nil).Get(&retrievedUser, "partial-update-user")
		assert.NoError(t, err)
		assert.Equal(t, "Partial Update User - Updated", retrievedUser.Name)
		assert.Equal(t, 30, retrievedUser.Age)
		assert.Equal(t, originalEmail, retrievedUser.Email)
		assert.Equal(t, originalIsActive, retrievedUser.IsActive)
	})
}

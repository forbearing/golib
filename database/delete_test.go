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

func TestDatabase_Delete(t *testing.T) {
	// Create test data
	testUsers := createTestUsers(t)

	t.Run("successfully delete existing record", func(t *testing.T) {
		// Delete first user
		err := database.Database[*TestUser](nil).Delete(testUsers[0])
		assert.NoError(t, err)

		// Verify record is soft deleted (deleted_at is not null)
		var user TestUser
		err = database.Database[*TestUser](nil).Get(&user, testUsers[0].ID)
		assert.NoError(t, err)
		assert.Empty(t, user.ID) // Cannot be retrieved through normal query after soft delete

		// Use Unscoped query to verify record still exists but marked as deleted
		err = database.Database[*TestUser](nil).WithDB(database.DB.Unscoped()).Get(&user, testUsers[0].ID)
		assert.NoError(t, err)
		assert.Equal(t, testUsers[0].ID, user.ID)
		assert.NotNil(t, user.DeletedAt) // Confirm deleted_at field is set
	})

	t.Run("permanently delete record with WithPurge", func(t *testing.T) {
		// Permanently delete second user
		err := database.Database[*TestUser](nil).WithPurge().Delete(testUsers[1])
		assert.NoError(t, err)

		// Verify record is permanently deleted
		var user TestUser
		err = database.Database[*TestUser](nil).Get(&user, testUsers[1].ID)
		assert.NoError(t, err)
		assert.Empty(t, user.ID)

		// Even with Unscoped, record cannot be found
		err = database.Database[*TestUser](nil).WithDB(database.DB.Unscoped()).Get(&user, testUsers[1].ID)
		assert.NoError(t, err)
		assert.Empty(t, user.ID)
	})

	t.Run("batch delete multiple records", func(t *testing.T) {
		// Create additional test users for batch deletion
		now := time.Now()
		batchUsers := []*TestUser{
			{
				Base: model.Base{
					ID:        "batch-user-001",
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Batch User 1",
				Email:    "batch1@example.com",
				Age:      25,
				IsActive: true,
			},
			{
				Base: model.Base{
					ID:        "batch-user-002",
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Batch User 2",
				Email:    "batch2@example.com",
				Age:      30,
				IsActive: true,
			},
			{
				Base: model.Base{
					ID:        "batch-user-003",
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Batch User 3",
				Email:    "batch3@example.com",
				Age:      28,
				IsActive: true,
			},
		}

		err := database.Database[*TestUser](nil).Create(batchUsers...)
		require.NoError(t, err)

		// Batch delete
		err = database.Database[*TestUser](nil).Delete(batchUsers...)
		assert.NoError(t, err)

		// Verify all records are deleted
		for _, user := range batchUsers {
			var retrievedUser TestUser
			err = database.Database[*TestUser](nil).Get(&retrievedUser, user.ID)
			assert.NoError(t, err)
			assert.Empty(t, retrievedUser.ID)
		}
	})

	t.Run("delete non-existent record", func(t *testing.T) {
		nonExistentUser := &TestUser{
			Base: model.Base{
				ID: "non-existent-user",
			},
		}

		// Deleting non-existent record should not error
		err := database.Database[*TestUser](nil).Delete(nonExistentUser)
		assert.NoError(t, err)
	})

	t.Run("delete empty list", func(t *testing.T) {
		// Passing empty list should not error
		err := database.Database[*TestUser](nil).Delete()
		assert.NoError(t, err)
	})

	t.Run("control batch delete size with WithBatchSize", func(t *testing.T) {
		// Create multiple test users
		now := time.Now()
		var batchSizeUsers []*TestUser
		for i := 0; i < 5; i++ {
			user := &TestUser{
				Base: model.Base{
					ID:        fmt.Sprintf("batch-size-user-%03d", i),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     fmt.Sprintf("Batch Size User %d", i),
				Email:    fmt.Sprintf("batchsize%d@example.com", i),
				Age:      20 + i,
				IsActive: true,
			}
			batchSizeUsers = append(batchSizeUsers, user)
		}

		err := database.Database[*TestUser](nil).Create(batchSizeUsers...)
		require.NoError(t, err)

		// Delete with small batch size
		err = database.Database[*TestUser](nil).WithBatchSize(2).Delete(batchSizeUsers...)
		assert.NoError(t, err)

		// Verify all records are deleted
		for _, user := range batchSizeUsers {
			var retrievedUser TestUser
			err = database.Database[*TestUser](nil).Get(&retrievedUser, user.ID)
			assert.NoError(t, err)
			assert.Empty(t, retrievedUser.ID)
		}
	})

	t.Run("test Product model deletion", func(t *testing.T) {
		// Create test product
		now := time.Now()
		product := &TestProduct{
			Base: model.Base{
				ID:        "delete-product-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Product to Delete",
			Description: "This is a test product to be deleted",
			Price:       199.99,
			CategoryID:  "category-001",
		}

		err := database.Database[*TestProduct](nil).Create(product)
		require.NoError(t, err)

		// Delete product
		err = database.Database[*TestProduct](nil).Delete(product)
		assert.NoError(t, err)

		// Verify product is deleted
		var retrievedProduct TestProduct
		err = database.Database[*TestProduct](nil).Get(&retrievedProduct, "delete-product-001")
		assert.NoError(t, err)
		assert.Empty(t, retrievedProduct.ID)
	})

	t.Run("test Category model deletion", func(t *testing.T) {
		// Create test category
		now := time.Now()
		category := &TestCategory{
			Base: model.Base{
				ID:        "delete-category-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Category to Delete",
			ParentID: "",
		}

		err := database.Database[*TestCategory](nil).Create(category)
		require.NoError(t, err)

		// Delete category
		err = database.Database[*TestCategory](nil).Delete(category)
		assert.NoError(t, err)

		// Verify category is deleted
		var retrievedCategory TestCategory
		err = database.Database[*TestCategory](nil).Get(&retrievedCategory, "delete-category-001")
		assert.NoError(t, err)
		assert.Empty(t, retrievedCategory.ID)
	})
}

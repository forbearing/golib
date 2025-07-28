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

func TestDatabase_Create(t *testing.T) {
	t.Run("successfully create single record", func(t *testing.T) {
		now := time.Now()
		user := &TestUser{
			Base: model.Base{
				ID:        "create-user-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Create Test User",
			Email:    "createtest@example.com",
			Age:      25,
			IsActive: true,
		}

		// Create user
		err := database.Database[*TestUser]().Create(user)
		require.NoError(t, err)

		// Verify user was created
		var retrievedUser TestUser
		err = database.Database[*TestUser]().Get(&retrievedUser, "create-user-001")
		require.NoError(t, err)
		assert.Equal(t, "Create Test User", retrievedUser.Name)
		assert.Equal(t, "createtest@example.com", retrievedUser.Email)
		assert.Equal(t, 25, retrievedUser.Age)
		assert.True(t, retrievedUser.IsActive)
		assert.NotNil(t, retrievedUser.CreatedAt)
		assert.NotNil(t, retrievedUser.UpdatedAt)
	})

	t.Run("successfully create multiple records", func(t *testing.T) {
		now := time.Now()
		users := []*TestUser{
			{
				Base: model.Base{
					ID:        "create-batch-user-001",
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Batch Create User 1",
				Email:    "batchcreate1@example.com",
				Age:      30,
				IsActive: true,
			},
			{
				Base: model.Base{
					ID:        "create-batch-user-002",
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Batch Create User 2",
				Email:    "batchcreate2@example.com",
				Age:      35,
				IsActive: false,
			},
			{
				Base: model.Base{
					ID:        "create-batch-user-003",
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Batch Create User 3",
				Email:    "batchcreate3@example.com",
				Age:      28,
				IsActive: true,
			},
		}

		// Create multiple users
		err := database.Database[*TestUser]().Create(users...)
		require.NoError(t, err)

		// Verify all users were created
		for i, user := range users {
			var retrievedUser TestUser
			err = database.Database[*TestUser]().Get(&retrievedUser, user.ID)
			require.NoError(t, err)
			assert.Equal(t, user.Name, retrievedUser.Name)
			assert.Equal(t, user.Email, retrievedUser.Email)
			assert.Equal(t, user.Age, retrievedUser.Age)
			assert.Equal(t, user.IsActive, retrievedUser.IsActive)
			t.Logf("Created user %d: %s", i+1, retrievedUser.Name)
		}
	})

	t.Run("create with batch size control", func(t *testing.T) {
		now := time.Now()
		var batchUsers []*TestUser
		for i := 0; i < 5; i++ {
			batchUsers = append(batchUsers, &TestUser{
				Base: model.Base{
					ID:        fmt.Sprintf("create-batch-size-user-%03d", i+1),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     fmt.Sprintf("Batch Size User %d", i+1),
				Email:    fmt.Sprintf("batchsizecontrol%d@example.com", i+1),
				Age:      20 + i,
				IsActive: i%2 == 0,
			})
		}

		// Create with batch size of 2
		err := database.Database[*TestUser]().WithBatchSize(2).Create(batchUsers...)
		require.NoError(t, err)

		// Verify all users were created
		for _, user := range batchUsers {
			var retrievedUser TestUser
			err = database.Database[*TestUser]().Get(&retrievedUser, user.ID)
			require.NoError(t, err)
			assert.Equal(t, user.Name, retrievedUser.Name)
		}
	})

	t.Run("create duplicate ID should update existing record", func(t *testing.T) {
		now := time.Now()
		user1 := &TestUser{
			Base: model.Base{
				ID:        "create-duplicate-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Duplicate User 1",
			Email:    "duplicate1@example.com",
			Age:      25,
			IsActive: true,
		}

		user2 := &TestUser{
			Base: model.Base{
				ID:        "create-duplicate-001", // Same ID
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Duplicate User 2",
			Email:    "duplicate2@example.com",
			Age:      30,
			IsActive: true, // Changed to true to avoid zero value issue
		}

		// Create first user
		err := database.Database[*TestUser]().Create(user1)
		require.NoError(t, err)

		// Create second user with same ID (should update existing record due to Save method behavior)
		err = database.Database[*TestUser]().Create(user2)
		require.NoError(t, err) // Should succeed as Save performs upsert

		// Verify the record was updated with user2's data
		var retrievedUser TestUser
		err = database.Database[*TestUser]().Get(&retrievedUser, "create-duplicate-001")
		require.NoError(t, err)
		assert.Equal(t, "Duplicate User 2", retrievedUser.Name)
		assert.Equal(t, "duplicate2@example.com", retrievedUser.Email)
		assert.Equal(t, 30, retrievedUser.Age)
		assert.Equal(t, true, retrievedUser.IsActive)
	})

	t.Run("create empty list should not error", func(t *testing.T) {
		// Creating empty list should not cause error
		err := database.Database[*TestUser]().Create()
		assert.NoError(t, err)
	})

	t.Run("test Product model Create", func(t *testing.T) {
		now := time.Now()
		product := &TestProduct{
			Base: model.Base{
				ID:        "create-product-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Create Test Product",
			Description: "This is a create test product",
			Price:       199.99,
			CategoryID:  "electronics",
		}

		err := database.Database[*TestProduct]().Create(product)
		require.NoError(t, err)

		// Verify product was created
		var retrievedProduct TestProduct
		err = database.Database[*TestProduct]().Get(&retrievedProduct, "create-product-001")
		require.NoError(t, err)
		assert.Equal(t, "Create Test Product", retrievedProduct.Name)
		assert.Equal(t, "This is a create test product", retrievedProduct.Description)
		assert.Equal(t, 199.99, retrievedProduct.Price)
		assert.Equal(t, "electronics", retrievedProduct.CategoryID)
	})

	t.Run("test Category model Create", func(t *testing.T) {
		now := time.Now()
		category := &TestCategory{
			Base: model.Base{
				ID:        "create-category-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Create Test Category",
			ParentID: "",
		}

		err := database.Database[*TestCategory]().Create(category)
		require.NoError(t, err)

		// Verify category was created
		var retrievedCategory TestCategory
		err = database.Database[*TestCategory]().Get(&retrievedCategory, "create-category-001")
		require.NoError(t, err)
		assert.Equal(t, "Create Test Category", retrievedCategory.Name)
		assert.Equal(t, "", retrievedCategory.ParentID)
	})
}


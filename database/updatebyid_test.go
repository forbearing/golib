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

// BenchmarkDatabase_UpdateById performance test for UpdateById operation
func BenchmarkDatabase_UpdateById(b *testing.B) {
	b.Run("update_single_field", func(b *testing.B) {
		for b.Loop() {
			// Create test data
			now := time.Now()
			user := &TestUser{
				Base: model.Base{
					ID:        fmt.Sprintf("benchmark-updatebyid-user-%d", b.N),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Benchmark UpdateById User",
				Email:    fmt.Sprintf("benchmarkupdatebyid%d@example.com", b.N),
				Age:      30,
				IsActive: true,
			}

			err := database.Database[*TestUser](nil).Create(user)
			if err != nil {
				b.Fatal(err)
			}

			// Execute UpdateById
			if err = database.Database[*TestUser](nil).UpdateById(user.ID, "name", "Benchmark UpdateById User - Updated"); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("update_multiple_fields_sequential", func(b *testing.B) {
		for b.Loop() {
			// Create test data
			now := time.Now()
			user := &TestUser{
				Base: model.Base{
					ID:        fmt.Sprintf("benchmark-updatebyid-multi-user-%d", b.N),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Benchmark UpdateById Multi User",
				Email:    fmt.Sprintf("benchmarkupdatebyidmulti%d@example.com", b.N),
				Age:      30,
				IsActive: true,
			}

			err := database.Database[*TestUser](nil).Create(user)
			if err != nil {
				b.Fatal(err)
			}

			// Update multiple fields sequentially
			if err = database.Database[*TestUser](nil).UpdateById(user.ID, "name", "Benchmark UpdateById Multi User - Updated"); err != nil {
				b.Fatal(err)
			}
			if err = database.Database[*TestUser](nil).UpdateById(user.ID, "age", 35); err != nil {
				b.Fatal(err)
			}
			if err = database.Database[*TestUser](nil).UpdateById(user.ID, "is_active", false); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestDatabase_UpdateById(t *testing.T) {
	// Create test data
	testUsers := createTestUsers(t)

	t.Run("successfully update single field by ID", func(t *testing.T) {
		// Update first user's name
		err := database.Database[*TestUser](nil).UpdateById(testUsers[0].ID, "name", "Zhang San Updated by ID")
		assert.NoError(t, err)

		// Verify update was successful
		var updatedUser TestUser
		err = database.Database[*TestUser](nil).Get(&updatedUser, testUsers[0].ID)
		assert.NoError(t, err)
		assert.Equal(t, "Zhang San Updated by ID", updatedUser.Name)
		// Other fields should remain unchanged
		assert.Equal(t, testUsers[0].Email, updatedUser.Email)
		assert.Equal(t, testUsers[0].Age, updatedUser.Age)
		assert.Equal(t, testUsers[0].IsActive, updatedUser.IsActive)
	})

	t.Run("update multiple different fields", func(t *testing.T) {
		// Update second user's age
		err := database.Database[*TestUser](nil).UpdateById(testUsers[1].ID, "age", 35)
		assert.NoError(t, err)

		// Verify age update
		var updatedUser TestUser
		err = database.Database[*TestUser](nil).Get(&updatedUser, testUsers[1].ID)
		assert.NoError(t, err)
		assert.Equal(t, 35, updatedUser.Age)

		// Update same user's email
		err = database.Database[*TestUser](nil).UpdateById(testUsers[1].ID, "email", "lisi_updated_by_id@example.com")
		assert.NoError(t, err)

		// Verify email update
		err = database.Database[*TestUser](nil).Get(&updatedUser, testUsers[1].ID)
		assert.NoError(t, err)
		assert.Equal(t, "lisi_updated_by_id@example.com", updatedUser.Email)
		assert.Equal(t, 35, updatedUser.Age) // Previous update should persist
	})

	t.Run("update boolean field", func(t *testing.T) {
		// Update third user's active status
		err := database.Database[*TestUser](nil).UpdateById(testUsers[2].ID, "is_active", false)
		assert.NoError(t, err)

		// Verify boolean field update
		var updatedUser TestUser
		err = database.Database[*TestUser](nil).Get(&updatedUser, testUsers[2].ID)
		assert.NoError(t, err)
		assert.False(t, updatedUser.IsActive)
	})

	t.Run("update non-existent record", func(t *testing.T) {
		// Try to update non-existent record
		err := database.Database[*TestUser](nil).UpdateById("non-existent-id", "name", "Non-existent User")
		assert.NoError(t, err) // GORM won't error for non-existent records

		// Verify record indeed doesn't exist
		var user TestUser
		err = database.Database[*TestUser](nil).Get(&user, "non-existent-id")
		assert.NoError(t, err)
		assert.Empty(t, user.ID)
	})

	t.Run("update with empty ID", func(t *testing.T) {
		// Using empty ID should not error but have no effect
		err := database.Database[*TestUser](nil).UpdateById("", "name", "Empty ID User")
		assert.NoError(t, err)
	})

	t.Run("update non-existent field", func(t *testing.T) {
		// Try to update non-existent field, should return error
		err := database.Database[*TestUser](nil).UpdateById(testUsers[0].ID, "non_existent_field", "some value")
		assert.Error(t, err) // Database should error for non-existent field
		assert.Contains(t, err.Error(), "no such column")
	})

	t.Run("test Product model UpdateById", func(t *testing.T) {
		// Create test product
		now := time.Now()
		product := &TestProduct{
			Base: model.Base{
				ID:        "updatebyid-product-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "UpdateById Product",
			Description: "This is an UpdateById test product",
			Price:       199.99,
			CategoryID:  "category-001",
		}

		err := database.Database[*TestProduct](nil).Create(product)
		require.NoError(t, err)

		// Update product price by ID
		err = database.Database[*TestProduct](nil).UpdateById("updatebyid-product-001", "price", 299.99)
		assert.NoError(t, err)

		// Verify product price was updated
		var retrievedProduct TestProduct
		err = database.Database[*TestProduct](nil).Get(&retrievedProduct, "updatebyid-product-001")
		assert.NoError(t, err)
		assert.Equal(t, 299.99, retrievedProduct.Price)
		// Other fields should remain unchanged
		assert.Equal(t, "UpdateById Product", retrievedProduct.Name)
		assert.Equal(t, "This is an UpdateById test product", retrievedProduct.Description)
	})

	t.Run("test Category model UpdateById", func(t *testing.T) {
		// Create test category
		now := time.Now()
		category := &TestCategory{
			Base: model.Base{
				ID:        "updatebyid-category-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "UpdateById Category",
			ParentID: "",
		}

		err := database.Database[*TestCategory](nil).Create(category)
		require.NoError(t, err)

		// Update category's parent ID by ID
		err = database.Database[*TestCategory](nil).UpdateById("updatebyid-category-001", "parent_id", "parent-category-001")
		assert.NoError(t, err)

		// Verify category parent ID was updated
		var retrievedCategory TestCategory
		err = database.Database[*TestCategory](nil).Get(&retrievedCategory, "updatebyid-category-001")
		assert.NoError(t, err)
		assert.Equal(t, "parent-category-001", retrievedCategory.ParentID)
		// Name should remain unchanged
		assert.Equal(t, "UpdateById Category", retrievedCategory.Name)
	})
}

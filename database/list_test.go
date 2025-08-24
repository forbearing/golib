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

func TestDatabase_List(t *testing.T) {
	// Create test data
	testUsers := createTestUsers(t)

	t.Run("list all users without conditions", func(t *testing.T) {
		var users []*TestUser
		err := database.Database[*TestUser]().List(&users)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), len(testUsers))

		// Verify that all test users are in the result
		userMap := make(map[string]*TestUser)
		for _, user := range users {
			userMap[user.ID] = user
		}

		for _, expectedUser := range testUsers {
			user, exists := userMap[expectedUser.ID]
			assert.True(t, exists, "User %s should exist in list", expectedUser.ID)
			if exists {
				assert.Equal(t, expectedUser.Name, user.Name)
				assert.Equal(t, expectedUser.Email, user.Email)
				assert.Equal(t, expectedUser.Age, user.Age)
				assert.Equal(t, expectedUser.IsActive, user.IsActive)
			}
		}
	})

	t.Run("list with query conditions", func(t *testing.T) {
		// Query by age
		var users []*TestUser
		query := &TestUser{Age: 25}
		err := database.Database[*TestUser]().WithQuery(query).List(&users)

		assert.NoError(t, err)
		assert.Greater(t, len(users), 0)
		for _, user := range users {
			assert.Equal(t, 25, user.Age)
		}
	})

	t.Run("list with query by name", func(t *testing.T) {
		var users []*TestUser
		query := &TestUser{Name: "Zhang San"}
		err := database.Database[*TestUser]().WithQuery(query).List(&users)

		assert.NoError(t, err)
		assert.Greater(t, len(users), 0)
		for _, user := range users {
			assert.Equal(t, "Zhang San", user.Name)
		}
	})

	t.Run("list with query by active status", func(t *testing.T) {
		var users []*TestUser
		query := &TestUser{IsActive: true}
		err := database.Database[*TestUser]().WithQuery(query).List(&users)

		assert.NoError(t, err)
		assert.Greater(t, len(users), 0)
		for _, user := range users {
			assert.True(t, user.IsActive)
		}
	})

	t.Run("list with fuzzy match", func(t *testing.T) {
		var users []*TestUser
		query := &TestUser{Name: "Zhang"}
		err := database.Database[*TestUser]().WithQuery(query, true).List(&users)

		assert.NoError(t, err)
		assert.Greater(t, len(users), 0)
		for _, user := range users {
			assert.Contains(t, user.Name, "Zhang")
		}
	})

	t.Run("list with order by name ascending", func(t *testing.T) {
		var users []*TestUser
		err := database.Database[*TestUser]().WithOrder("name ASC").List(&users)

		assert.NoError(t, err)
		assert.Greater(t, len(users), 1)

		// Verify ascending order
		for i := 1; i < len(users); i++ {
			assert.LessOrEqual(t, users[i-1].Name, users[i].Name)
		}
	})

	t.Run("list with order by age descending", func(t *testing.T) {
		var users []*TestUser
		err := database.Database[*TestUser]().WithOrder("age DESC").List(&users)

		assert.NoError(t, err)
		assert.Greater(t, len(users), 1)

		// Verify descending order
		for i := 1; i < len(users); i++ {
			assert.GreaterOrEqual(t, users[i-1].Age, users[i].Age)
		}
	})

	t.Run("list with limit", func(t *testing.T) {
		var users []*TestUser
		err := database.Database[*TestUser]().WithLimit(2).List(&users)

		assert.NoError(t, err)
		assert.LessOrEqual(t, len(users), 2)
	})

	t.Run("list with scope pagination", func(t *testing.T) {
		// First page
		var page1Users []*TestUser
		err := database.Database[*TestUser]().WithScope(1, 2).List(&page1Users)

		assert.NoError(t, err)
		assert.LessOrEqual(t, len(page1Users), 2)

		// Second page
		var page2Users []*TestUser
		err = database.Database[*TestUser]().WithScope(2, 2).List(&page2Users)

		assert.NoError(t, err)
		assert.LessOrEqual(t, len(page2Users), 2)

		// Verify no overlap between pages
		if len(page1Users) > 0 && len(page2Users) > 0 {
			page1IDs := make(map[string]bool)
			for _, user := range page1Users {
				page1IDs[user.ID] = true
			}

			for _, user := range page2Users {
				assert.False(t, page1IDs[user.ID], "User %s should not appear in both pages", user.ID)
			}
		}
	})

	t.Run("list with cache enabled", func(t *testing.T) {
		var users1 []*TestUser
		err := database.Database[*TestUser]().WithCache(true).List(&users1)
		assert.NoError(t, err)

		var users2 []*TestUser
		err = database.Database[*TestUser]().WithCache(true).List(&users2)
		assert.NoError(t, err)

		// Results should be the same
		assert.Equal(t, len(users1), len(users2))
	})

	t.Run("list with multiple conditions", func(t *testing.T) {
		var users []*TestUser
		query := &TestUser{Age: 25, IsActive: true}
		err := database.Database[*TestUser]().
			WithQuery(query).
			WithOrder("name ASC").
			WithLimit(10).
			List(&users)

		assert.NoError(t, err)
		assert.LessOrEqual(t, len(users), 10)

		for _, user := range users {
			assert.Equal(t, 25, user.Age)
			assert.True(t, user.IsActive)
		}

		// Verify order
		if len(users) > 1 {
			for i := 1; i < len(users); i++ {
				assert.LessOrEqual(t, users[i-1].Name, users[i].Name)
			}
		}
	})

	t.Run("list empty result", func(t *testing.T) {
		var users []*TestUser
		query := &TestUser{Name: "NonExistentUser"}
		err := database.Database[*TestUser]().WithQuery(query).List(&users)

		assert.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("list with nil pointer", func(t *testing.T) {
		err := database.Database[*TestUser]().List(nil)
		assert.NoError(t, err)
	})
}

func TestDatabase_List_Products(t *testing.T) {
	// Create test products
	now := time.Now()
	products := []*TestProduct{
		{
			Base: model.Base{
				ID:        "product-list-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Laptop",
			Description: "High-performance laptop",
			Price:       1299.99,
			CategoryID:  "electronics",
		},
		{
			Base: model.Base{
				ID:        "product-list-002",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Mouse",
			Description: "Wireless mouse",
			Price:       29.99,
			CategoryID:  "electronics",
		},
		{
			Base: model.Base{
				ID:        "product-list-003",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:        "Book",
			Description: "Programming book",
			Price:       49.99,
			CategoryID:  "books",
		},
	}

	for _, product := range products {
		err := database.Database[*TestProduct]().Create(product)
		require.NoError(t, err)
	}

	t.Run("list all products", func(t *testing.T) {
		var products []*TestProduct
		err := database.Database[*TestProduct]().List(&products)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(products), 3)
	})

	t.Run("list products by category", func(t *testing.T) {
		var products []*TestProduct
		query := &TestProduct{CategoryID: "electronics"}
		err := database.Database[*TestProduct]().WithQuery(query).List(&products)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(products), 2)
		for _, product := range products {
			assert.Equal(t, "electronics", product.CategoryID)
		}
	})

	t.Run("list products ordered by price", func(t *testing.T) {
		var products []*TestProduct
		err := database.Database[*TestProduct]().WithOrder("price ASC").List(&products)

		assert.NoError(t, err)
		assert.Greater(t, len(products), 1)

		// Verify ascending order by price
		for i := 1; i < len(products); i++ {
			assert.LessOrEqual(t, products[i-1].Price, products[i].Price)
		}
	})

	t.Run("list_products_with_name_like_match", func(t *testing.T) {
		var products []*TestProduct
		query := &TestProduct{Name: "Product"}
		err := database.Database[*TestProduct]().WithQuery(query, true).List(&products)
		assert.NoError(t, err)
		assert.Greater(t, len(products), 0, "Should find products with name containing 'Product'")

		// 验证所有返回的产品名称都包含 "Product"
		for _, product := range products {
			assert.Contains(t, product.Name, "Product")
		}
	})
}

func TestDatabase_List_Categories(t *testing.T) {
	// 清理现有数据 - 先查询所有分类，然后删除
	var existingCategories []*TestCategory
	database.Database[*TestCategory]().List(&existingCategories)
	if len(existingCategories) > 0 {
		database.Database[*TestCategory]().WithPurge(true).Delete(existingCategories...)
	}

	// Create test categories with unique names
	now := time.Now()
	timestamp := now.Unix()
	categories := []*TestCategory{
		{
			Base: model.Base{
				ID:        "category-list-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     fmt.Sprintf("Electronics_%d", timestamp),
			ParentID: "",
		},
		{
			Base: model.Base{
				ID:        "category-list-002",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     fmt.Sprintf("Computers_%d", timestamp),
			ParentID: "category-list-001",
		},
		{
			Base: model.Base{
				ID:        "category-list-003",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     fmt.Sprintf("Books_%d", timestamp),
			ParentID: "",
		},
	}

	for _, category := range categories {
		err := database.Database[*TestCategory]().Create(category)
		require.NoError(t, err)
	}

	t.Run("list all categories", func(t *testing.T) {
		var categories []*TestCategory
		err := database.Database[*TestCategory]().List(&categories)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(categories), 3)
	})

	t.Run("list root categories", func(t *testing.T) {
		var categories []*TestCategory
		// 使用 WithQueryRaw 来正确查询 ParentID 为空字符串的记录
		err := database.Database[*TestCategory]().WithQueryRaw("parent_id = ?", "").List(&categories)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(categories), 2)
		for _, category := range categories {
			assert.Empty(t, category.ParentID)
		}
	})

	t.Run("list subcategories", func(t *testing.T) {
		var categories []*TestCategory
		query := &TestCategory{ParentID: "category-list-001"}
		err := database.Database[*TestCategory]().WithQuery(query).List(&categories)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(categories), 1)
		for _, category := range categories {
			assert.Equal(t, "category-list-001", category.ParentID)
		}
	})

	t.Run("list categories ordered by name", func(t *testing.T) {
		var categories []*TestCategory
		err := database.Database[*TestCategory]().WithOrder("name ASC").List(&categories)

		assert.NoError(t, err)
		assert.Greater(t, len(categories), 1)

		// Verify ascending order by name
		for i := 1; i < len(categories); i++ {
			assert.LessOrEqual(t, categories[i-1].Name, categories[i].Name)
		}
	})
}

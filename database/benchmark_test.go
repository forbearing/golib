package database_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
)

// BenchmarkDatabase_Get performance test for Get method
func BenchmarkDatabase_Get(b *testing.B) {
	b.Run("nocache", func(b *testing.B) {
		// Create test data
		now := time.Now()
		user := &TestUser{
			Base: model.Base{
				ID:        "benchmark-user",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Benchmark User",
			Email:    "benchmark@example.com",
			Age:      30,
			IsActive: true,
		}

		err := database.Database[*TestUser](nil).Create(user)
		if err != nil {
			b.Fatal(err)
		}

		for b.Loop() {
			var retrievedUser TestUser
			if err = database.Database[*TestUser](nil).Get(&retrievedUser, "benchmark-user"); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("withcache", func(b *testing.B) {
		// Create test data
		now := time.Now()
		user := &TestUser{
			Base: model.Base{
				ID:        "benchmark-user-cached",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Cached Benchmark User",
			Email:    "benchmark-cached@example.com",
			Age:      30,
			IsActive: true,
		}

		err := database.Database[*TestUser](nil).Create(user)
		if err != nil {
			b.Fatal(err)
		}

		// Warm up cache
		var warmupUser TestUser
		if err = database.Database[*TestUser](nil).Get(&warmupUser, "benchmark-user-cached"); err != nil {
			b.Fatal(err)
		}

		for b.Loop() {
			var retrievedUser TestUser
			if err := database.Database[*TestUser](nil).WithCache(true).Get(&retrievedUser, "benchmark-user-cached"); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDatabase_Create performance test for Create method
func BenchmarkDatabase_Create(b *testing.B) {
	b.Run("single_create", func(b *testing.B) {
		for b.Loop() {
			now := time.Now()
			user := &TestUser{
				Base: model.Base{
					ID:        fmt.Sprintf("benchmark-create-user-%d", b.N),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Benchmark Create User",
				Email:    fmt.Sprintf("benchmarkcreate%d@example.com", b.N),
				Age:      30,
				IsActive: true,
			}

			if err := database.Database[*TestUser](nil).Create(user); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("batch_create", func(b *testing.B) {
		for b.Loop() {
			now := time.Now()
			var users []*TestUser
			for i := 0; i < 10; i++ {
				user := &TestUser{
					Base: model.Base{
						ID:        fmt.Sprintf("benchmark-batch-create-user-%d-%d", b.N, i),
						CreatedAt: &now,
						UpdatedAt: &now,
					},
					Name:     fmt.Sprintf("Batch Create User %d", i),
					Email:    fmt.Sprintf("benchmarkbatchcreate%d-%d@example.com", b.N, i),
					Age:      20 + i,
					IsActive: true,
				}
				users = append(users, user)
			}

			if err := database.Database[*TestUser](nil).Create(users...); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDatabase_List performance test for List method
func BenchmarkDatabase_List(b *testing.B) {
	// Setup test data
	now := time.Now()
	var users []*TestUser
	for i := 0; i < 100; i++ {
		user := &TestUser{
			Base: model.Base{
				ID:        fmt.Sprintf("benchmark-list-user-%d", i),
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     fmt.Sprintf("List User %d", i),
			Email:    fmt.Sprintf("benchmarklist%d@example.com", i),
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
		}
		users = append(users, user)
	}

	err := database.Database[*TestUser](nil).Create(users...)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("list_all", func(b *testing.B) {
		for b.Loop() {
			var result []*TestUser
			if err := database.Database[*TestUser](nil).List(&result); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("list_with_condition", func(b *testing.B) {
		for b.Loop() {
			var result []*TestUser
			if err := database.Database[*TestUser](nil).WithQueryRaw("age > ?", 30).List(&result); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("list_with_limit", func(b *testing.B) {
		for b.Loop() {
			var result []*TestUser
			if err := database.Database[*TestUser](nil).WithLimit(10).List(&result); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDatabase_Count performance test for Count method
func BenchmarkDatabase_Count(b *testing.B) {
	// Setup test data
	now := time.Now()
	var users []*TestUser
	for i := 0; i < 100; i++ {
		user := &TestUser{
			Base: model.Base{
				ID:        fmt.Sprintf("benchmark-count-user-%d", i),
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     fmt.Sprintf("Count User %d", i),
			Email:    fmt.Sprintf("benchmarkcount%d@example.com", i),
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
		}
		users = append(users, user)
	}

	err := database.Database[*TestUser](nil).Create(users...)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("count_all", func(b *testing.B) {
		for b.Loop() {
			var count int64
			if err := database.Database[*TestUser](nil).Count(&count); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("count_with_condition", func(b *testing.B) {
		for b.Loop() {
			var count int64
			if err := database.Database[*TestUser](nil).WithQueryRaw("age > ?", 30).Count(&count); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDatabase_First performance test for First method
func BenchmarkDatabase_First(b *testing.B) {
	// Setup test data
	now := time.Now()
	var users []*TestUser
	for i := 0; i < 100; i++ {
		user := &TestUser{
			Base: model.Base{
				ID:        fmt.Sprintf("benchmark-first-user-%d", i),
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     fmt.Sprintf("First User %d", i),
			Email:    fmt.Sprintf("benchmarkfirst%d@example.com", i),
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
		}
		users = append(users, user)
	}

	err := database.Database[*TestUser](nil).Create(users...)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("first_record", func(b *testing.B) {
		for b.Loop() {
			var user TestUser
			if err := database.Database[*TestUser](nil).First(&user); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("first_with_condition", func(b *testing.B) {
		for b.Loop() {
			var user TestUser
			if err := database.Database[*TestUser](nil).WithQueryRaw("age > ?", 30).First(&user); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDatabase_Last performance test for Last method
func BenchmarkDatabase_Last(b *testing.B) {
	// Setup test data
	now := time.Now()
	var users []*TestUser
	for i := 0; i < 100; i++ {
		user := &TestUser{
			Base: model.Base{
				ID:        fmt.Sprintf("benchmark-last-user-%d", i),
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     fmt.Sprintf("Last User %d", i),
			Email:    fmt.Sprintf("benchmarklast%d@example.com", i),
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
		}
		users = append(users, user)
	}

	err := database.Database[*TestUser](nil).Create(users...)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("last_record", func(b *testing.B) {
		for b.Loop() {
			var user TestUser
			if err := database.Database[*TestUser](nil).Last(&user); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("last_with_condition", func(b *testing.B) {
		for b.Loop() {
			var user TestUser
			if err := database.Database[*TestUser](nil).WithQueryRaw("age > ?", 30).Last(&user); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDatabase_Take performance test for Take method
func BenchmarkDatabase_Take(b *testing.B) {
	// Setup test data
	now := time.Now()
	var users []*TestUser
	for i := 0; i < 100; i++ {
		user := &TestUser{
			Base: model.Base{
				ID:        fmt.Sprintf("benchmark-take-user-%d", i),
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     fmt.Sprintf("Take User %d", i),
			Email:    fmt.Sprintf("benchmarktake%d@example.com", i),
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
		}
		users = append(users, user)
	}

	err := database.Database[*TestUser](nil).Create(users...)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("take_record", func(b *testing.B) {
		for b.Loop() {
			var user TestUser
			if err := database.Database[*TestUser](nil).Take(&user); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("take_with_condition", func(b *testing.B) {
		for b.Loop() {
			var user TestUser
			if err := database.Database[*TestUser](nil).WithQueryRaw("age > ?", 30).Take(&user); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDatabase_Update performance test for Update method
func BenchmarkDatabase_Update(b *testing.B) {
	b.Run("single_update", func(b *testing.B) {
		for b.Loop() {
			// Create test data
			now := time.Now()
			user := &TestUser{
				Base: model.Base{
					ID:        fmt.Sprintf("benchmark-update-user-%d", b.N),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Benchmark Update User",
				Email:    fmt.Sprintf("benchmarkupdate%d@example.com", b.N),
				Age:      30,
				IsActive: true,
			}

			err := database.Database[*TestUser](nil).Create(user)
			if err != nil {
				b.Fatal(err)
			}

			// Modify user information
			user.Name = "Benchmark Update User - Updated"
			user.Age = 31

			// Execute update
			if err = database.Database[*TestUser](nil).Update(user); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("batch_update", func(b *testing.B) {
		for b.Loop() {
			// Create batch test data
			now := time.Now()
			var users []*TestUser
			for i := 0; i < 10; i++ {
				user := &TestUser{
					Base: model.Base{
						ID:        fmt.Sprintf("benchmark-batch-update-user-%d-%d", b.N, i),
						CreatedAt: &now,
						UpdatedAt: &now,
					},
					Name:     fmt.Sprintf("Batch Update User %d", i),
					Email:    fmt.Sprintf("benchmarkbatchupdate%d-%d@example.com", b.N, i),
					Age:      20 + i,
					IsActive: true,
				}
				users = append(users, user)
			}

			err := database.Database[*TestUser](nil).Create(users...)
			if err != nil {
				b.Fatal(err)
			}

			// Modify all user information
			for _, user := range users {
				user.Name += " - Updated"
				user.Age += 5
			}

			// Execute batch update
			if err = database.Database[*TestUser](nil).Update(users...); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDatabase_Delete performance test for Delete method
func BenchmarkDatabase_Delete(b *testing.B) {
	b.Run("soft_delete", func(b *testing.B) {
		for b.Loop() {
			// Create test data
			now := time.Now()
			user := &TestUser{
				Base: model.Base{
					ID:        fmt.Sprintf("benchmark-delete-user-%d", b.N),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Benchmark Delete User",
				Email:    fmt.Sprintf("benchmarkdelete%d@example.com", b.N),
				Age:      30,
				IsActive: true,
			}

			err := database.Database[*TestUser](nil).Create(user)
			if err != nil {
				b.Fatal(err)
			}

			// Execute soft delete
			if err = database.Database[*TestUser](nil).Delete(user); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("purge_delete", func(b *testing.B) {
		for b.Loop() {
			// Create test data
			now := time.Now()
			user := &TestUser{
				Base: model.Base{
					ID:        fmt.Sprintf("benchmark-purge-user-%d", b.N),
					CreatedAt: &now,
					UpdatedAt: &now,
				},
				Name:     "Benchmark Purge User",
				Email:    fmt.Sprintf("benchmarkpurge%d@example.com", b.N),
				Age:      30,
				IsActive: true,
			}

			err := database.Database[*TestUser](nil).Create(user)
			if err != nil {
				b.Fatal(err)
			}

			// Execute permanent delete
			if err = database.Database[*TestUser](nil).WithPurge().Delete(user); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("batch_delete", func(b *testing.B) {
		for b.Loop() {
			// Create batch test data
			now := time.Now()
			var users []*TestUser
			for i := 0; i < 10; i++ {
				user := &TestUser{
					Base: model.Base{
						ID:        fmt.Sprintf("benchmark-batch-user-%d-%d", b.N, i),
						CreatedAt: &now,
						UpdatedAt: &now,
					},
					Name:     fmt.Sprintf("Batch Delete User %d", i),
					Email:    fmt.Sprintf("benchmarkbatch%d-%d@example.com", b.N, i),
					Age:      20 + i,
					IsActive: true,
				}
				users = append(users, user)
			}

			err := database.Database[*TestUser](nil).Create(users...)
			if err != nil {
				b.Fatal(err)
			}

			// Execute batch delete
			if err = database.Database[*TestUser](nil).Delete(users...); err != nil {
				b.Fatal(err)
			}
		}
	})
}

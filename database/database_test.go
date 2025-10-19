package database_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/forbearing/gst/bootstrap"
	"github.com/forbearing/gst/config"
	"github.com/forbearing/gst/database"
	"github.com/forbearing/gst/model"
	"github.com/forbearing/gst/types"
	"github.com/forbearing/gst/types/consts"
	"github.com/stretchr/testify/suite"
)

// TestUser test user model
type TestUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active"`

	model.Base
}

// UpdateBefore sets the UpdatedAt timestamp before update operations
func (u *TestUser) UpdateBefore(_ *types.ModelContext) error {
	now := time.Now()
	u.UpdatedAt = &now
	return nil
}

// TestProduct test product model
type TestProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`

	model.Base
}

// TestCategory test category model
type TestCategory struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`

	model.Base
}

// DatabaseTestSuite defines the test suite for database operations
type DatabaseTestSuite struct {
	suite.Suite
	userDB     types.Database[*TestUser]
	productDB  types.Database[*TestProduct]
	categoryDB types.Database[*TestCategory]
}

// SetupSuite runs once before all tests in the suite
func (suite *DatabaseTestSuite) SetupSuite() {
	suite.userDB = database.Database[*TestUser](nil)
	suite.productDB = database.Database[*TestProduct](nil)
	suite.categoryDB = database.Database[*TestCategory](nil)
}

// SetupTest runs before each test
func (suite *DatabaseTestSuite) SetupTest() {
	// Clean up test data before each test
	_ = suite.userDB.WithPurge().Delete(&TestUser{})
	_ = suite.productDB.WithPurge().Delete(&TestProduct{})
	_ = suite.categoryDB.WithPurge().Delete(&TestCategory{})
}

func init() {
	os.Setenv(config.LOGGER_DIR, "/tmp/test_database")
	os.Setenv(config.DATABASE_TYPE, string(config.DBSqlite))
	os.Setenv(config.SQLITE_IS_MEMORY, "false")
	os.Setenv(config.SQLITE_PATH, "/tmp/test.db")

	_ = os.Remove("/tmp/test.db")

	// os.Setenv(config.DATABASE_TYPE, string(config.DBMySQL))
	// os.Setenv(config.MYSQL_DATABASE, "test")
	// os.Setenv(config.MYSQL_USERNAME, "test")
	// os.Setenv(config.MYSQL_PASSWORD, "test")

	model.Register[*TestUser]()
	model.Register[*TestProduct]()
	model.Register[*TestCategory]()

	if err := bootstrap.Bootstrap(); err != nil {
		panic(err)
	}
}

// TestDatabaseSuite runs the database test suite
func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

// TestDatabase tests the Database factory function
func (suite *DatabaseTestSuite) TestDatabase() {
	// Test creating database instance
	db := database.Database[*TestUser](nil)
	suite.NotNil(db)

	// Test that multiple calls return different instances
	db2 := database.Database[*TestUser](nil)
	suite.NotNil(db2)
	// Note: We can't directly compare instances as they may be different
}

// TestWithDB tests the WithDB method
func (suite *DatabaseTestSuite) TestWithDB() {
	db := suite.userDB

	// Test with valid database
	result := db.WithDB(database.DB)
	suite.NotNil(result)

	// Test with nil database
	result = db.WithDB(nil)
	suite.NotNil(result)
}

// TestWithTable tests the WithTable method
func (suite *DatabaseTestSuite) TestWithTable() {
	db := suite.userDB

	// Test setting custom table name
	result := db.WithTable("custom_users")
	suite.NotNil(result)

	// Test with empty table name
	result = db.WithTable("")
	suite.NotNil(result)
}

// TestWithBatchSize tests the WithBatchSize method
func (suite *DatabaseTestSuite) TestWithBatchSize() {
	db := suite.userDB

	// Test setting batch size
	result := db.WithBatchSize(100)
	suite.NotNil(result)

	// Test with zero batch size
	result = db.WithBatchSize(0)
	suite.NotNil(result)

	// Test with negative batch size
	result = db.WithBatchSize(-1)
	suite.NotNil(result)
}

// TestWithDebug tests the WithDebug method
func (suite *DatabaseTestSuite) TestWithDebug() {
	db := suite.userDB

	// Test enabling debug mode
	result := db.WithDebug()
	suite.NotNil(result)
}

// TestWithAnd tests the WithAnd method
func (suite *DatabaseTestSuite) TestWithAnd() {
	db := suite.userDB

	// Test setting AND mode (default)
	result := db.WithAnd()
	suite.NotNil(result)

	// Test with explicit true
	result = db.WithAnd(true)
	suite.NotNil(result)

	// Test with false
	result = db.WithAnd(false)
	suite.NotNil(result)
}

// TestWithOr tests the WithOr method
func (suite *DatabaseTestSuite) TestWithOr() {
	db := suite.userDB

	// Test setting OR mode
	result := db.WithOr()
	suite.NotNil(result)

	// Test with explicit true
	result = db.WithOr(true)
	suite.NotNil(result)

	// Test with false
	result = db.WithOr(false)
	suite.NotNil(result)
}

// TestWithIndex tests the WithIndex method
func (suite *DatabaseTestSuite) TestWithIndex() {
	db := suite.userDB

	// Test default behavior - single index name defaults to USE INDEX
	result := db.WithIndex("idx_name")
	suite.NotNil(result)

	// Test explicit USE INDEX
	result = db.WithIndex("idx_name", consts.IndexHintUse)
	suite.NotNil(result)

	// Test FORCE INDEX
	result = db.WithIndex("idx_name", consts.IndexHintForce)
	suite.NotNil(result)

	// Test IGNORE INDEX
	result = db.WithIndex("idx_name", consts.IndexHintIgnore)
	suite.NotNil(result)

	// Test with empty index name (should return unchanged)
	result = db.WithIndex("")
	suite.NotNil(result)

	// Test with whitespace-only index name (should return unchanged)
	result = db.WithIndex("  ")
	suite.NotNil(result)
}

// TestWithQuery tests the WithQuery method
func (suite *DatabaseTestSuite) TestWithQuery() {
	db := suite.userDB

	// Test with struct query
	user := &TestUser{Name: "John", Age: 25}
	result := db.WithQuery(user)
	suite.NotNil(result)

	// Test with empty struct query
	emptyUser := &TestUser{}
	result = db.WithQuery(emptyUser)
	suite.NotNil(result)

	// Test with fuzzy matching
	result = db.WithQuery(user, true)
	suite.NotNil(result)
}

// TestWithQueryRaw tests the WithQueryRaw method
func (suite *DatabaseTestSuite) TestWithQueryRaw() {
	db := suite.userDB

	// Test with raw query
	result := db.WithQueryRaw("name = ?", "John")
	suite.NotNil(result)

	// Test with multiple parameters
	result = db.WithQueryRaw("name = ? AND age > ?", "John", 18)
	suite.NotNil(result)

	// Test with empty query
	result = db.WithQueryRaw("")
	suite.NotNil(result)
}

// TestWithCursor tests the WithCursor method
func (suite *DatabaseTestSuite) TestWithCursor() {
	db := suite.userDB

	// Test with default cursor field (id)
	result := db.WithCursor("123", true)
	suite.NotNil(result)

	// Test with custom cursor field
	result = db.WithCursor("123", true, "created_at")
	suite.NotNil(result)

	// Test with previous page
	result = db.WithCursor("123", false)
	suite.NotNil(result)

	// Test with empty cursor value
	result = db.WithCursor("", true)
	suite.NotNil(result)
}

// TestWithTimeRange tests the WithTimeRange method
func (suite *DatabaseTestSuite) TestWithTimeRange() {
	db := suite.userDB

	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now

	// Test with both start and end time
	result := db.WithTimeRange("created_at", start, end)
	suite.NotNil(result)

	// Test with only start time (using zero end time)
	result = db.WithTimeRange("created_at", start, time.Time{})
	suite.NotNil(result)

	// Test with only end time (using zero start time)
	result = db.WithTimeRange("created_at", time.Time{}, end)
	suite.NotNil(result)

	// Test with empty column name
	result = db.WithTimeRange("", start, end)
	suite.NotNil(result)
}

// TestWithSelect tests the WithSelect method
func (suite *DatabaseTestSuite) TestWithSelect() {
	db := suite.userDB

	// Test with single field
	result := db.WithSelect("name")
	suite.NotNil(result)

	// Test with multiple fields
	result = db.WithSelect("name", "email", "age")
	suite.NotNil(result)

	// Test with empty fields
	result = db.WithSelect()
	suite.NotNil(result)
}

// TestWithSelectRaw tests the WithSelectRaw method
func (suite *DatabaseTestSuite) TestWithSelectRaw() {
	db := suite.userDB

	// Test with raw select
	result := db.WithSelectRaw("COUNT(*) as count")
	suite.NotNil(result)

	// Test with complex raw select
	result = db.WithSelectRaw("name, COUNT(*) as count, AVG(age) as avg_age")
	suite.NotNil(result)

	// Test with empty raw select
	result = db.WithSelectRaw("")
	suite.NotNil(result)
}

// TestWithLimit tests the WithLimit method
func (suite *DatabaseTestSuite) TestWithLimit() {
	db := suite.userDB

	// Test with positive limit
	result := db.WithLimit(10)
	suite.NotNil(result)

	// Test with zero limit
	result = db.WithLimit(0)
	suite.NotNil(result)

	// Test with negative limit
	result = db.WithLimit(-1)
	suite.NotNil(result)
}

// TestWithExpand tests the WithExpand method
func (suite *DatabaseTestSuite) TestWithExpand() {
	db := suite.userDB

	// Test with single association
	result := db.WithExpand([]string{"Profile"})
	suite.NotNil(result)

	// Test with multiple associations
	result = db.WithExpand([]string{"Profile", "Orders"})
	suite.NotNil(result)

	// Test with empty associations
	result = db.WithExpand([]string{})
	suite.NotNil(result)
}

// TestWithExclude tests the WithExclude method
func (suite *DatabaseTestSuite) TestWithExclude() {
	db := suite.userDB

	// Test with single field exclusion
	result := db.WithExclude(map[string][]any{"password": {}})
	suite.NotNil(result)

	// Test with multiple field exclusions
	result = db.WithExclude(map[string][]any{"password": {}, "secret": {}})
	suite.NotNil(result)

	// Test with empty exclusions
	result = db.WithExclude(map[string][]any{})
	suite.NotNil(result)
}

// TestWithLock tests the WithLock method
func (suite *DatabaseTestSuite) TestWithLock() {
	db := suite.userDB

	// Test with default lock (FOR UPDATE)
	result := db.WithLock()
	suite.NotNil(result)

	// Test with specific lock types using constants
	result = db.WithLock(consts.LockShare)
	suite.NotNil(result)

	result = db.WithLock(consts.LockUpdateNoWait)
	suite.NotNil(result)

	result = db.WithLock(consts.LockShareNoWait)
	suite.NotNil(result)

	result = db.WithLock(consts.LockUpdateSkipLocked)
	suite.NotNil(result)

	result = db.WithLock(consts.LockShareSkipLocked)
	suite.NotNil(result)
}

// TestWithJoinRaw tests the WithJoinRaw method
func (suite *DatabaseTestSuite) TestWithJoinRaw() {
	db := suite.userDB

	// Test with INNER JOIN
	result := db.WithJoinRaw("INNER JOIN profiles ON users.id = profiles.user_id")
	suite.NotNil(result)

	// Test with LEFT JOIN and parameters
	result = db.WithJoinRaw("LEFT JOIN orders ON users.id = orders.user_id AND orders.status = ?", "active")
	suite.NotNil(result)

	// Test with empty join
	result = db.WithJoinRaw("")
	suite.NotNil(result)
}

// Note: WithGroup and WithHaving methods are not yet implemented in the interface
// These tests are commented out until the methods are available

// TestWithOrder tests the WithOrder method
func (suite *DatabaseTestSuite) TestWithOrder() {
	db := suite.userDB

	// Test with single field ascending
	result := db.WithOrder("name")
	suite.NotNil(result)

	// Test with single field descending
	result = db.WithOrder("name DESC")
	suite.NotNil(result)

	// Test with multiple fields (comma separated)
	result = db.WithOrder("name ASC, age DESC")
	suite.NotNil(result)

	// Test with empty order
	result = db.WithOrder("")
	suite.NotNil(result)
}

// TestWithPagination tests the WithPagination method
func (suite *DatabaseTestSuite) TestWithPagination() {
	db := suite.userDB

	// Test with page and size parameters
	result := db.WithPagination(1, 10)
	suite.NotNil(result)

	// Test with different page and size
	result = db.WithPagination(2, 20)
	suite.NotNil(result)

	// Test with zero values
	result = db.WithPagination(0, 0)
	suite.NotNil(result)
}

// TestWithPurge tests the WithPurge method
func (suite *DatabaseTestSuite) TestWithPurge() {
	db := suite.userDB

	// Test enabling purge mode
	result := db.WithPurge()
	suite.NotNil(result)

	// Test with explicit true
	result = db.WithPurge(true)
	suite.NotNil(result)

	// Test with false
	result = db.WithPurge(false)
	suite.NotNil(result)
}

// TestWithCache tests the WithCache method
func (suite *DatabaseTestSuite) TestWithCache() {
	db := suite.userDB

	// Test enabling cache
	result := db.WithCache()
	suite.NotNil(result)

	// Test with explicit true
	result = db.WithCache(true)
	suite.NotNil(result)

	// Test with false
	result = db.WithCache(false)
	suite.NotNil(result)
}

// TestWithOmit tests the WithOmit method
func (suite *DatabaseTestSuite) TestWithOmit() {
	db := suite.userDB

	// Test with single field
	result := db.WithOmit("created_at")
	suite.NotNil(result)

	// Test with multiple fields
	result = db.WithOmit("created_at", "updated_at")
	suite.NotNil(result)

	// Test with empty fields
	result = db.WithOmit()
	suite.NotNil(result)
}

// TestWithTryRun tests the WithTryRun method
func (suite *DatabaseTestSuite) TestWithTryRun() {
	db := suite.userDB

	// Test enabling try run mode
	result := db.WithTryRun()
	suite.NotNil(result)

	// Test with explicit true
	result = db.WithTryRun(true)
	suite.NotNil(result)

	// Test with false
	result = db.WithTryRun(false)
	suite.NotNil(result)
}

// TestWithoutHook tests the WithoutHook method
func (suite *DatabaseTestSuite) TestWithoutHook() {
	db := suite.userDB

	// Test disabling hooks
	result := db.WithoutHook()
	suite.NotNil(result)

	// Test without hook
	result = db.WithoutHook()
	suite.NotNil(result)
}

// TestCreate tests the Create method
func (suite *DatabaseTestSuite) TestCreate() {
	db := suite.userDB

	// Test creating single user
	user := &TestUser{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      25,
		IsActive: true,
	}
	err := db.Create(user)
	suite.NoError(err)
	suite.NotEmpty(user.ID)
	suite.NotZero(user.CreatedAt)

	// Test creating multiple users
	users := []*TestUser{
		{Name: "Alice", Email: "alice@example.com", Age: 30, IsActive: true},
		{Name: "Bob", Email: "bob@example.com", Age: 28, IsActive: false},
	}
	err = db.Create(users...)
	suite.NoError(err)
	for _, u := range users {
		suite.NotEmpty(u.ID)
		suite.NotZero(u.CreatedAt)
	}

	// Test creating with batch size
	batchUsers := make([]*TestUser, 5)
	for i := range 5 {
		batchUsers[i] = &TestUser{
			Name:  fmt.Sprintf("User%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   20 + i,
		}
	}
	err = db.WithBatchSize(2).Create(batchUsers...)
	suite.NoError(err)
}

// TestDelete tests the Delete method
func (suite *DatabaseTestSuite) TestDelete() {
	db := suite.userDB

	// Create test data
	user := &TestUser{Name: "ToDelete", Email: "delete@example.com", Age: 25}
	err := db.Create(user)
	suite.NoError(err)

	// Test soft delete
	err = db.Delete(user)
	suite.NoError(err)

	// Verify soft delete
	var deletedUser TestUser
	err = db.Get(&deletedUser, user.ID)
	suite.Error(err) // Should not find soft deleted record

	// Test hard delete with purge
	user2 := &TestUser{Name: "ToHardDelete", Email: "harddelete@example.com", Age: 30}
	err = db.Create(user2)
	suite.NoError(err)

	err = db.WithPurge().Delete(user2)
	suite.NoError(err)

	// Test delete by condition
	users := []*TestUser{
		{Name: "DeleteMe1", Email: "del1@example.com", Age: 25},
		{Name: "DeleteMe2", Email: "del2@example.com", Age: 25},
	}
	err = db.Create(users...)
	suite.NoError(err)

	err = db.WithQueryRaw("age = ?", 25).Delete(&TestUser{})
	suite.NoError(err)
}

// TestUpdate tests the Update method
func (suite *DatabaseTestSuite) TestUpdate() {
	db := suite.userDB

	// Create test data
	user := &TestUser{Name: "Original", Email: "original@example.com", Age: 25}
	err := db.Create(user)
	suite.NoError(err)
	originalUpdatedAt := user.UpdatedAt

	// Test update
	time.Sleep(time.Millisecond) // Ensure different timestamp
	user.Name = "Updated"
	user.Age = 30
	err = db.Update(user)
	suite.NoError(err)
	suite.Equal("Updated", user.Name)
	suite.Equal(30, user.Age)
	suite.True(user.UpdatedAt.After(*originalUpdatedAt))

	// Test batch update
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 20},
		{Name: "User2", Email: "user2@example.com", Age: 22},
	}
	err = db.Create(users...)
	suite.NoError(err)

	for i, u := range users {
		u.Age = 25 + i
	}
	err = db.Update(users...)
	suite.NoError(err)
}

// TestUpdateById tests the UpdateById method
func (suite *DatabaseTestSuite) TestUpdateById() {
	db := suite.userDB

	// Create test data
	user := &TestUser{Name: "Original", Email: "original@example.com", Age: 25}
	err := db.Create(user)
	suite.NoError(err)

	// Test update by ID - update name field
	err = db.UpdateByID(user.ID, "name", "UpdatedById")
	suite.NoError(err)

	// Test update by ID - update age field
	err = db.UpdateByID(user.ID, "age", 35)
	suite.NoError(err)

	// Verify the updates
	var updatedUser TestUser
	err = db.Get(&updatedUser, user.ID)
	suite.NoError(err)
	suite.Equal("UpdatedById", updatedUser.Name)
	suite.Equal(35, updatedUser.Age)
}

// TestList tests the List method
func (suite *DatabaseTestSuite) TestList() {
	db := suite.userDB

	// Create test data
	users := []*TestUser{
		{Name: "User1", Email: "user1@example.com", Age: 20, IsActive: true},
		{Name: "User2", Email: "user2@example.com", Age: 25, IsActive: false},
		{Name: "User3", Email: "user3@example.com", Age: 30, IsActive: true},
	}
	err := db.Create(users...)
	suite.NoError(err)

	// Test list all
	var result []*TestUser
	err = db.List(&result)
	suite.NoError(err)
	suite.GreaterOrEqual(len(result), 3)

	// Test list with limit
	var limitedResult []*TestUser
	err = db.WithLimit(2).List(&limitedResult)
	suite.NoError(err)
	suite.Equal(2, len(limitedResult))

	// Test list with condition
	var activeUsers []*TestUser
	queryUser := &TestUser{IsActive: true}
	err = db.WithQuery(queryUser).List(&activeUsers)
	suite.NoError(err)
	suite.GreaterOrEqual(len(activeUsers), 2)

	// Test list with select
	var selectedUsers []*TestUser
	err = db.WithSelect("id", "name").List(&selectedUsers)
	suite.NoError(err)
	suite.GreaterOrEqual(len(selectedUsers), 3)
	for _, u := range selectedUsers {
		suite.NotEmpty(u.ID)
		suite.NotEmpty(u.Name)
		suite.Empty(u.Email) // Should be empty due to select
	}
}

// TestGet tests the Get method
func (suite *DatabaseTestSuite) TestGet() {
	db := suite.userDB

	// Create test data
	user := &TestUser{Name: "GetTest", Email: "get@example.com", Age: 25}
	err := db.Create(user)
	suite.NoError(err)

	// Test get by ID
	var result TestUser
	err = db.Get(&result, user.ID)
	suite.NoError(err)
	suite.Equal(user.Name, result.Name)
	suite.Equal(user.Email, result.Email)
	suite.Equal(user.Age, result.Age)

	// Test get non-existent record
	var notFound TestUser
	err = db.Get(&notFound, "non-existent-id")
	suite.Error(err)

	// Test get with select
	var selectedResult TestUser
	err = db.WithSelect("id", "name").Get(&selectedResult, user.ID)
	suite.NoError(err)
	suite.Equal(user.Name, selectedResult.Name)
	suite.Empty(selectedResult.Email) // Should be empty due to select
}

// TestCount tests the Count method
func (suite *DatabaseTestSuite) TestCount() {
	db := suite.userDB

	// Create test data
	users := []*TestUser{
		{Name: "Count1", Email: "count1@example.com", Age: 20, IsActive: true},
		{Name: "Count2", Email: "count2@example.com", Age: 25, IsActive: false},
		{Name: "Count3", Email: "count3@example.com", Age: 30, IsActive: true},
	}
	err := db.Create(users...)
	suite.NoError(err)

	// Test count all
	var count int64
	err = db.Count(&count)
	suite.NoError(err)
	suite.GreaterOrEqual(count, int64(3))

	// Test count with query
	activeQuery := &TestUser{IsActive: true}
	var activeCount int64
	err = db.WithQuery(activeQuery).Count(&activeCount)
	suite.NoError(err)
	suite.GreaterOrEqual(activeCount, int64(2))

	// Test count with age range
	var ageCount int64
	err = db.WithQueryRaw("age BETWEEN ? AND ?", 20, 30).Count(&ageCount)
	suite.NoError(err)
	suite.GreaterOrEqual(ageCount, int64(3))
}

// TestFirst tests the First method
func (suite *DatabaseTestSuite) TestFirst() {
	db := suite.userDB

	// Create test data
	users := []*TestUser{
		{Name: "First1", Email: "first1@example.com", Age: 20},
		{Name: "First2", Email: "first2@example.com", Age: 25},
	}
	err := db.Create(users...)
	suite.NoError(err)

	// Test first without condition
	var result TestUser
	err = db.First(&result)
	suite.NoError(err)
	suite.NotEmpty(result.ID)

	// Test first with condition
	var firstActive TestUser
	ageQuery := &TestUser{Age: 25}
	err = db.WithQuery(ageQuery).First(&firstActive)
	suite.NoError(err)
	suite.Equal(25, firstActive.Age)

	// Test first with order
	var orderedFirst TestUser
	err = db.WithOrder("age DESC").First(&orderedFirst)
	suite.NoError(err)
	suite.NotEmpty(orderedFirst.ID)
}

// TestLast tests the Last method
func (suite *DatabaseTestSuite) TestLast() {
	db := suite.userDB

	// Create test data
	users := []*TestUser{
		{Name: "Last1", Email: "last1@example.com", Age: 20},
		{Name: "Last2", Email: "last2@example.com", Age: 25},
	}
	err := db.Create(users...)
	suite.NoError(err)

	// Test last without condition
	var result TestUser
	err = db.Last(&result)
	suite.NoError(err)
	suite.NotEmpty(result.ID)

	// Test last with condition
	var lastAge TestUser
	ageQuery := &TestUser{Age: 20}
	err = db.WithQuery(ageQuery).Last(&lastAge)
	suite.NoError(err)
	suite.Equal(20, lastAge.Age)

	// Test last with order
	var orderedLast TestUser
	err = db.WithOrder("age ASC").Last(&orderedLast)
	suite.NoError(err)
	suite.NotEmpty(orderedLast.ID)
}

// TestTake tests the Take method
func (suite *DatabaseTestSuite) TestTake() {
	db := suite.userDB

	// Create test data
	user := &TestUser{Name: "TakeTest", Email: "take@example.com", Age: 25}
	err := db.Create(user)
	suite.NoError(err)

	// Test take without condition
	var result TestUser
	err = db.Take(&result)
	suite.NoError(err)
	suite.NotEmpty(result.ID)

	// Test take with condition
	var takeResult TestUser
	emailQuery := &TestUser{Email: "take@example.com"}
	err = db.WithQuery(emailQuery).Take(&takeResult)
	suite.NoError(err)
	suite.Equal("TakeTest", takeResult.Name)
	suite.Equal("take@example.com", takeResult.Email)

	// Test take non-existent
	var notFound TestUser
	notFoundQuery := &TestUser{Email: "nonexistent@example.com"}
	err = db.WithQuery(notFoundQuery).Take(&notFound)
	suite.Error(err)
}

// TestCleanup tests the Cleanup method
func (suite *DatabaseTestSuite) TestCleanup() {
	db := suite.userDB

	// Create and soft delete test data
	user := &TestUser{Name: "ToCleanup", Email: "cleanup@example.com", Age: 25}
	err := db.Create(user)
	suite.NoError(err)

	// Soft delete the user
	err = db.Delete(user)
	suite.NoError(err)

	// Verify soft delete
	var deletedUser TestUser
	err = db.Get(&deletedUser, user.ID)
	suite.Error(err) // Should not find soft deleted record

	// Test cleanup (permanent delete)
	err = db.Cleanup()
	suite.NoError(err)

	// After cleanup, the record should be permanently removed
	// We can't easily verify this without accessing the database directly
	// but the method should not return an error
}

// TestHealth tests the Health method
func (suite *DatabaseTestSuite) TestHealth() {
	db := suite.userDB

	// Test health check
	err := db.Health()
	suite.NoError(err)

	// Test health with context timeout
	// Note: Health method doesn't support context timeout in current implementation
	// This is a basic health check test
}

// Benchmark tests

// BenchmarkCreate benchmarks the Create method
func BenchmarkCreate(b *testing.B) {
	db := database.Database[*TestUser](nil)

	for i := 0; b.Loop(); i++ {
		user := &TestUser{
			Name:     fmt.Sprintf("BenchUser%d", i),
			Email:    fmt.Sprintf("bench%d@example.com", i),
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
		}
		err := db.Create(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCreateBatch benchmarks batch creation
func BenchmarkCreateBatch(b *testing.B) {
	db := database.Database[*TestUser](nil)
	batchSize := 100

	for i := 0; b.Loop(); i++ {
		users := make([]*TestUser, batchSize)
		for j := range batchSize {
			users[j] = &TestUser{
				Name:     fmt.Sprintf("BatchUser%d_%d", i, j),
				Email:    fmt.Sprintf("batch%d_%d@example.com", i, j),
				Age:      20 + (j % 50),
				IsActive: j%2 == 0,
			}
		}
		err := db.WithBatchSize(50).Create(users...)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGet benchmarks the Get method
func BenchmarkGet(b *testing.B) {
	db := database.Database[*TestUser](nil)

	// Create test data
	user := &TestUser{Name: "BenchGet", Email: "benchget@example.com", Age: 25}
	err := db.Create(user)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		var result TestUser
		err := db.Get(&result, user.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkList benchmarks the List method
func BenchmarkList(b *testing.B) {
	db := database.Database[*TestUser](nil)

	// Create test data
	users := make([]*TestUser, 100)
	for i := range 100 {
		users[i] = &TestUser{
			Name:     fmt.Sprintf("ListUser%d", i),
			Email:    fmt.Sprintf("list%d@example.com", i),
			Age:      20 + (i % 50),
			IsActive: i%2 == 0,
		}
	}
	err := db.Create(users...)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		var result []*TestUser
		err := db.WithLimit(10).List(&result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpdate benchmarks the Update method
func BenchmarkUpdate(b *testing.B) {
	db := database.Database[*TestUser](nil)

	// Create test data
	user := &TestUser{Name: "BenchUpdate", Email: "benchupdate@example.com", Age: 25}
	err := db.Create(user)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; b.Loop(); i++ {
		user.Age = 25 + (i % 50)
		err := db.Update(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCount benchmarks the Count method
func BenchmarkCount(b *testing.B) {
	db := database.Database[*TestUser](nil)

	// Create test data
	users := make([]*TestUser, 50)
	for i := range 50 {
		users[i] = &TestUser{
			Name:     fmt.Sprintf("CountUser%d", i),
			Email:    fmt.Sprintf("count%d@example.com", i),
			Age:      20 + (i % 30),
			IsActive: i%2 == 0,
		}
	}
	err := db.Create(users...)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		var count int64
		err := db.Count(&count)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFirst benchmarks the First method
func BenchmarkFirst(b *testing.B) {
	db := database.Database[*TestUser](nil)

	// Create test data
	users := make([]*TestUser, 20)
	for i := range 20 {
		users[i] = &TestUser{
			Name:     fmt.Sprintf("FirstUser%d", i),
			Email:    fmt.Sprintf("first%d@example.com", i),
			Age:      20 + i,
			IsActive: i%2 == 0,
		}
	}
	err := db.Create(users...)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		var result TestUser
		err := db.WithOrder("age ASC").First(&result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkQueryBuilder benchmarks query building methods
func BenchmarkQueryBuilder(b *testing.B) {
	db := database.Database[*TestUser](nil)

	for b.Loop() {
		queryUser := &TestUser{Age: 25, IsActive: true}
		_ = db.WithQuery(queryUser).
			WithOrder("created_at DESC").
			WithLimit(10).
			WithSelect("id", "name", "email")
	}
}

package types

import (
	"io"
	"time"

	"github.com/forbearing/golib/types/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Initalizer interface is used to initialize configuration, flag arguments, logger, or other components.
// This interface is commonly implemented by bootstrap components that need to perform
// initialization tasks during application startup.
//
// Example implementations:
//   - Configuration loaders
//   - Logger initializers
//   - Database connection setup
//   - Cache initialization
type Initalizer interface {
	Init() error
}

// StandardLogger interface provides standard logging methods for custom logger implementations.
// This interface follows the traditional logging pattern with both simple and formatted logging methods.
//
// Usage:
//   - Implement this interface to create custom loggers
//   - Use Debug/Info/Warn/Error for simple logging
//   - Use Debugf/Infof/Warnf/Errorf for formatted logging
//   - Fatal methods should terminate the program after logging
type StandardLogger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)

	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

// StructuredLogger interface provides structured logging methods with key-value pairs.
// This interface is designed for structured logging where additional context can be
// attached to log messages as key-value pairs.
//
// Usage:
//
//	logger.Infow("User login", "userID", 123, "ip", "192.168.1.1")
//	logger.Errorw("Database error", "error", err, "query", sql)
//
// The 'w' suffix stands for "with" (structured data).
type StructuredLogger interface {
	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Fatalw(msg string, keysAndValues ...any)
}

// ZapLogger interface provides zap-specific logging methods with structured fields.
// This interface is designed for integration with the uber-go/zap logging library,
// offering high-performance structured logging capabilities.
//
// Usage:
//
//	logger.Infoz("Request processed", zap.String("method", "GET"), zap.Int("status", 200))
//	logger.Errorz("Database connection failed", zap.Error(err), zap.String("host", dbHost))
//
// The 'z' suffix distinguishes these methods from other logging interfaces.
type ZapLogger interface {
	Debugz(msg string, fields ...zap.Field)
	Infoz(msg string, fields ...zap.Field)
	Warnz(msg string, fields ...zap.Field)
	Errorz(msg string, fields ...zap.Field)
	Fatalz(msg string, fields ...zap.Field)
}

// Logger interface combines all logging capabilities into a unified interface.
// This interface provides comprehensive logging functionality by embedding
// StandardLogger, StructuredLogger, and ZapLogger interfaces, along with
// context-aware logging methods.
//
// Key features:
//   - Standard logging (Debug, Info, Warn, Error, Fatal)
//   - Structured logging with key-value pairs (Debugw, Infow, etc.)
//   - Zap-specific structured logging with typed fields
//   - Context-aware logging for controllers, services, and database operations
//   - Support for complex object and array marshaling
//
// This unified approach allows flexible logging usage throughout the application.
type Logger interface {
	With(fields ...string) Logger

	WithObject(name string, obj zapcore.ObjectMarshaler) Logger
	WithArray(name string, arr zapcore.ArrayMarshaler) Logger

	WithControllerContext(*ControllerContext, consts.Phase) Logger
	WithServiceContext(*ServiceContext, consts.Phase) Logger
	WithDatabaseContext(*DatabaseContext, consts.Phase) Logger

	StandardLogger
	StructuredLogger
	ZapLogger
}

// Database interface provides comprehensive database operations for any model type.
// This interface is constrained by the Model interface, ensuring type safety
// for database operations across different model implementations.
//
// Generic constraint:
//
//	M must implement the Model interface (typically by embedding model.Base)
//
// Core operations:
//   - CRUD operations with automatic timestamp management
//   - Flexible querying with various finder methods
//   - Health monitoring and cleanup capabilities
//   - Optional caching support for improved performance
//
// The interface embeds DatabaseOption[M] to provide chainable query building.
type Database[M Model] interface {
	// Create one or multiple record.
	// Pass M to create one record,
	// Pass []M to create multiple record.
	// It will update the "created_at" and "updated_at" field.
	Create(objs ...M) error
	// Delete one or multiple record.
	// Pass M to delete one record.
	// Pass []M to delete multiple record.
	Delete(objs ...M) error
	// Update one or multiple record, if record doesn't exist, it will be created.
	// Pass M to update one record.
	// Pass []M to update multiple record.
	// It will just update the "updated_at" field.
	Update(objs ...M) error
	// UpdateById only update one record with specific id.
	// its not invoke model hook.
	UpdateById(id string, key string, value any) error
	// List all records and write to dest.
	List(dest *[]M, cache ...*[]byte) error
	// Get one record with specific id and write to dest.
	Get(dest M, id string, cache ...*[]byte) error
	// First finds the first record ordered by primary key.
	First(dest M, cache ...*[]byte) error
	// Last finds the last record ordered by primary key
	Last(dest M, cache ...*[]byte) error
	// Take finds the first record returned by the database in no specified order.
	Take(dest M, cache ...*[]byte) error
	// Count returns the total number of records with the given query condition.
	Count(*int64) error
	// Cleanup delete all records that column 'deleted_at' is not null.
	Cleanup() error
	// Health checks the database connectivity and basic operations.
	// It returns nil if the database is healthy, otherwise returns an error.
	Health() error

	DatabaseOption[M]
}

// DatabaseOption interface.
// WithXXX setting database options.
type DatabaseOption[M Model] interface {
	// WithDB returns a new database manipulator, only support *gorm.DB.
	WithDB(any) Database[M]

	// WithTable multiple custom table, always used with the method `WithDB`.
	WithTable(name string) Database[M]

	// WithDebug setting debug mode, the priority is higher than config.Server.LogLevel and default value(false).
	WithDebug() Database[M]

	// WithQuery is where condition.
	WithQuery(query M, fuzzyMatch ...bool) Database[M]

	// WithQueryRaw is where condition.
	// database.WithQueryRaw(xxx) same as database.WithQuery(xxx) and provides more flexible query.
	// Examples:
	// - WithQueryRaw("name = ?", "hybfkuf")
	// - WithQueryRaw("name <> ?", "hybfkuf")
	// - WithQueryRaw("name IN (?)", []string{"hybfkuf", "hybfkuf 2"})
	// - WithQueryRaw("name LIKE ?", "%hybfkuf%")
	// - WithQueryRaw("name = ? AND age >= ?", "hybfkuf", "100")
	// - WithQueryRaw("updated_at > ?", lastWeek)
	// - WithQueryRaw("created_at BETWEEN ? AND ?", lastWeek, today)
	WithQueryRaw(query any, args ...any) Database[M]

	// WithCursor enables cursor-based pagination.
	// cursorValue is the value of the last record in the previous page.
	// next indicates the direction of pagination:
	//   - true: fetch records after the cursor (next page)
	//   - false: fetch records before the cursor (previous page)
	//
	// Example:
	//
	//	// First page (no cursor)
	//	database.Database[*model.User]().WithLimit(10).List(&users)
	//	// Next page (using last user's ID as cursor)
	//	lastID := users[len(users)-1].ID
	//	database.Database[*model.User]().WithCursor(lastID, true).WithLimit(10).List(&nextUsers)
	//	// Next page (using last user id as cursor)
	//	database.Database[*model.User]().WithCursor(lastID, true, "user_id").WithLimit(10).List(&nextUsers)
	WithCursor(string, bool, ...string) Database[M]

	// WithAnd with AND query condition(default).
	// It must be called before WithQuery.
	WithAnd(...bool) Database[M]

	// WithAnd with OR query condition.
	// It must be called before WithQuery.
	WithOr(...bool) Database[M]

	// WithTimeRange applies a time range filter to the query based on the specified column name.
	// It restricts the results to records where the column's value falls within the specified start and end times.
	// This method is designed to be used in a chainable manner, allowing for the construction of complex queries.
	//
	// Parameters:
	// - columnName: The name of the column to apply the time range filter on. This should be a valid date/time column in the database.
	// - startTime: The beginning of the time range. Records with the column's value equal to or later than this time will be included.
	// - endTime: The end of the time range. Records with the column's value equal to or earlier than this time will be included.
	//
	// Returns: A modified Database instance that includes the time range filter in its query conditions.
	WithTimeRange(columnName string, startTime time.Time, endTime time.Time) Database[M]

	// WithSelect specify fields that you want when querying, creating, updating
	// default select all fields.
	WithSelect(columns ...string) Database[M]

	// WithSelectRaw
	WithSelectRaw(query any, args ...any) Database[M]

	// WithIndex use specific index to query.
	WithIndex(index string) Database[M]

	// WithTransaction executes operations within a transaction.
	WithTransaction(tx any) Database[M]

	// WithJoinRaw
	WithJoinRaw(query string, args ...any) Database[M]

	// TODO:
	// WithGroup(name string) Database[M]
	// WithHaving(query any, args ...any) Database[M]

	// WithLock adds locking clause to SELECT statement.
	// It must be used within a transaction (WithTransaction).
	WithLock(mode ...string) Database[M]

	// WithBatchSize set batch size for bulk operations. affects Create, Update, Delete.
	WithBatchSize(size int) Database[M]

	// WithScope applies pagination parameters to the query, useful for retrieving data in pages.
	// This method enables front-end applications to request a specific subset of records,
	// based on the desired page number and the number of records per page.
	//
	// Parameters:
	// - page: The page number being requested. Page numbers typically start at 1.
	// - size: The number of records to return per page. This determines the "size" of each page.
	//
	// The pagination logic calculates the offset based on the page number and size,
	// and applies it along with the limit (size) to the query. This facilitates efficient
	// data fetching suitable for front-end pagination displays.
	//
	// Returns: A modified Database instance that includes pagination parameters in its query conditions.
	WithScope(page, size int) Database[M]

	// WithLimit determines how much record should retrieve.
	// limit is 0 or -1 means no limit.
	WithLimit(limit int) Database[M]

	// WithExclude excludes records that matchs a condition within a list.
	// For example:
	//   - If you want exlcude users with specific ids from your query,
	//     you can use WithExclude(excludes),
	//     excludes: "id" as key, ["myid1", "myid2", "myid3"] as value.
	//   - If you want excludes users that id not ["myid1", "myid2"] and not not ["root", "noname"],
	//     the `excludes` should be:
	//     excludes := make(map[string][]any)
	//     excludes["id"] = []any{"myid1", "myid2"}
	//     excludes["name"] = []any{"root", "noname"}.
	WithExclude(map[string][]any) Database[M]

	// WithOrder
	// For example:
	// - WithOrder("name") // default ASC.
	// - WithOrder("name desc")
	// - WithOrder("created_at")
	// - WithOrder("updated_at desc")
	// NOTE: you cannot using the mysql keyword, such as: "order", "limit".
	WithOrder(order string) Database[M]

	// WithExpand, for "foreign key".
	WithExpand(expand []string, order ...string) Database[M]

	// WithPurge tells the database manipulator to delete resource in database permanently.
	WithPurge(...bool) Database[M]
	// WithCache tells the database manipulator to retrieve resource from cache.
	WithCache(...bool) Database[M]
	// WithOmit omit specific columns when create/update.
	WithOmit(...string) Database[M]
	// WithTryRun only executes model hooks without performing actual database operations.
	// Also logs the SQL statements that would have been executed.
	WithTryRun(...bool) Database[M]
	// WithoutHook tells the database manipulator not invoke model hooks.
	WithoutHook() Database[M]
}

// Model interface defines the contract for all data models in the framework.
// This interface ensures consistent behavior across different model implementations
// and provides comprehensive functionality for database operations, logging, and lifecycle hooks.
//
// Implementation requirements:
//  1. Must be a pointer to a struct (e.g., *User, not User) - otherwise causes panic
//  2. Must have an "ID" field as the primary key in the database
//  3. Should embed model.Base to inherit common fields and methods
//
// Example implementation:
//
//	type User struct {
//	    model.Base
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//
//	func (u *User) GetTableName() string {
//	    return "users"
//	}
//
// Core functionality:
//   - Table and ID management for database operations
//   - Audit trail with created/updated timestamps and user tracking
//   - Relationship management through Expands() for foreign key preloading
//   - Query filtering through Excludes() for conditional operations
//   - Structured logging support via zap.ObjectMarshaler
//   - Lifecycle hooks for custom business logic during CRUD operations
type Model interface {
	GetTableName() string // GetTableName returns the table name.
	GetID() string
	SetID(id ...string) // SetID method will automatically set the id if id is empty.
	GetCreatedBy() string
	GetUpdatedBy() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreatedBy(string)
	SetUpdatedBy(string)
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
	Expands() []string // Expands returns the foreign keys should preload.
	Excludes() map[string][]any
	MarshalLogObject(zapcore.ObjectEncoder) error // MarshalLogObject implement zap.ObjectMarshaler

	Hooker
}

type (
	Request  interface{}
	Response interface{}
)

// Service interface provides comprehensive business logic operations for model types.
// This interface defines the service layer that sits between controllers and database operations,
// implementing business rules, validation, complex operations, and lifecycle management.
//
// Implementation requirements:
//   - The implementing object must be a pointer to struct
//
// Generic constraints:
//   - M: Must implement the Model interface
//   - REQ: Request type (typically DTOs or request structures)
//   - RSP: Response type (typically DTOs or response structures)
//
// Core operations:
//   - CRUD operations: Create, Delete, Update, Patch, List, Get
//   - Batch operations: CreateMany, DeleteMany, UpdateMany, PatchMany
//   - Lifecycle hooks: Before/After methods for each operation
//   - Data operations: Import/Export for bulk data management
//   - Filtering: Custom filtering logic for queries
//
// ServiceContext provides:
//   - HTTP request/response context
//   - Database transaction management
//   - User authentication and authorization context
//   - Request validation and data binding
//   - Logging and tracing capabilities
//
// Hook methods allow custom business logic:
//   - Before hooks: Validation, authorization, data transformation
//   - After hooks: Notifications, caching, audit logging, cleanup
type Service[M Model, REQ Request, RSP Response] interface {
	Create(*ServiceContext, REQ) (RSP, error)
	Delete(*ServiceContext, REQ) (RSP, error)
	Update(*ServiceContext, REQ) (RSP, error)
	Patch(*ServiceContext, REQ) (RSP, error)
	List(*ServiceContext, REQ) (RSP, error)
	Get(*ServiceContext, REQ) (RSP, error)

	CreateMany(*ServiceContext, REQ) (RSP, error)
	DeleteMany(*ServiceContext, REQ) (RSP, error)
	UpdateMany(*ServiceContext, REQ) (RSP, error)
	PatchMany(*ServiceContext, REQ) (RSP, error)

	CreateBefore(*ServiceContext, M) error
	CreateAfter(*ServiceContext, M) error
	DeleteBefore(*ServiceContext, M) error
	DeleteAfter(*ServiceContext, M) error
	UpdateBefore(*ServiceContext, M) error
	UpdateAfter(*ServiceContext, M) error
	PatchBefore(*ServiceContext, M) error
	PatchAfter(*ServiceContext, M) error
	ListBefore(*ServiceContext, *[]M) error
	ListAfter(*ServiceContext, *[]M) error
	GetBefore(*ServiceContext, M) error
	GetAfter(*ServiceContext, M) error

	CreateManyBefore(*ServiceContext, ...M) error
	CreateManyAfter(*ServiceContext, ...M) error
	DeleteManyBefore(*ServiceContext, ...M) error
	DeleteManyAfter(*ServiceContext, ...M) error
	UpdateManyBefore(*ServiceContext, ...M) error
	UpdateManyAfter(*ServiceContext, ...M) error
	PatchManyBefore(*ServiceContext, ...M) error
	PatchManyAfter(*ServiceContext, ...M) error

	Import(*ServiceContext, io.Reader) ([]M, error)
	Export(*ServiceContext, ...M) ([]byte, error)

	Filter(*ServiceContext, M) M
	FilterRaw(*ServiceContext) string

	Logger
}

// Hooker interface defines lifecycle hooks that can be executed at various points
// during database operations. This interface enables custom business logic,
// validation, auditing, and side effects to be executed automatically.
//
// Hook execution order:
//  1. Before hooks are called first (validation, authorization)
//  2. Main operation is performed (database CRUD)
//  3. After hooks are called last (notifications, caching, cleanup)
//
// Common use cases:
//   - CreateBefore: Validate data, set defaults, check permissions
//   - CreateAfter: Send notifications, update caches, log audit trail
//   - UpdateBefore: Validate changes, check business rules
//   - UpdateAfter: Invalidate caches, trigger workflows
//   - DeleteBefore: Check dependencies, backup data
//   - DeleteAfter: Clean up related data, update statistics
//
// Error handling:
//   - Before hooks can prevent the operation by returning an error
//   - After hooks should handle errors gracefully to avoid rollbacks
type Hooker interface {
	CreateBefore() error
	CreateAfter() error
	DeleteBefore() error
	DeleteAfter() error
	UpdateBefore() error
	UpdateAfter() error
	ListBefore() error
	ListAfter() error
	GetBefore() error
	GetAfter() error
}

// Cache interface provides generic caching operations for any data type.
// This interface supports multiple caching strategies and implementations,
// including in-memory caches, distributed caches, and hybrid approaches.
//
// Generic type T can be any serializable data type.
//
// Available implementations:
//   - freecache: High-performance in-memory cache with zero GC overhead
//   - gocache: Multi-tier caching with various backends
//   - smap: Simple concurrent map-based cache
//   - cmap: Concurrent map with advanced features
//   - ccache: LRU cache with size limits
//
// Operations:
//   - Set: Store value with TTL (time-to-live)
//   - Get: Retrieve value and mark as accessed (affects LRU)
//   - Peek: Retrieve value without affecting access order
//   - Delete: Remove specific key
//   - Exists: Check if key exists without retrieving value
//   - Len: Get current number of cached items
//   - Clear: Remove all cached items
//
// Usage patterns:
//   - Application-level caching for expensive computations
//   - Database query result caching
//   - Session storage
//   - Rate limiting counters
//   - Temporary data storage
type Cache[T any] interface {
	Set(key string, values T, ttl time.Duration)
	Get(key string) (T, bool)
	Peek(key string) (T, bool)
	Delete(key string)
	Exists(key string) bool
	Len() int
	Clear()

	// Increment(key string, delta int64) (int64, error)
	// Decrement(key string, delta int64) (int64, error)
}

// RBAC interface defines comprehensive role-based access control operations.
// This interface provides a complete RBAC system supporting roles, permissions,
// and subject assignments with flexible resource and action management.
//
// RBAC Model Components:
//   - Subject: Users or entities that need access (e.g., "user:123", "service:api")
//   - Role: Named collection of permissions (e.g., "admin", "editor", "viewer")
//   - Resource: Protected objects or endpoints (e.g., "users", "posts", "/api/v1/users")
//   - Action: Operations on resources (e.g., "read", "write", "delete", "create")
//
// Permission Model:
//   - Permissions are defined as (role, resource, action) tuples
//   - Subjects are assigned roles, inheriting all role permissions
//   - Supports hierarchical roles and resource patterns
//
// Implementation:
//   - Typically backed by Casbin for policy enforcement
//   - Supports both file-based and database-backed policy storage
//   - Can integrate with external identity providers
//
// Usage patterns:
//   - API endpoint authorization
//   - Resource-level access control
//   - Multi-tenant permission management
type RBAC interface {
	AddRole(name string) error
	RemoveRole(name string) error

	GrantPermission(role string, resource string, action string) error
	RevokePermission(role string, resource string, action string) error

	AssignRole(subject string, role string) error
	UnassignRole(subject string, role string) error
}

// ESDocumenter represents a document that can be indexed into Elasticsearch.
// Types implementing this interface should be able to convert themselves
// into a document format suitable for Elasticsearch indexing.
type ESDocumenter interface {
	// Document returns a map representing an Elasticsearch document.
	// The returned map should contain all fields to be indexed, where:
	//   - keys are field names (string type)
	//   - values are field values (any type)
	//
	// Implementation notes:
	//   1. The returned map should only contain JSON-serializable values.
	//   2. Field names should match those defined in the Elasticsearch mapping.
	//   3. Complex types (like nested objects or arrays) should be correctly
	//      represented in the returned map.
	//
	// Example:
	//   return map[string]any{
	//       "id":    "1234",
	//       "title": "Sample Document",
	//       "tags":  []string{"tag1", "tag2"},
	//   }
	Document() map[string]any

	// GetID returns a string that uniquely identifies the document.
	// This ID is typically used as the Elasticsearch document ID.
	//
	// Implementation notes:
	//   1. The ID should be unique within the index.
	//   2. If no custom ID is needed, consider returning an empty string
	//      to let Elasticsearch auto-generate an ID.
	//   3. The ID should be a string, even if it's originally a numeric value.
	//
	// Example:
	//   return "user_12345"
	GetID() string
}

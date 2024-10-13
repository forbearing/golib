package types

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Initalizer interface used to initial configuration, flag arguments or logger, etc.
type Initalizer interface {
	Init() error
}

// StandardLogger interface used for user to custom themselves logger.
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

// StructuredLogger is structured logger interface used for user to custom themselves logger.
type StructuredLogger interface {
	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Fatalw(msg string, keysAndValues ...any)
}
type ZapLogger interface {
	Debugz(msg string, fields ...zap.Field)
	Infoz(msg string, fields ...zap.Field)
	Warnz(msg string, feilds ...zap.Field)
	Errorz(msg string, fields ...zap.Field)
	Fatalz(msg string, fields ...zap.Field)
}

type Logger interface {
	With(key, value string) Logger

	StandardLogger
	StructuredLogger
	ZapLogger
}

// Database interface.
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
	UpdateById(id any, key string, value any) error
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
	// WithoutHook tells the database manipulator not invoke model hooks.
	WithoutHook() Database[M]
}

// Model interface.
// The two principles must be follwed before you implements Model interface.
//   - The object that implements the Model interface must be pointer to structure,
//     otherwise cause panic.
//   - The structure must have feild "ID" and the field must be primaryKey in database.
type Model interface {
	GetTableName() string // GetTableName returns the table name.
	GetID() string
	SetID(id ...string) // SetID method will automatically set the id if id is empty.
	GetCreatedBy() string
	GetUpdatedBy() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreatedBy(s string)
	SetUpdatedBy(s string)
	SetCreatedAt(t time.Time)
	SetUpdatedAt(t time.Time)
	Expands() []string // Expands returns the foreign keys should preload.
	Excludes() map[string][]any
	MarshalLogObject(zapcore.ObjectEncoder) error // MarshalLogObject implement zap.ObjectMarshaler

	Hooker
}

// Service interface.
// The object that implements this interface must be pointer to struct.
//
// xxxBefore
// - 权限校验
// xxxAfter
// - 操作行为记录
type Service[M Model] interface {
	CreateBefore(*ServiceContext, ...M) error
	CreateAfter(*ServiceContext, ...M) error
	DeleteBefore(*ServiceContext, ...M) error
	DeleteAfter(*ServiceContext, ...M) error
	UpdateBefore(*ServiceContext, ...M) error
	UpdateAfter(*ServiceContext, ...M) error
	UpdatePartialBefore(*ServiceContext, ...M) error
	UpdatePartialAfter(*ServiceContext, ...M) error
	ListBefore(*ServiceContext, *[]M) error // 必须是指针类型, 因为有时候需要修改原数据
	ListAfter(*ServiceContext, *[]M) error  // 必须是指针类型, 因为有时候需要修改原数据
	GetBefore(*ServiceContext, ...M) error
	GetAfter(*ServiceContext, ...M) error
	// Import file.
	Import(*ServiceContext, io.Reader) ([]M, error)
	// Export records from database and write into excel.
	Export(*ServiceContext, ...M) ([]byte, error)

	Filter(*ServiceContext, M) M
	FilterRaw(*ServiceContext) string

	Logger
}

// Hooker interface.
type Hooker interface {
	// CreateBefore is a hook that will be invoked before create records in database.
	// For example:
	// - Make sure the user mobile and email is valid before create in database.
	// - Set record default value.
	CreateBefore() error
	// CreateAfter is a hook that will be invoked after create in database.
	CreateAfter() error

	// DeleteBefore is a hook that will be invoked before delete records in database.
	DeleteBefore() error
	// DeleteAfter is a hook that will be invoked after delete in database.
	DeleteAfter() error

	// UpdateBefore is a hook that will be invoked before update records in database.
	UpdateBefore() error
	// UpdateAfter is a hook that will be invoked after update in database.
	UpdateAfter() error

	// UpdatePartialBefore is a hook that will be invoked before update records in database.
	UpdatePartialBefore() error
	// UpdatePartialAfter is a hook that will be invoked after update in database.
	UpdatePartialAfter() error

	// ListBefore is a hook that will be invoked before list records in database.
	ListBefore() error
	// ListAfter is a hook that will be invoked after list in database.
	// For examples:
	// - clean user password before responsed to frontend.
	ListAfter() error

	// GetBefore is a hook that will be invoked before get records in database.
	GetBefore() error
	// For examples:
	// - clean user password before responsed to frontend.
	// GetAfter is a hook that will be invoked after get in database.
	GetAfter() error
}

// Cache interface defines the standard operations for a generic cache mechanism.
type Cache[M Model] interface {
	// Set stores a Model identified by a key in the cache, with an expiration time.
	// If the cache already contains the key, the existing Model is overwritten.
	// The ttl (time to live) determines how long the key stays in the cache before it is auto-evicted.
	// Returns an error if the operation fails.
	Set(key string, values ...M)
	// SetInt stores an integer value identified by a key in the cache, with an expiration time.
	// If the cache already contains the key, the existing value is overwritten.
	// The ttl (time to live) determines how long the key stays in the cache before it is auto-evicted.
	// This method is specifically for storing integer values.
	SetInt(key string, values ...int64)
	// SetBool stores a boolean value identified by a key in the cache, with an expiration time.
	// If the cache already contains the key, the existing value is overwritten.
	// The ttl (time to live) determines how long the key stays in the cache before it is auto-evicted.
	SetBool(key string, values ...bool)
	// SetFloat stores a floating-point value identified by a key in the cache, with an expiration time.
	// If the cache already contains the key, the existing value is overwritten.
	// The ttl (time to live) determines how long the key stays in the cache before it is auto-evicted.
	SetFloat(key string, values ...float64)
	// SetString stores a string value identified by a key in the cache, with an expiration time.
	// If the cache already contains the key, the existing value is overwritten.
	// The ttl (time to live) determines how long the key stays in the cache before it is auto-evicted.
	SetString(key string, values ...string)
	// SetAny stores a value identified by a key in the cache, with an expiration time.
	// If the cache already contains the key, the existing value is overwritten.
	// The ttl (time to live) determines how long the key stays in the cache before it is auto-evicted.
	SetAny(key string, value any)

	// Get retrieves a Model based on its key from the cache.
	// Returns the Model associated with the key and a boolean indicating whether the key was found.
	// If the key does not exist, the returned error will be non-nil.
	Get(key string) ([]M, bool)
	// GetInt retrieves an integer value based on its key from the cache.
	// Returns the integer value associated with the key and a boolean indicating whether the key was found.
	// If the key does not exist or the value is not an integer, the returned error will be non-nil.
	GetInt(key string) ([]int64, bool)
	// GetBool retrieves a boolean value based on its key from the cache.
	// Returns the boolean value associated with the key and a boolean indicating whether the key was found.
	// If the key does not exist, the returned error will be non-nil.
	GetBool(key string) ([]bool, bool)
	// GetFloat retrieves a floating-point value based on its key from the cache.
	// Returns the floating-point value associated with the key and a boolean indicating whether the key was found.
	// If the key does not exist, the returned error will be non-nil.
	GetFloat(key string) ([]float64, bool)
	// GetString retrieves a string value based on its key from the cache.
	// Returns the string value associated with the key and a boolean indicating whether the key was found.
	// If the key does not exist, the returned error will be non-nil.
	GetString(key string) ([]string, bool)
	// GetAny retrieves a value based on its key from the cache.
	// Returns the value associated with the key and a boolean indicating whether the key was found.
	// If the key does not exist, the returned error will be non-nil.
	GetAny(key string) (any, bool)

	// GetAll returns a map of all key-value pairs currently in the cache.
	GetAll() map[string][]M
	// GetAllInt returns a map containing the key-value pairs where the values are integers in the cache.
	GetAllInt() map[string][]int64
	// GetAllBool returns a map containing the key-value pairs where the values are booleans in the cache.
	GetAllBool() map[string][]bool
	// GetAllFloat returns a map containing the key-value pairs where the values are floating-point numbers in the cache.
	GetAllFloat() map[string][]float64

	// GetAllString returns a map containing the key-value pairs where the values are strings in the cache.
	GetAllString() map[string][]string

	// GetAllAny returns a map containing all the key-value pairs present in the cache, regardless of the value types.
	GetAllAny() map[string]any

	// Peek retrieves the Model associated with the given key without updating the key's access time.
	// This is particularly useful in caches where the age of an item might determine its eviction from the cache.
	// Returns the Model associated with the key and a boolean indicating whether the key was found.
	Peek(key string) ([]M, bool)
	// PeekInt retrieves the integer value associated with the given key from the cache without updating the key's access time.
	// This is particularly useful in caches where the age of an item might determine its eviction from the cache.
	// Returns the integer value associated with the key and a boolean indicating whether the key was found.
	PeekInt(key string) ([]int64, bool)
	// PeekBool retrieves the boolean value associated with the given key from the cache without updating the key's access time.
	// This is particularly useful in caches where the age of an item might determine its eviction from the cache.
	// Returns the boolean value associated with the key and a boolean indicating whether the key was found.
	PeekBool(key string) ([]bool, bool)
	// PeekString retrieves the string value associated with the given key from the cache without updating the key's access time.
	// This is particularly useful in caches where the age of an item might determine its eviction from the cache.
	// Returns the string value associated with the key and a boolean indicating whether the key was found.
	PeekString(key string) ([]string, bool)
	// PeekFloat retrieves the floating-point value associated with the given key from the cache without updating the key's access time.
	// This is particularly useful in caches where the age of an item might determine its eviction from the cache.
	// Returns the floating-point value associated with the key and a boolean indicating whether the key was found.
	PeekFloat(key string) ([]float64, bool)
	// PeekAny retrieves the value associated with the given key from the cache without updating the key's access time.
	// This is particularly useful in caches where the age of an item might determine its eviction from the cache.
	// Returns the value associated with the key and a boolean indicating whether the key was found.
	PeekAny(key string) (any, bool)

	// Remove removes the Model associated with a key from the cache.
	Remove(key string)
	// RemoveInt removes the integer value associated with the given key from the cache.
	RemoveInt(key string)
	// RemoveBool removes the boolean value associated with the given key from the cache.
	RemoveBool(key string)
	// RemoveFloat removes the floating-point value associated with the given key from the cache.
	RemoveFloat(key string)
	// RemoveString removes the string value associated with the given key from the cache.
	RemoveString(key string)
	// RemoveAny removes the value associated with the given key from the cache.
	RemoveAny(key string)

	// Exists checks if the cache contains a Model associated with the given key.
	// Returns a boolean indicating whether the key exists in the cache.
	Exists(key string) bool
	// ExistsInt checks if the cache contains an integer value associated with the given key.
	// Returns a boolean indicating whether the key exists in the cache and the value is an integer.
	ExistsInt(key string) bool
	// ExistsBool checks if the cache contains a boolean value associated with the given key.
	// Returns a boolean indicating whether the key exists in the cache and the value is a boolean.
	ExistsBool(key string) bool
	// ExistsFloat checks if the cache contains a floating-point value associated with the given key.
	// Returns a boolean indicating whether the key exists in the cache and the value is a floating-point number.
	ExistsFloat(key string) bool
	// ExistsString checks if the cache contains a string value associated with the given key.
	// Returns a boolean indicating whether the key exists in the cache and the value is a string.
	ExistsString(key string) bool
	// ExistsAny checks if the cache contains a value associated with the given key.
	// Returns a boolean indicating whether the key exists in the cache.
	ExistsAny(key string) bool

	// Keys returns a slice of all the keys present in the cache.
	Keys() []string
	// KeysInt returns a slice of keys that have integer values associated with them in the cache.
	KeysInt() []string
	// KeysBool returns a slice of keys that have boolean values associated with them in the cache.
	KeysBool() []string
	// KeysFloat returns a slice of keys that have floating-point values associated with them in the cache.
	KeysFloat() []string
	// KeysString returns a slice of keys that have string values associated with them in the cache.
	KeysString() []string
	// KeysAny returns a slice of keys that have values of any type associated with them in the cache.
	KeysAny() []string

	// Flush clears all entries in the cache.
	Flush()

	// // Count returns the number of items currently in the cache.
	// // This can be useful for monitoring cache usage or for debugging.
	// // Returns the count of items and an error if the operation fails.
	// Count() int

	// // Increment atomically increases the integer value of a key by delta.
	// // The function returns the new value after the increment and an error if the operation fails or if the value is not an integer.
	// Increment(key string, delta int64) (int64, error)

	// // Decrement atomically decreases the integer value of a key by delta.
	// // The function returns the new value after the decrement and an error if the operation fails or if the value is not an integer.
	// Decrement(key string, delta int64) (int64, error)
}

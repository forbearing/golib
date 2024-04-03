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
	// Pass T to create one record,
	// Pass []M to create multiple record.
	// It will update the "created_at" and "updated_at" field.
	Create(objs ...M) error
	// Delete one or multiple record.
	// Pass T to delete one record.
	// Pass []M to delete multiple record.
	Delete(objs ...M) error
	// Update one or multiple record, if record doesn't exist, it will be created.
	// Pass T to update one record.
	// Pass []M to update multiple record.
	// It will update the "created_at" and "updated_at" field.
	Update(objs ...M) error
	// UpdateById only update one record with specific id.
	UpdateById(id any, key string, value any) error

	// List all record and write to dest.
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
	// database.WithQuery(xxx).List(xxx) equal to database.Find(xxx).
	WithQuery(query M, fuzzyMatch ...bool) Database[M]

	// WithQueryRaw is where condition.
	// database.WithQueryRaw(xxx) same as database.WithQuery(xxx) and provides more flexible query.
	WithQueryRaw(query any, args ...any) Database[M]

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

	// WithPurge
	WithPurge(...bool) Database[M]
	// WithCache
	WithCache(...bool) Database[M]
	// WithOmit
	WithOmit(...string) Database[M]
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

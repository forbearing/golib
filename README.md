## Description

🚀 Golang Lightning Backend Framework

WARNING: Library under active development - expect significant API changes.

## Examples

1.  [basic usage example](./examples/simple/main.go)

2.  [example external project](https://github.com/forbearing/glpi)



## Documents

-   [Router usage](./examples/demo/main.go)
-   [Model usage](./examples/demo/model/user.go)
-   [Database Create](./examples/demo/controller/user_create.go)
-   [Database Delete](./examples/demo/controller/user_delete.go)
-   [Database Update](./examples/demo/controller/user_update.go)
-   [Database List](./examples/demo/controller/user_list.go)
-   [Database Get](./examples/demo/controller/user_get.go)
-   [Controller usage](./controller/READMD.md)
-   [Service usage](./examples/demo/service/user.go)
-   [Client usage](./pkg/client/client_test.go)
-   [tunnel usage](./tunnel/session_test.go)
-   Package usage
    -   lru
    -   cmap
    -   sqlite,postgres,mysql
    -   redis
    -   elastic
    -   mongo
    -   minio
    -   mqtt
    -   task



## Data Structure

-   [list](./ds/list)
    -   [arraylist](./ds/list/arraylist/list.go)
    -   [linkedlist](./ds/list/linkedlist/list.go)
    -   [skiplist](./ds/list/skiplist/skiplist.go)
-   [stack](./ds/stack)
    -   [arraystack](./ds/stack/arraystack/stack.go)
    -   [linkedstack](./ds/stack/linkedstack/stack.go)
-   [queue](./ds/queue)
    -   [arrayqueue](./ds/queue/arrayqueue/queue.go)
    -   [linkedqueue](./ds/queue/linkedqueue/queue.go)
    -   [priorityqueue](./ds/queue/priorityqueue/queue.go)
    -   [circularbuffer](./ds/queue/circularbuffer/circularbuffer.go)
-   [tree](./ds/tree)
    -   [rbtree](./ds/tree/rbtree/rbtree.go)
    -   [avltree](./ds/tree/avltree/avltree.go)
    -   [splaytree](./ds/tree/splaytree/splaytree.go)
    -   [trie](./ds/tree/trie/trie.go)
-   [heap](./ds/heap)
    -   [binaryheap](./ds/heap/binaryheap/binaryheap.go)
-   [mapset](./ds/mapset/set.go)
-   [multimap](./ds/multimap/multimap.go)


## Interface

### Logger

```go
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
```

### Database

```go
type Database[M Model] interface {
	Create(objs ...M) error
	Delete(objs ...M) error
	Update(objs ...M) error
	UpdateById(id string, key string, value any) error
	List(dest *[]M, cache ...*[]byte) error
	Get(dest M, id string, cache ...*[]byte) error
	First(dest M, cache ...*[]byte) error
	Last(dest M, cache ...*[]byte) error
	Take(dest M, cache ...*[]byte) error
	Count(*int64) error
	Cleanup() error
	Health() error

	DatabaseOption[M]
}

type DatabaseOption[M Model] interface {
	WithDB(any) Database[M]
	WithTable(name string) Database[M]
	WithDebug() Database[M]
	WithQuery(query M, fuzzyMatch ...bool) Database[M]
	WithQueryRaw(query any, args ...any) Database[M]
	WithAnd(...bool) Database[M]
	WithOr(...bool) Database[M]
	WithTimeRange(columnName string, startTime time.Time, endTime time.Time) Database[M]
	WithSelect(columns ...string) Database[M]
	WithSelectRaw(query any, args ...any) Database[M]
	WithIndex(index string) Database[M]
	WithTransaction(tx any) Database[M]
	WithJoinRaw(query string, args ...any) Database[M]
	WithLock(mode ...string) Database[M]
	WithBatchSize(size int) Database[M]
	WithScope(page, size int) Database[M]
	WithLimit(limit int) Database[M]
	WithExclude(map[string][]any) Database[M]
	WithOrder(order string) Database[M]
	WithExpand(expand []string, order ...string) Database[M]
	WithPurge(...bool) Database[M]
	WithCache(...bool) Database[M]
	WithOmit(...string) Database[M]
	WithTryRun(...bool) Database[M]
	WithoutHook() Database[M]
}
```

### Modal,Service

```go
type Model interface {
	GetTableName() string
	GetID() string
	SetID(id ...string)
	GetCreatedBy() string
	GetUpdatedBy() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreatedBy(string)
	SetUpdatedBy(string)
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
	Expands() []string
	Excludes() map[string][]any
	MarshalLogObject(zapcore.ObjectEncoder) error

	Hooker
}

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
```

### Service

```go
type Service[M Model] interface {
	CreateBefore(*ServiceContext, M) error
	CreateAfter(*ServiceContext, M) error
	DeleteBefore(*ServiceContext, M) error
	DeleteAfter(*ServiceContext, M) error
	UpdateBefore(*ServiceContext, M) error
	UpdateAfter(*ServiceContext, M) error
	UpdatePartialBefore(*ServiceContext, M) error
	UpdatePartialAfter(*ServiceContext, M) error
	ListBefore(*ServiceContext, *[]M) error
	ListAfter(*ServiceContext, *[]M) error
	GetBefore(*ServiceContext, M) error
	GetAfter(*ServiceContext, M) error

	BatchCreateBefore(*ServiceContext, ...M) error
	BatchCreateAfter(*ServiceContext, ...M) error
	BatchDeleteBefore(*ServiceContext, ...M) error
	BatchDeleteAfter(*ServiceContext, ...M) error
	BatchUpdateBefore(*ServiceContext, ...M) error
	BatchUpdateAfter(*ServiceContext, ...M) error
	BatchUpdatePartialBefore(*ServiceContext, ...M) error
	BatchUpdatePartialAfter(*ServiceContext, ...M) error

	Import(*ServiceContext, io.Reader) ([]M, error)
	Export(*ServiceContext, ...M) ([]byte, error)

	Filter(*ServiceContext, M) M
	FilterRaw(*ServiceContext) string

	Logger
}
```

### RBAC

```go
type RBAC interface {
	AddRole(name string) error
	RemoveRole(name string) error

	GrantPermission(role string, resource string, action string) error
	RevokePermission(role string, resource string, action string) error

	AssignRole(subject string, role string) error
	UnassignRole(subject string, role string) error
}
```

### Cache

```go
type Cache[T any] interface {
	Set(key string, values T)
	Get(key string) (T, bool)
	Peek(key string) (T, bool)
	Remove(key string)
	Exists(key string) bool
	Keys() []string
	Count() int
	Flush()
}
```

### ESDocumenter

```go
type ESDocumenter interface {
	Document() map[string]any
	GetID() string
}
```

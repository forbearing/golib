## Description

## Examples

1.   [basic usage example](./examples/simple)
2.   [simple project](./examples/myproject)

## Datatabase operation

### Create

```go
database.Database[M].WithExpand(req.Expands()).Create(req)
```

### Delete

```go
database.Database[M].WithExpand(req.Expands()).Create(req)
```

### Update/update_partial

```go
database.Database[M].Update(req)
```

### List

```go
database.Database[M].
  WithScope(page, size).
  WithOr(or).
  WithQuery(svc.Filter(svcCtx, m), fuzzy).
  WithQueryRaw(svc.FilterRaw(svcCtx)).
  WithExclude(m.Excludes()).
  WithExpand(expands, sortBy).
  WithOrder(sortBy).
  WithTimeRange(columnName, startTime, endTime).
  WithCache(!nocache).
  List(&data)
```

### Get

```go
database.Database[M].WithExpand(expands).WithCache(!nocache).Get(m, id)
```

## Router

```go
router.API.POST("/category", controller.Create[*model.Category])
router.API.DELETE("/category", controller.Delete[*model.Category])
router.API.DELETE("/category/:id", controller.Delete[*model.Category])
router.API.PUT("/category", controller.Update[*model.Category])
router.API.PUT("/category/:id", controller.Update[*model.Category])
router.API.PATCH("/category", controller.UpdatePartial[*model.Category])
router.API.PATCH("/category/:id", controller.UpdatePartial[*model.Category])
router.API.GET("/category", controller.List[*model.Category])
router.API.GET("/category/:id", controller.Get[*model.Category])
router.API.GET("/category/export", controller.Export[*model.Category])
router.API.POST("/category/import", controller.Import[*model.Category])
```

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

type Hooker interface {
	CreateBefore() error
	CreateAfter() error
	DeleteBefore() error
	DeleteAfter() error
	UpdateBefore() error
	UpdateAfter() error
	UpdatePartialBefore() error
	UpdatePartialAfter() error
	ListBefore() error
	ListAfter() error
	GetBefore() error
	GetAfter() error
}
```

### Service

```go
type Service[M Model] interface {
	CreateBefore(*ServiceContext, ...M) error
	CreateAfter(*ServiceContext, ...M) error
	DeleteBefore(*ServiceContext, ...M) error
	DeleteAfter(*ServiceContext, ...M) error
	UpdateBefore(*ServiceContext, ...M) error
	UpdateAfter(*ServiceContext, ...M) error
	UpdatePartialBefore(*ServiceContext, ...M) error
	UpdatePartialAfter(*ServiceContext, ...M) error
	ListBefore(*ServiceContext, *[]M) error
	ListAfter(*ServiceContext, *[]M) error
	GetBefore(*ServiceContext, ...M) error
	GetAfter(*ServiceContext, ...M) error
	Import(*ServiceContext, io.Reader) ([]M, error)
	Export(*ServiceContext, ...M) ([]byte, error)
	Filter(*ServiceContext, M) M
	FilterRaw(*ServiceContext) string

	Logger
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







## TODO

- [ ] Use reflect caching to optimizeize performance
- [ ] Configuration `config/config.go` supports multiple instances, including MySQL, PostgresSQL, Elastic, Elasticsearch, MongoDB, etc.
- [ ] Add `ldap` package and integrate it into `boostrap.Boostrap`
- [ ] Controler layer utilizes github.com/araddon/dateparse to handle arbitrary date time formats.

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
-   Service usage
-   Logger usage
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

## [tunnel usage](./tunnel/session_test.go)

```go
package tunnel_test

import (
	"net"
	"testing"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/tunnel"
	"github.com/forbearing/golib/types/consts"
	"github.com/stretchr/testify/assert"
)

var (
	addr   = "0.0.0.0:12345"
	doneCh = make(chan struct{}, 1)

	Bye   = tunnel.NewCmd("bye", 1000)
	Hello = tunnel.NewCmd("hello", 1001)
)

type ByePaylod struct {
	Field1 string
	Field2 uint64
}
type HelloPaylod struct {
	Field3 string
	Field4 float64
}

var (
	byePayload1   = ByePaylod{Field1: "bye1", Field2: 123}
	byePayload2   = ByePaylod{Field1: "bye2", Field2: 456}
	helloPayload1 = HelloPaylod{Field3: "hello1", Field4: 3.14}
	helloPayload2 = HelloPaylod{Field3: "hello2", Field4: 3.14}
)

func TestSession(t *testing.T) {
	assert.NoError(t, bootstrap.Bootstrap())
	go server(t)
	client(t)
	<-doneCh
}

func server(t *testing.T) {
	l, err := net.Listen("tcp", addr)
	assert.NoError(t, err)
	defer l.Close()

	doneCh <- struct{}{}

	conn, err := l.Accept()
	assert.NoError(t, err)
	defer conn.Close()

	session, _ := tunnel.NewSession(conn, consts.Server)
	for {
		event, err := session.Read()
		assert.NoError(t, err)
		switch event.Cmd {
		case tunnel.Ping:
			t.Log("client ping")
			session.Write(&tunnel.Event{Cmd: tunnel.Pong})
		case Hello:
			payload, err := tunnel.DecodePayload[HelloPaylod](event.Payload)
			assert.NoError(t, err)
			t.Logf("client hello: %+v\n", payload)
			session.Write(&tunnel.Event{Cmd: Hello, Payload: helloPayload2})
		case Bye:
			payload, err := tunnel.DecodePayload[ByePaylod](event.Payload)
			assert.NoError(t, err)
			t.Logf("client bye: %+v\n", payload)
			session.Write(&tunnel.Event{Cmd: Bye, Payload: byePayload2})
		}
	}
}

func client(t *testing.T) {
	<-doneCh

	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	defer conn.Close()

	session, _ := tunnel.NewSession(conn, consts.Client)
	session.Write(&tunnel.Event{Cmd: tunnel.Ping})

	for {
		event, err := session.Read()
		assert.NoError(t, err)

		switch event.Cmd {
		case tunnel.Pong:
			t.Log("server pong")
			session.Write(&tunnel.Event{Cmd: Hello, Payload: helloPayload1})
		case Hello:
			payload, err := tunnel.DecodePayload[HelloPaylod](event.Payload)
			assert.NoError(t, err)
			t.Logf("server hello: %+v\n", payload)
			session.Write(&tunnel.Event{Cmd: Bye, Payload: byePayload1})
		case Bye:
			payload, err := tunnel.DecodePayload[ByePaylod](event.Payload)
			assert.NoError(t, err)
			t.Logf("server bye: %+v\n", payload)
			doneCh <- struct{}{}
			return
		}
	}
}
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
- [ ] Add more middleware

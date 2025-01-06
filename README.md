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

## [client usage](./pkg/client/client_test.go)

```go
package client_test

import (
	"fmt"
	"testing"

	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/pkg/client"
	"github.com/stretchr/testify/assert"
)

var (
	token = "-"
	addr2 = "http://localhost:8002/api/user//"

	id1     = "user1"
	id2     = "user2"
	id3     = "user3"
	id4     = "user4"
	id5     = "user5"
	name1   = id1
	name2   = id2
	name3   = id3
	name4   = id4
	name5   = id5
	email1  = "user1@gmail.com"
	email2  = "user2@gmail.com"
	email3  = "user3@gmail.com"
	email4  = "user4@gmail.com"
	email5  = "user5@gmail.com"
	avatar1 = "avatar1"
	avatar2 = "avatar2"
	avatar3 = "avatar3"
	avatar4 = "avatar4"
	avatar5 = "avatar5"

	name1Modified   = id1 + "_modified"
	email1Modified  = email1 + "_modified"
	avatar1Modified = avatar1 + "_modified"

	user1 = User{Name: name1, Email: email1, Avatar: avatar1, Base: model.Base{ID: id1}}
	user2 = User{Name: name2, Email: email2, Avatar: avatar2, Base: model.Base{ID: id2}}
	user3 = User{Name: name3, Email: email3, Avatar: avatar3, Base: model.Base{ID: id3}}
	user4 = User{Name: name4, Email: email4, Avatar: avatar4, Base: model.Base{ID: id4}}
	user5 = User{Name: name5, Email: email5, Avatar: avatar5, Base: model.Base{ID: id5}}
)

func TestClient(t *testing.T) {
	client, err := client.New(addr2, client.WithToken(token), client.WithQueryPagination(1, 2))
	assert.NoError(t, err)
	fmt.Println(client.QueryString())
	fmt.Println(client.RequestURL())

	client.Create(user1)
	client.Create(user2)
	client.Create(user3)
	client.Create(user4)
	client.Create(user5)

	users := make([]User, 0)
	total := new(int64)
	user := new(User)

	// test List
	{
		assert.NoError(t, client.List(&users, total))
		assert.Equal(t, 2, len(users))
		assert.Equal(t, int64(5), *total)
	}

	// test Get
	{
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1, user.Name)
		assert.Equal(t, email1, user.Email)
		assert.Equal(t, avatar1, user.Avatar)
		assert.NoError(t, client.Get(id2, user))
		assert.Equal(t, id2, user.ID)
		assert.Equal(t, name2, user.Name)
		assert.Equal(t, email2, user.Email)
		assert.Equal(t, avatar2, user.Avatar)
		assert.NoError(t, client.Get(id3, user))
		assert.Equal(t, id3, user.ID)
		assert.Equal(t, name3, user.Name)
		assert.Equal(t, email3, user.Email)
		assert.Equal(t, avatar3, user.Avatar)
		assert.NoError(t, client.Get(id4, user))
		assert.Equal(t, id4, user.ID)
		assert.Equal(t, name4, user.Name)
		assert.Equal(t, email4, user.Email)
		assert.Equal(t, avatar4, user.Avatar)
		assert.NoError(t, client.Get(id5, user))
		assert.Equal(t, id5, user.ID)
		assert.Equal(t, name5, user.Name)
		assert.Equal(t, email5, user.Email)
		assert.Equal(t, avatar5, user.Avatar)
	}

	// Test Update
	{
		_, err = client.Update(&User{Name: name1Modified, Email: email1Modified, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1Modified, user.Name)
		assert.Equal(t, email1Modified, user.Email)
		assert.Empty(t, user.Avatar)
	}

	// Test UpdatePartial
	{
		_, err = client.UpdatePartial(&User{Avatar: avatar1Modified, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1Modified, user.Name)
		assert.Equal(t, email1Modified, user.Email)
		assert.Equal(t, avatar1Modified, user.Avatar)

		_, err = client.UpdatePartial(&User{Name: name1, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1, user.Name)
		assert.Equal(t, email1Modified, user.Email)
		assert.Equal(t, avatar1Modified, user.Avatar)

		_, err = client.UpdatePartial(&User{Email: email1, Avatar: avatar1, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1, user.Name)
		assert.Equal(t, email1, user.Email)
		assert.Equal(t, avatar1, user.Avatar)

	}
}

type User struct {
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Avatar string `json:"avatar,omitempty"`

	model.Base
}
```

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
			payload := new(HelloPaylod)
			assert.NoError(t, tunnel.DecodePayload(event.Payload, payload))
			assert.Equal(t, helloPayload1, *payload)
			t.Logf("client hello: %+v\n", *payload)
			session.Write(&tunnel.Event{Cmd: Hello, Payload: helloPayload2})
		case Bye:
			payload := new(ByePaylod)
			assert.NoError(t, tunnel.DecodePayload(event.Payload, payload))
			assert.Equal(t, byePayload1, *payload)
			t.Logf("client bye: %+v\n", *payload)
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
			payload := new(HelloPaylod)
			assert.NoError(t, tunnel.DecodePayload(event.Payload, payload))
			assert.Equal(t, helloPayload2, *payload)
			t.Logf("server hello: %+v\n", *payload)
			session.Write(&tunnel.Event{Cmd: Bye, Payload: byePayload1})
		case Bye:
			payload := new(ByePaylod)
			assert.NoError(t, tunnel.DecodePayload(event.Payload, payload))
			assert.Equal(t, byePayload2, *payload)
			t.Logf("server bye: %+v\n", *payload)
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

## Description

🚀 Golang Lightning Backend Framework

文档正在快马加鞭补充中



## 文档

### 1.路由使用

### 2.controller 使用

### 4.model 使用

### 4.service 使用

### 4.数据库操作



## Full Example

### main.go

```go
package main

import (
	"net/http"
	"time"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/middleware"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/task"
	. "github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Prepare
	// Setup configuration.
	config.SetConfigFile("./config.ini")
	config.SetConfigName("config")
	config.SetConfigType("ini")
	// Add tasks.
	task.Register(SayHello, 1*time.Second, "say hello")
	task.Register(SayGoodbye, 1*time.Second, "say goodbye")
	RunOrDie(bootstrap.Bootstrap)

	zap.S().Infow("successfully initialized", "addr", AppConf.MqttConfig.Addr, "username", AppConf.MqttConfig.Username)
	logger.Controller.Infow("successfully initialized", "addr", AppConf.MqttConfig.Addr, "username", AppConf.MqttConfig.Username)
	logger.Service.Infow("successfully initialized", "addr", AppConf.MqttConfig.Addr, "username", AppConf.MqttConfig.Username)

	// use Base router.
	router.Base.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	router.Base.GET("/hello", func(c *gin.Context) { c.String(http.StatusOK, "hello world!") })

	// without auth
	router.API.GET("/noauth/user", controller.List[*User])
	router.API.GET("/noauth/user/:id", controller.Get[*User])
	router.API.Use(middleware.JwtAuth(), middleware.RateLimiter())

	// with auth
	router.API.POST("/user", controller.Create[*User])
	router.API.DELETE("/user", controller.Delete[*User])
	router.API.DELETE("/user/:id", controller.Delete[*User])
	router.API.PUT("/user", controller.Update[*User])
	router.API.PUT("/user/:id", controller.Update[*User])
	router.API.PATCH("/user", controller.UpdatePartial[*User])
	router.API.PATCH("/user/:id", controller.UpdatePartial[*User])
	router.API.GET("/user", controller.List[*User])
	router.API.GET("/user/:id", controller.Get[*User])
	router.API.GET("/user/export", controller.Export[*User])
	router.API.POST("/user/import", controller.Import[*User])

	router.API.POST("/group", controller.Create[*Group])
	router.API.DELETE("/group", controller.Delete[*Group])
	router.API.DELETE("/group/:id", controller.Delete[*Group])
	router.API.PUT("/group", controller.Update[*Group])
	router.API.PUT("/group/:id", controller.Update[*Group])
	router.API.PATCH("/group", controller.UpdatePartial[*Group])
	router.API.PATCH("/group/:id", controller.UpdatePartial[*Group])
	router.API.GET("/group", controller.List[*Group])
	router.API.GET("/group/:id", controller.Get[*Group])
	router.API.GET("/group/export", controller.Export[*Group])
	router.API.POST("/group/import", controller.Import[*Group])

	// Run server.
	RunOrDie(router.Run)
}


var AppConf = new(Config)

func InitConfig() (err error) {
	config.SetDefaultValue()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(AppConf); err != nil {
		return
	}
	return nil
}

type Config struct {
	MqttConfig            `json:"mqtt" mapstructure:"mqtt" ini:"mqtt" yaml:"mqtt"`
	config.ServerConfig   `json:"server" mapstructure:"server" ini:"server" yaml:"server"`
	config.AuthConfig     `json:"auth" mapstructure:"auth" ini:"auth" yaml:"auth"`
	config.SqliteConfig   `json:"sqlite" mapstructure:"sqlite" ini:"sqlite" yaml:"sqlite"`
	config.PostgreConfig  `json:"postgres" mapstructure:"postgres" ini:"postgres" yaml:"postgres"`
	config.MySQLConfig    `json:"mysql" mapstructure:"mysql" ini:"mysql" yaml:"mysql"`
	config.RedisConfig    `json:"redis" mapstructure:"redis" ini:"redis" yaml:"redis"`
	config.MinioConfig    `json:"minio" mapstructure:"minio" ini:"minio" yaml:"minio"`
	config.S3Config       `json:"s3" mapstructure:"s3" ini:"s3" yaml:"s3"`
	config.LoggerConfig   `json:"logger" mapstructure:"logger" ini:"logger" yaml:"logger"`
	config.LdapConfig     `json:"ldap" mapstructure:"ldap" ini:"ldap" yaml:"ldap"`
	config.InfluxdbConfig `json:"influxdb" mapstructure:"influxdb" ini:"influxdb" yaml:"influxdb"`
	config.FeishuConfig   `json:"feishu" mapstructure:"feishu" ini:"feishu" yaml:"feishu"`
}

type MqttConfig struct {
	Addr     string `json:"addr" mapstructure:"addr" ini:"addr" yaml:"addr"`
	Username string `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
}
```

### model.go

```go
package main

import "github.com/forbearing/golib/model"

func init() {
	model.Register[*User]()
	model.Register[*Group]()
}

type User struct {
	Name   string `json:"name,omitempty" schema:"name" gorm:"unique" binding:"required"`
	Email  string `json:"email,omitempty" schema:"email" gorm:"unique" binding:"required"`
	Avatar string `json:"avatar,omitempty" schema:"avatar"`

	model.Base
}

type Group struct {
	Name string `json:"name,omitempty" schema:"name" gorm:"unique" binding:"required"`

	model.Base
}
```

### task.go

```go
func SayHello() error {
	// fmt.Println("hello world!")
	logger.Task.Info("hello world!")
	return nil
}

func SayGoodbye() error {
	// fmt.Println("goodbye world!")
	logger.Task.Info("goodbye world!")
	return nil
}
```



## config example

```ini
[server]
mode = dev
port = 8002
# token_expire_duration = 12h

[auth]
none_expire_token = "-"

[logger]
log_level = info
log_file = ""
# log_format = "console"

[sqlite]
path = "/tmp/data.db"
; is_memory = true

[postgres]
port = 15432
username = "postgres"
password = "admin"

[mysql]
database = mydb
password = random_password

[redis]
host = localhost
port = 6379
password = random_password
enable = false
# expiration = 10m

[minio]
endpoint = localhost:9000
access_key = my_access_key
secret_key = my_secret_key 
region = shjd-oss
bucket = asset
use_ssl = false

[mqtt]
addr = tcp://localhost:1883
username = myuser
password = mypass
```



## Datatabase operation

### Create

```go
if err := database.Database[M].WithExpand(req.Expands()).Create(req); err != nil {
  log.Error(err)
  ResponseJSON(c, CodeFailure)
  return
}
```

### Delete

```go
if err := database.Database[M].WithExpand(req.Expands()).Create(req); err != nil {
  log.Error(err)
  ResponseJSON(c, CodeFailure)
  return
}
```

### Update/update_partial

```go
if err := database.Database[M].Update(req); err != nil {
  log.Error(err)
  ResponseJSON(c, CodeFailure)
  return
}
```

### List

```go
if err = database.Database[M].
  WithScope(page, size).
  WithOr(or).
  WithQuery(svc.Filter(svcCtx, m), fuzzy).
  WithQueryRaw(svc.FilterRaw(svcCtx)).
  WithExclude(m.Excludes()).
  WithExpand(expands, sortBy).
  WithOrder(sortBy).
  WithTimeRange(columnName, startTime, endTime).
  WithCache(!nocache).
  List(&data, &cache); err != nil {
  log.Error(err)
  ResponseJSON(c, CodeFailure)
  return
}
```

### Get

```go
if err = database.Database[M].WithExpand(expands).WithCache(!nocache).Get(m, c.Param(PARAM_ID), &cache); err != nil {
  log.Error(err)
  ResponseJSON(c, CodeFailure)
  return
}
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
	With(key, value string) Logger

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
	UpdateById(id any, key string, value any) error
	List(dest *[]M, cache ...*[]byte) error
	Get(dest M, id string, cache ...*[]byte) error
	First(dest M, cache ...*[]byte) error
	Last(dest M, cache ...*[]byte) error
	Take(dest M, cache ...*[]byte) error
	Count(*int64) error
	Cleanup() error

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
	WithIndex(index string) Database[M]
	WithBatchSize(size int) Database[M]
	WithScope(page, size int) Database[M]
	WithLimit(limit int) Database[M]
	WithExclude(map[string][]any) Database[M]
	WithOrder(order string) Database[M]
	WithExpand(expand []string, order ...string) Database[M]
	WithPurge(...bool) Database[M]
	WithCache(...bool) Database[M]
	WithOmit(...string) Database[M]
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

### Cache

```go
type Cache[T any] interface {
	Set(key string, values T)
	Get(key string) (T, bool)
	Remove(key string)
	Exists(key string) bool
	Keys() []string
	Count() int
	Flush()
}
```

## TODO

- [x] database support postgresql
- [x] database support sqlite
- [ ] dateparse parse anytime \_start_time, \_end_time
- [ ] limit recursive query/update in Hook.
- [ ] config support toml
- [ ] Join
- [x] WithSelect, WithIndex
- [ ] frontend

## Description

Golang library

## Full Example

### main.go

```go
package main

import (
	"time"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/database/cache"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/database/redis"
	"github.com/forbearing/golib/examples/myproject/model"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/logger/logrus"
	pkgzap "github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/metrics"
	"github.com/forbearing/golib/middleware"
	"github.com/forbearing/golib/minio"
	"github.com/forbearing/golib/rbac"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/task"
	. "github.com/forbearing/golib/util"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 1.Register initializer.
	bootstrap.Register(
		config.Init,
		InitConfig,
		pkgzap.Init,
		logrus.Init,
		metrics.Init,
		cache.Init,
		mysql.Init,
		redis.Init,
		rbac.Init,
		service.Init,
		minio.Init,
		router.Init,
		task.Init,
	)
	bootstrap.RegisterGo(router.Run)

	// 2.Prepare
	// Setup configuration.
	config.SetConfigFile("./config.ini")
	config.SetConfigName("config")
	config.SetConfigType("ini")
	// Add tasks.
	task.Register(SayHello, 1*time.Second, "say hello")
	task.Register(SayGoodbye, 1*time.Second, "say goodbye")

	// 3.Initialize.
	RunOrDie(bootstrap.Init)

	zap.S().Infow("successfully initialized", "addr", AppConf.MqttConfig.Addr, "username", AppConf.MqttConfig.Username)
	logger.Controller.Infow("successfully initialized", "addr", AppConf.MqttConfig.Addr, "username", AppConf.MqttConfig.Username)
	logger.Service.Infow("successfully initialized", "addr", AppConf.MqttConfig.Addr, "username", AppConf.MqttConfig.Username)

	// 4.setup apis.
	// without auth
	router.API.GET("/noauth/category", controller.List[*model.Category])
	router.API.GET("/noauth/category/:id", controller.Get[*model.Category])
	router.API.Use(middleware.JwtAuth(), middleware.RateLimiter())
	// with auth
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

	// 5.Run server.
	RunOrDie(bootstrap.Go)
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
	MqttConfig            `json:"mqtt" mapstructure:"mqtt" init:"mqtt" yaml:"mqtt"`
	config.ServerConfig   `json:"server" mapstructure:"server" ini:"server" yaml:"server"`
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
	Addr     string `json:"addr" mapstructure:"addr" init:"addr" yaml:"addr"`
	Username string `json:"username" mapstructure:"username" init:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" init:"password" yaml:"password"`
}
```

### model/category.go

```go
package model

import (
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"go.uber.org/zap/zapcore"
)

var (
	// 根分类, 不会展示出来 所有的一级分类的 parentId 都是这个根分类.
	// 数据库初始化时, 会自动创建这个记录
	categoryRoot = &Category{Name: RootName, Status: util.Pointer(uint(0)), ParentId: RootId, Base: model.Base{ID: RootId}}
	// 未知分类, 导入数据时,如果发现 parentId 为空,将将其 parentId 指向这个未知分类
	// 数据库初始化时, 会自动创建这个记录
	categoryUnknown = &Category{Name: UnknownName, Status: util.Pointer(uint(0)), ParentId: UnknownId, Base: model.Base{ID: UnknownId}}
	categoryNone    = &Category{Name: NoneName, Status: util.Pointer(uint(0)), ParentId: RootId, Base: model.Base{ID: NoneId}}
)

func init() {
	model.Register[*Category](categoryRoot, categoryUnknown, categoryNone)
	model.RegisterRoutes[*Category]("category")
}

type Category struct {
	Name     string     `json:"name,omitempty" gorm:"unique" schema:"name"`
	Status   *uint      `json:"status,omitempty" gorm:"type:tinyint(1);comment:status(0: disabled, 1: enable)" schema:"status"`
	ParentId string     `json:"parent_id,omitempty" gorm:"size:191" schema:"parent_id"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentId"`
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentId;references:ID"`

	model.Base
}

func (*Category) Expands() []string {
	return []string{"Children", "Parent"}
}
func (*Category) Excludes() map[string][]any {
	return map[string][]any{KeyName: {RootName, UnknownName, NoneName}}
}
func (c *Category) CreateBefore() error {
	if len(c.ParentId) == 0 {
		c.ParentId = RootId
	}
	return nil
}
func (c *Category) UpdateBefore() error {
	return c.CreateBefore()
}

func (c *Category) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if c == nil {
		return nil
	}
	enc.AddString("id", c.ID)
	enc.AddString("name", c.Name)
	enc.AddUint("status", util.Depointer(c.Status))
	enc.AddString("parent_id", c.ParentId)
	enc.AddObject("base", &c.Base)
	return nil
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
	WithAnd(...bool) Database[M]
	WithOr(...bool) Database[M]
	WithQueryRaw(query any, args ...any) Database[M]
	WithTimeRange(columnName string, startTime time.Time, endTime time.Time) Database[M]
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


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


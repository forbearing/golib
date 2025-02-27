package main

import (
	"net/http"
	"time"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	pkgcontroller "github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/cronjob"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/examples/myproject/model"
	"github.com/forbearing/golib/middleware"
	pkgmodel "github.com/forbearing/golib/model"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/task"
	"github.com/forbearing/golib/types"
	. "github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	_ "github.com/forbearing/golib/examples/myproject/service"
)

func main() {
	// Prepare
	// Setup configuration.
	config.SetConfigFile("./config.ini")
	config.SetConfigName("config")
	config.SetConfigType("ini")

	// Register task and cronjob before bootstrap.
	task.Register(SayHello, 1*time.Second, "say hello")
	cronjob.Register(SayHello, "*/1 * * * * *", "say hello")

	RunOrDie(bootstrap.Bootstrap)

	// Register task and cronjob after bootstrap.
	task.Register(SayGoodbye, 1*time.Second, "say goodbye")
	cronjob.Register(SayGoodbye, "*/1 * * * * *", "say goodbye")

	zap.S().Infow("successfully initialized", "addr", AppConf.MqttConfig.Addr, "username", AppConf.MqttConfig.Username)

	// use Base router.
	router.Base.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	router.Base.GET("/hello", func(c *gin.Context) { c.String(http.StatusOK, "hello world!") })

	router.API.POST("/login", pkgcontroller.User.Login)
	router.API.POST("/signup", pkgcontroller.User.Signup)
	router.API.DELETE("/logout", middleware.JwtAuth(), pkgcontroller.User.Logout)
	router.API.POST("token/refresh", pkgcontroller.User.RefreshToken)
	router.API.GET("/debug/debug", Debug.Debug)

	router.API.Use(
		middleware.JwtAuth(),
		// middleware.RateLimiter(),
	)

	router.Register[*pkgmodel.User](router.API, "/user")
	router.Register[*model.Group](router.API, "/group")

	// router.API.POST("/user", controller.Create[*User])
	// router.API.DELETE("/user", controller.Delete[*User])
	// router.API.DELETE("/user/:id", controller.Delete[*User])
	// router.API.PUT("/user", controller.Update[*User])
	// router.API.PUT("/user/:id", controller.Update[*User])
	// router.API.PATCH("/user", controller.UpdatePartial[*User])
	// router.API.PATCH("/user/:id", controller.UpdatePartial[*User])
	// router.API.GET("/user", controller.List[*User])
	// router.API.GET("/user/:id", controller.Get[*User])
	// router.API.GET("/user/export", controller.Export[*User])
	// router.API.POST("/user/import", controller.Import[*User])

	// router.API.POST("/group", controller.Create[*Group])
	// router.API.DELETE("/group", controller.Delete[*Group])
	// router.API.DELETE("/group/:id", controller.Delete[*Group])
	// router.API.PUT("/group", controller.Update[*Group])
	// router.API.PUT("/group/:id", controller.Update[*Group])
	// router.API.PATCH("/group", controller.UpdatePartial[*Group])
	// router.API.PATCH("/group/:id", controller.UpdatePartial[*Group])
	// router.API.GET("/group", controller.List[*Group])
	// router.API.GET("/group/:id", controller.Get[*Group])
	// router.API.GET("/group/export", controller.Export[*Group])
	// router.API.POST("/group/import", controller.Import[*Group])

	cfg := config.MySQLConfig{}
	cfg.Host = "127.0.0.1"
	cfg.Port = 3306
	cfg.Database = "golib"
	cfg.Username = "golib"
	cfg.Password = "golib"
	cfg.Charset = "utf8mb4"
	db, err := mysql.New(cfg)
	if err != nil {
		panic(err)
	}
	// It's your responsibility to ensure the table already exists.
	router.RegisterWithConfig(&types.ControllerConfig[*pkgmodel.User]{DB: db}, router.API, "/external/user")
	router.RegisterWithConfig(&types.ControllerConfig[*model.Group]{DB: db}, router.API, "/external/group")

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

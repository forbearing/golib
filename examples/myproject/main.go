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

	router.API.GET("/debug/debug", Debug.Debug)

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

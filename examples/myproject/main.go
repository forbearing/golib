package main

import (
	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/database/cache"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/database/redis"
	"github.com/forbearing/golib/examples/myproject/model"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/logger/logrus"
	"github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/metrics"
	"github.com/forbearing/golib/minio"
	"github.com/forbearing/golib/rbac"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/task"
	. "github.com/forbearing/golib/util"
	"github.com/spf13/viper"
)

func main() {
	bootstrap.Register(
		config.Init,
		InitConfig,
		zap.Init,
		logrus.Init,
		metrics.Init,
		cache.Init,
		mysql.Init,
		redis.Init,
		rbac.Init,
		service.Init,
		minio.Init,
		task.Init,
		router.Init,
	)
	bootstrap.RegisterGo(router.Run)

	RunOrDie(bootstrap.Init)

	logger.Infow("successfully initialized",
		"addr", AppConf.MqttConfig.Addr, "username",
		AppConf.MqttConfig.Username, "password",
		AppConf.MqttConfig.Password)

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

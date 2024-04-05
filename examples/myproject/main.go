package main

import (
	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/database/redis"
	"github.com/forbearing/golib/examples/myproject/model"
	"github.com/forbearing/golib/logger/logrus"
	"github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/minio"
	"github.com/forbearing/golib/rbac"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/task"
	. "github.com/forbearing/golib/util"
)

func main() {
	bootstrap.Register(
		config.Init,
		zap.Init,
		logrus.Init,
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

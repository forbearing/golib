package bootstrap

import (
	"github.com/forbearing/golib/cmap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/database/postgres"
	"github.com/forbearing/golib/database/redis"
	"github.com/forbearing/golib/database/sqlite"
	"github.com/forbearing/golib/elastic"
	"github.com/forbearing/golib/logger/logrus"
	"github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/lru"
	"github.com/forbearing/golib/metrics"
	"github.com/forbearing/golib/minio"
	"github.com/forbearing/golib/mqtt"
	"github.com/forbearing/golib/rbac"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/task"
	"github.com/forbearing/golib/util"
)

var keys = []string{
	util.GetFunctionName(config.Init),
	util.GetFunctionName(zap.Init),
	util.GetFunctionName(logrus.Init),
	util.GetFunctionName(metrics.Init),
	util.GetFunctionName(cmap.Init),
	util.GetFunctionName(sqlite.Init),
	util.GetFunctionName(postgres.Init),
	util.GetFunctionName(mysql.Init),
	util.GetFunctionName(redis.Init),
	util.GetFunctionName(rbac.Init),
	util.GetFunctionName(service.Init),
	util.GetFunctionName(minio.Init),
	util.GetFunctionName(router.Init),
	util.GetFunctionName(task.Init),
}

func Bootstrap() error {
	// config.Init,
	// InitConfig,
	// pkgzap.Init,
	// logrus.Init,
	// metrics.Init,
	// cache.Init,
	// // mysql.Init,
	// // sqlite.Init,
	// postgres.Init,
	// redis.Init,
	// rbac.Init,
	// service.Init,
	// minio.Init,
	// router.Init,
	// task.Init,

	Register(
		config.Init,
		zap.Init,
		logrus.Init,
		metrics.Init,
		lru.Init,
		cmap.Init,
		// mysql.Init,
		// sqlite.Init,
		postgres.Init,
		elastic.Init,
		redis.Init,
		mqtt.Init,
		rbac.Init,
		service.Init,
		minio.Init,
		router.Init,
		task.Init,
	)

	return Init()
}

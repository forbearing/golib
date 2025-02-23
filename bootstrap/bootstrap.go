package bootstrap

import (
	"github.com/forbearing/golib/cmap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/database/postgres"
	"github.com/forbearing/golib/database/redis"
	"github.com/forbearing/golib/database/sqlite"
	"github.com/forbearing/golib/elastic"
	"github.com/forbearing/golib/jwt"
	"github.com/forbearing/golib/logger/logrus"
	"github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/lru"
	"github.com/forbearing/golib/metrics"
	"github.com/forbearing/golib/minio"
	"github.com/forbearing/golib/mongo"
	"github.com/forbearing/golib/mqtt"
	"github.com/forbearing/golib/rbac"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/task"
)

func Bootstrap() error {
	Register(
		config.Init,
		zap.Init,
		logrus.Init,
		metrics.Init,
		lru.Init,
		cmap.Init,
		sqlite.Init,
		postgres.Init,
		mysql.Init,
		redis.Init,
		elastic.Init,
		mongo.Init,
		minio.Init,
		mqtt.Init,
		rbac.Init,
		service.Init,
		router.Init,
		jwt.Init,
		task.Init,
	)

	return Init()
}

package bootstrap

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/forbearing/golib/cache/cmap"
	"github.com/forbearing/golib/cache/lru"
	"github.com/forbearing/golib/cache/redis"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/cronjob"
	"github.com/forbearing/golib/database/clickhouse"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/database/postgres"
	"github.com/forbearing/golib/database/sqlite"
	"github.com/forbearing/golib/database/sqlserver"
	"github.com/forbearing/golib/debug/gops"
	"github.com/forbearing/golib/debug/pprof"
	"github.com/forbearing/golib/debug/statsviz"
	"github.com/forbearing/golib/grpc"
	"github.com/forbearing/golib/jwt"
	"github.com/forbearing/golib/logger/logrus"
	"github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/metrics"
	"github.com/forbearing/golib/provider/cassandra"
	"github.com/forbearing/golib/provider/elastic"
	"github.com/forbearing/golib/provider/etcd"
	"github.com/forbearing/golib/provider/influxdb"
	"github.com/forbearing/golib/provider/kafka"
	"github.com/forbearing/golib/provider/minio"
	"github.com/forbearing/golib/provider/mongo"
	"github.com/forbearing/golib/provider/mqtt"
	"github.com/forbearing/golib/provider/nats"
	"github.com/forbearing/golib/rbac"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/task"
	"go.uber.org/automaxprocs/maxprocs"
)

var (
	initialized bool
	mu          sync.Mutex
)

func Bootstrap() error {
	maxprocs.Set(maxprocs.Logger(zap.New().Infof))

	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return nil
	}

	Register(
		config.Init,
		zap.Init,
		logrus.Init,
		metrics.Init,

		// cache
		lru.Init,
		cmap.Init,
		redis.Init,

		// database
		sqlite.Init,
		postgres.Init,
		mysql.Init,
		clickhouse.Init,
		sqlserver.Init,

		// provider
		elastic.Init,
		mongo.Init,
		minio.Init,
		nats.Init,
		mqtt.Init,
		etcd.Init,
		nats.Init,
		kafka.Init,
		cassandra.Init,
		influxdb.Init,

		// service
		rbac.Init,
		service.Init,
		router.Init,
		grpc.Init,
		jwt.Init,

		// job
		task.Init,
		cronjob.Init,
	)

	RegisterExitHandler(redis.Close)
	RegisterExitHandler(influxdb.Close)
	RegisterExitHandler(nats.Close)
	RegisterExitHandler(kafka.Close)
	RegisterExitHandler(etcd.Close)

	initialized = true
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigCh
		Cleanup()
	}()

	return Init()
}

func Run() error {
	RegisterGo(
		router.Run,
		grpc.Run,
		statsviz.Run,
		pprof.Run,
		gops.Run,
	)
	return Go()
}

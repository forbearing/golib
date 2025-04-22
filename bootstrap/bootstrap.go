package bootstrap

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/forbearing/golib/authn/jwt"
	"github.com/forbearing/golib/authz/rbac/basic"
	"github.com/forbearing/golib/authz/rbac/tenant"
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
	"github.com/forbearing/golib/logger/logrus"
	pkgzap "github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/metrics"
	"github.com/forbearing/golib/middleware"
	"github.com/forbearing/golib/provider/cassandra"
	"github.com/forbearing/golib/provider/elastic"
	"github.com/forbearing/golib/provider/etcd"
	"github.com/forbearing/golib/provider/feishu"
	"github.com/forbearing/golib/provider/influxdb"
	"github.com/forbearing/golib/provider/kafka"
	"github.com/forbearing/golib/provider/ldap"
	"github.com/forbearing/golib/provider/memcached"
	"github.com/forbearing/golib/provider/minio"
	"github.com/forbearing/golib/provider/mongo"
	"github.com/forbearing/golib/provider/mqtt"
	"github.com/forbearing/golib/provider/nats"
	"github.com/forbearing/golib/provider/rethinkdb"
	"github.com/forbearing/golib/provider/rocketmq"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/service"
	service_authz "github.com/forbearing/golib/service/authz"
	service_log "github.com/forbearing/golib/service/log"
	"github.com/forbearing/golib/task"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

var (
	initialized bool
	mu          sync.Mutex
)

func Bootstrap() error {
	maxprocs.Set(maxprocs.Logger(pkgzap.New().Infof))

	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return nil
	}

	Register(
		config.Init,
		pkgzap.Init,
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
		kafka.Init,
		etcd.Init,
		nats.Init,
		cassandra.Init,
		influxdb.Init,
		memcached.Init,
		rethinkdb.Init,
		rocketmq.Init,
		feishu.Init,
		ldap.Init,

		// Authorization and Authentication
		basic.Init,
		tenant.Init,
		jwt.Init,

		// service
		service.Init,
		service_authz.Init,
		service_log.Init,

		middleware.Init,
		router.Init,
		grpc.Init,

		// job
		task.Init,
		cronjob.Init,
	)

	RegisterCleanup(config.Clean)
	RegisterCleanup(pkgzap.Clean)
	RegisterCleanup(redis.Close)
	RegisterCleanup(kafka.Close)
	RegisterCleanup(etcd.Close)
	RegisterCleanup(nats.Close)
	RegisterCleanup(cassandra.Close)
	RegisterCleanup(influxdb.Close)
	RegisterCleanup(memcached.Close)
	RegisterCleanup(rethinkdb.Close)
	RegisterCleanup(rocketmq.Close)
	RegisterCleanup(ldap.Close)

	initialized = true

	return Init()
}

func Run() error {
	defer Cleanup()

	RegisterGo(
		router.Run,
		grpc.Run,
		statsviz.Run,
		pprof.Run,
		gops.Run,
	)

	RegisterCleanup(router.Stop)
	RegisterCleanup(grpc.Stop)
	RegisterCleanup(statsviz.Stop)
	RegisterCleanup(pprof.Stop)
	RegisterCleanup(gops.Stop)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	errCh := make(chan error, 1)

	go func() {
		errCh <- Go()
	}()
	select {
	case sig := <-sigCh:
		zap.S().Infow("cancelled by signal", "signal", sig)
		return nil
	case err := <-errCh:
		return err
	}
}

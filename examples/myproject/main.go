package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/cache/redis"
	"github.com/forbearing/golib/config"
	pkgcontroller "github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/cronjob"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/examples/myproject/model"
	"github.com/forbearing/golib/middleware"
	pkgmodel "github.com/forbearing/golib/model"
	"github.com/forbearing/golib/provider/etcd"
	pkgnats "github.com/forbearing/golib/provider/nats"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/task"
	"github.com/forbearing/golib/types"
	. "github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	_ "github.com/forbearing/golib/examples/myproject/service"
)

var redisCluster = []string{
	"127.0.0.1:6379",
	"127.0.0.1:6380",
	"127.0.0.1:6381",
	"127.0.0.1:6382",
	"127.0.0.1:6383",
	"127.0.0.1:6384",
}

func main() {
	os.Setenv(config.DEBUG_PPROF_ENABLE, "true")
	os.Setenv(config.DEBUG_STATSVIZ_ENABLE, "true")
	os.Setenv(config.DEBUG_GOPS_ENABLE, "true")

	os.Setenv(config.SERVER_PORT, "8002")
	os.Setenv(config.SERVER_DB, config.DBMySQL)
	os.Setenv(config.AUTH_NONE_EXPIRE_TOKEN, "-")
	os.Setenv(config.LOGGER_DIR, "/tmp/myproject/logs")
	os.Setenv(config.DATABASE_MAX_IDLE_CONNS, "100")
	os.Setenv(config.DATABASE_MAX_OPEN_CONNS, "1000")
	os.Setenv(config.SQLITE_PATH, "/tmp/myproject/data.db")

	os.Setenv(config.POSTGRES_ENABLE, "true")
	os.Setenv(config.POSTGRES_PORT, "15432")
	os.Setenv(config.POSTGRES_USERNAME, "postgres")
	os.Setenv(config.POSTGRES_PASSWORD, "admin")

	os.Setenv(config.MYSQL_ENABLE, "true")
	os.Setenv(config.MYSQL_DATABASE, "myproject")
	os.Setenv(config.MYSQL_USERNAME, "myproject")
	os.Setenv(config.MYSQL_PASSWORD, "myproject")

	os.Setenv(config.CLICKHOUSE_ENABLE, "true")
	os.Setenv(config.CLICKHOUSE_DATABASE, "default")
	os.Setenv(config.CLICKHOUSE_USERNAME, "default")
	os.Setenv(config.CLICKHOUSE_PASSWORD, "clickhouse")

	os.Setenv(config.SQLSERVER_ENABLE, "true")
	os.Setenv(config.SQLSERVER_DATABASE, "myproject")
	os.Setenv(config.SQLSERVER_USERNAME, "sa")
	os.Setenv(config.SQLSERVER_PASSWORD, "Passw0rd")

	os.Setenv(config.ELASTICSEARCH_ENABLE, "true")
	os.Setenv(config.ELASTICSEARCH_ENABLE_DEBUG_LOGGER, "true")
	os.Setenv(config.ELASTICSEARCH_ADDRS, "http://127.0.0.1:9200")
	os.Setenv(config.ELASTICSEARCH_USERNAME, "elastic")
	os.Setenv(config.ELASTICSEARCH_PASSWORD, "changeme")

	os.Setenv(config.REDIS_ENABLE, "true")
	os.Setenv(config.REDIS_CLUSTER_MODE, "false")
	os.Setenv(config.REDIS_PASSWORD, "password123")
	os.Setenv(config.REDIS_EXPIRATION, "8h")
	os.Setenv(config.REDIS_ADDRS, strings.Join(redisCluster, ","))
	os.Setenv(config.REDIS_ADDR, "127.0.0.1:6378")

	os.Setenv(config.MONGO_ENABLE, "true")
	os.Setenv(config.MONGO_USERNAME, "mongo")
	os.Setenv(config.MONGO_PASSWORD, "changeme")

	os.Setenv(config.MINIO_ENABLE, "false")
	os.Setenv(config.MINIO_ENDPOINT, "localhost:9000")
	os.Setenv(config.MINIO_ACCESS_KEY, "my_access_key")
	os.Setenv(config.MINIO_SECRET_KEY, "my_secret_key")
	os.Setenv(config.MINIO_REGION, "myregion")
	os.Setenv(config.MINIO_USE_SSL, "false")

	os.Setenv(config.MQTT_ENABLE, "true")
	os.Setenv(config.MQTT_CLIENT_PREFIX, "golib")

	os.Setenv(config.INFLUXDB_ENABLE, "true")
	os.Setenv(config.INFLUXDB_TOKEN, "influxdb")
	os.Setenv(config.INFLUXDB_ORG, "golib.com")

	os.Setenv(config.KAFKA_ENABLE, "true")
	os.Setenv(config.KAFKA_BROKERS, "127.0.0.1:9092,127.0.0.1:9093,127.0.0.1:9094")

	os.Setenv(config.ETCD_ENABLE, "true")
	os.Setenv(config.ETCD_ENDPOINTS, "127.0.0.1:2379,127.0.0.1:12379,127.0.0.1:32379")

	os.Setenv(config.NATS_ADDRS, "nats://127.0.0.1:4222,nats://127.0.0.1:4223,nats://127.0.0.1:4224")
	os.Setenv(config.NATS_ENABLE, "true")

	os.Setenv(config.CASSANDRA_ENABLE, "true")
	os.Setenv(config.CASSANDRA_USERNAME, "cassandra")
	os.Setenv(config.CASSANDRA_PASSWORD, "cassandra")

	os.Setenv(config.LDAP_ENABLE, "true")
	os.Setenv(config.LDAP_PORT, "1389")
	os.Setenv(config.LDAP_BASE_DN, "dc=example,dc=org")
	os.Setenv(config.LDAP_BIND_DN, "cn=admin,dc=example,dc=org")
	os.Setenv(config.LDAP_BIND_PASSWORD, "adminpassword")
	os.Setenv(config.LDAP_USER_DN, "ou=users,dc=example,dc=org")
	os.Setenv(config.LDAP_GROUP_DN, "ou=groups,dc=example,dc=org")
	os.Setenv(config.LDAP_USER_ATTRIBUTE, "uid")
	os.Setenv(config.LDAP_GROUP_ATTRIBUTE, "member")
	os.Setenv(config.LDAP_USER_FILTER, "(objectClass=inetOrgPerson)")
	os.Setenv(config.LDAP_GROUP_FILTER, "(objectClass=groupOfNames)")
	os.Setenv(config.LDAP_SCOPE, "2")
	os.Setenv(config.LDAP_ATTRIBUTES, "uid,cn,sn,mail,memberOf")
	os.Setenv(config.LDAP_DEREF, "0")
	os.Setenv(config.LDAP_PAGE_SIZE, "100")
	os.Setenv(config.LDAP_REQUEST_TIMEOUT, "10s")
	os.Setenv(config.LDAP_CONN_TIMEOUT, "5s")
	os.Setenv(config.LDAP_HEARTBEAT, "30s")

	config.SetConfigFile("./config.ini")
	config.SetConfigName("config")
	config.SetConfigType("ini")

	// Register config before bootstrap.
	config.Register[WechatConfig]("wechat")

	// Register task and cronjob before bootstrap.
	task.Register(SayHello, 1*time.Second, "say hello")
	cronjob.Register(SayHello, "*/1 * * * * *", "say hello")

	//
	//
	//
	RunOrDie(bootstrap.Bootstrap)
	//
	//
	//

	// Register config after bootstrap.
	// config.Register[*NatsConfig]("nats")
	zap.S().Infof("%+v", config.Get[*WechatConfig]("wechat"))
	// zap.S().Infof("%+v", config.Get[*NatsConfig]("nats"))

	// Register task and cronjob after bootstrap.
	task.Register(SayGoodbye, 1*time.Second, "say goodbye")
	cronjob.Register(SayGoodbye, "*/1 * * * * *", "say goodbye")

	// redis
	{
		g1 := &model.Group{Name: "group01"}
		g2 := &model.Group{Name: "group02"}
		groups := []*model.Group{g1, g2}
		if err := redis.SetM("group", g1); err != nil {
			zap.S().Error(err)
		}
		if err := redis.SetML("groups", groups); err != nil {
			zap.S().Error(err)
		}
	}

	// nats
	{
		nc := pkgnats.Conn()
		// 订阅主题
		sub, err := nc.Subscribe("greetings", func(msg *nats.Msg) {
			fmt.Printf("Received: %s\n", string(msg.Data))
			msg.Respond([]byte("Hello back!"))
		})
		if err != nil {
			zap.S().Error(err)
		}
		defer sub.Unsubscribe()

		// 发布消息
		if err = nc.Publish("greetings", []byte("Hello NATS!")); err != nil {
			zap.S().Error(err)
		}

		// 发送请求
		reply, err := nc.Request("greetings", []byte("Hello"), time.Second)
		if err != nil {
			zap.S().Error(err)
		} else {
			fmt.Printf("Reply: %s\n", string(reply.Data))
		}

		time.Sleep(time.Second)
	}
	// etcd
	{
		if _, err := etcd.Client().Put(context.TODO(), "key1", "value1"); err != nil {
			zap.S().Fatal(err)
		}
		fmt.Println("Successfully put key1")
		getResp, err := etcd.Client().Get(context.TODO(), "key1")
		if err != nil {
			zap.S().Fatal(err)
		}
		fmt.Printf("%+v\n", getResp.Kvs)

		if _, err := etcd.Client().Put(context.TODO(), "prefix/key1", "prefixed-value1"); err != nil {
			zap.S().Fatal(err)
		}
		if _, err := etcd.Client().Put(context.TODO(), "prefix/key2", "prefixed-value2"); err != nil {
			zap.S().Fatal(err)
		}
		getResp, err = etcd.Client().Get(context.TODO(), "prefix/", clientv3.WithPrefix())
		if err != nil {
			zap.S().Fatal(err)
		}
		fmt.Printf("%+v\n", getResp.Kvs)

	}

	//
	//
	//
	// router
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

	cfg := config.MySQL{}
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

	RunOrDie(bootstrap.Run)
}

type WechatConfig struct {
	AppID     string `json:"app_id" mapstructure:"app_id" default:"myappid"`
	AppSecret string `json:"app_secret" mapstructure:"app_secret" default:"myappsecret"`
	Enable    bool   `json:"enable" mapstructure:"enable"`
}

// type NatsConfig struct {
// 	URL      string        `json:"url" mapstructure:"url" default:"nats://127.0.0.1:4222"`
// 	Username string        `json:"username" mapstructure:"username" default:"nats"`
// 	Password string        `json:"password" mapstructure:"password" default:"nats"`
// 	Timeout  time.Duration `json:"timeout" mapstructure:"timeout" default:"5s"`
// 	Enable   bool          `json:"enable" mapstructure:"enable"`
// }

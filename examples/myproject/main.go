package main

import (
	"fmt"
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

	_ "github.com/forbearing/golib/examples/myproject/service"
)

func main() {
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
	config.Register[*NatsConfig]("nats")
	fmt.Printf("%+v\n", config.Get[*WechatConfig]("wechat"))
	fmt.Printf("%+v\n", config.Get[NatsConfig]("nats"))

	// Register task and cronjob after bootstrap.
	task.Register(SayGoodbye, 1*time.Second, "say goodbye")
	cronjob.Register(SayGoodbye, "*/1 * * * * *", "say goodbye")

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

type WechatConfig struct {
	AppID     string `json:"app_id" mapstructure:"app_id" default:"myappid"`
	AppSecret string `json:"app_secret" mapstructure:"app_secret" default:"myappsecret"`
	Enable    bool   `json:"enable" mapstructure:"enable"`
}

type NatsConfig struct {
	URL      string        `json:"url" mapstructure:"url" default:"nats://127.0.0.1:4222"`
	Username string        `json:"username" mapstructure:"username" default:"nats"`
	Password string        `json:"password" mapstructure:"password" default:"nats"`
	Timeout  time.Duration `json:"timeout" mapstructure:"timeout" default:"5s"`
	Enable   bool          `json:"enable" mapstructure:"enable"`
}

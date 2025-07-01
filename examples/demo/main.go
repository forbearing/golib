package main

import (
	"fmt"
	"os"
	"time"

	"demo/model"
	model_asset "demo/model/asset"
	model_instance "demo/model/instance"
	_ "demo/service"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/cronjob"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/middleware"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/task"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	. "github.com/forbearing/golib/util"
	"go.uber.org/zap"
)

func main() {
	// Default config is searched in current directory with supported formats(ini,json,yaml,etc.)
	// To specify a custom config file path instead:
	//     config.SetConfigFile("path/to/config.ini")
	// Note: Missing configuration file will not prevent program execution.

	// Set environment variables to override config that read from file and default config value.
	os.Setenv(config.MYSQL_ENABLE, "true")
	os.Setenv(config.MYSQL_DATABASE, "golib")
	os.Setenv(config.MYSQL_USERNAME, "golib")
	os.Setenv(config.MYSQL_PASSWORD, "golib")
	os.Setenv("WECHAT_APP_ID", "id1234567")
	os.Setenv("WECHAT_APP_SECRET", "secret1234567")
	os.Setenv("ALIPAY_APP_ID", "id1234567")
	os.Setenv("ALIPAY_APP_SECRET", "secret1234567")

	// Register cronjob/task before application bootstrap.
	cronjob.Register(func() error { zap.S().Info("cronjob register before bootstrap"); return nil }, "*/1 * * * * *", "cronjob1")
	task.Register(func() error { zap.S().Info("task register before bootstrap"); return nil }, 1*time.Second, "task1")
	// Register custom config before application bootstrap.
	config.Register[WechatConfig]("wechat") // Register accepts either WechatConfig or *WechatConfig as generic type.

	//
	//
	//
	RunOrDie(bootstrap.Bootstrap) // Application Bootstrap.
	//
	//
	//

	// Register cronjob/task after application bootstrap.
	cronjob.Register(func() error { zap.S().Info("cronjob register after bootstrap"); return nil }, "*/1 * * * * *", "cronjob2")
	task.Register(func() error { zap.S().Info("task register after bootstrap"); return nil }, 1*time.Second, "task2")
	// Register custom config after application bootstrap.
	config.Register[*AlipayConfig]("alipay") // Register accepts either AlipayConfig or *AlipayConfig as generic type.

	//
	//
	// Get custom config must after application bootstrap.
	wechat := config.Get[*WechatConfig]("wechat") // Get accepts either WechatConfig or *WechatConfig as generic type.
	alipay := config.Get[AlipayConfig]("alipay")  // Get accepts either AlipayConfig or *AlipayConfig as generic type.

	// WechatConfig.AppId and WechatConfig.AppSecret use values from "default" struct tags
	// These defaults can be overridden by environment variables or config file settings
	fmt.Printf("%+v\n", wechat)
	// Output
	// &{AppId:id1234567 AppSecret:secret1234567}

	// Alipay.AppId and AlipayConfig.AppSecret use values from "default" struct tags
	// These defaults can be overridden by environment variables or config file settings
	fmt.Printf("%+v\n", alipay)
	// Output
	// {AppId:id1234567 AppSecret:secret1234567}
	//
	//
	//

	// use middleware.
	router.API().Use(
		middleware.RateLimiter(),

		// middleware.JwtAuth(),
		// middleware.Gzip(),

		// // setup config.App.AuthConfig.BaseAuthUsername and config.App.AuthConfig.BaseAuthPassword before use this middleware
		// // config.App.Auth.BaseAuthUsername = "admin"
		// // config.App.Auth.BaseAuthPassword = "admin"
		// middleware.BaseAuth(),
	)

	// NOTE: router must register after application bootstrap.
	//
	// Registers basic CURD operation for the model `User`.
	//
	// Generated endpoints:
	//   POST   /api/user     - create user
	//   DELETE /api/user     - delete user
	//   DELETE /api/user/:id - delete user
	//   PUT    /api/user     - update user
	//   PUT    /api/user/:id - update user
	//   PATCH  /api/user     - update partial user
	//   PATCH  /api/user/:id - update partial user
	//   GET    /api/user     - List users
	//   GET    /api/user/:id - Get user by ID
	//
	// Equvalent to:
	// router.Register[*model.User](router.API(), "/user")
	// router.Register[*model.User](router.API(), "user", consts.Most)
	router.Register[*model.User](router.API(), "user")

	// Generated endpoints:
	//   POST    /api/category/batch - batch create category
	//   DELETE  /api/category/batch - batch delete category
	//   PUT     /api/category/batch - batch update category
	//   PATCH   /api/category/batch - batch update partial category
	router.Register[*model.Category](router.API(), "category", consts.MostBatch)

	// Only register `Create` operation for the model `Group`.
	// Generated endpoints:
	//   POST   /api/group     - create group
	//   POST   /api/group/:id - create group
	router.Register[*model.Group](router.API(), "group", consts.Create)

	// Only register `Delete` operation for the model `Group`.
	// Generated endpoints:
	//   DELETE /api/group     - delete group
	//   DELETE /api/group/:id - delete group
	router.Register[*model.Group](router.API(), "group", consts.Delete)

	// Only register `Update` operation for the model `Group`.
	// Generated endpoints:
	//   PUT    /api/group     - update group
	//   PUT    /api/group/:id - update group
	router.Register[*model.Group](router.API(), "group", consts.Update)

	// Only register `Update` operation for the model `Group`.
	// Generated endpoints:
	//   PATCH  /api/group     - update partial group
	//   PATCH  /api/group/:id - update partial group
	router.Register[*model.Group](router.API(), "group", consts.UpdatePartial)

	// Only register `List` operation for the model `Group`.
	// Generated endpoints:
	//   GET    /api/group/    - List users
	router.Register[*model.Group](router.API(), "group", consts.List)

	// Only register `Get` operation for the model `Group`.
	// Generated endpoints:
	//   GET    /api/group/:id - Get group by ID
	router.Register[*model.Group](router.API(), "group", consts.Get)

	// Only register `BatchCreate` operation for the model `Group`.
	// Generated endpoints:
	//   POST   /api/group/batch - batch create group
	router.Register[*model.Group](router.API(), "group", consts.BatchCreate)

	// Only register `BatchDelete` operation for the model `Group`.
	// Generated endpoints:
	//   DELETE /api/group/batch - batch delete group
	router.Register[*model.Group](router.API(), "group", consts.BatchDelete)

	// Only register `BatchUpdate` operation for the model `Group`.
	// Generated endpoints:
	//   PUT    /api/group/batch - batch update group
	router.Register[*model.Group](router.API(), "group", consts.BatchUpdate)

	// Only register `BatchUpdatePartial` operation for the model `Group`.
	// Generated endpoints:
	//   PATCH  /api/group/batch - batch update partial group
	router.Register[*model.Group](router.API(), "group", consts.BatchUpdatePartial)

	// Generated endpoints:
	//   POST   /api/group/import - import groups.
	router.Register[*model.Group](router.API(), "group", consts.Import)

	// Generated endpoints:
	//   GET    /api/group/export - export groups.
	router.Register[*model.Group](router.API(), "group", consts.Export)

	// Manual RESTful API route configuration for model `Contact`.
	router.API().POST("/contact", controller.Create[*model.Contact])                    // create
	router.API().DELETE("/contact", controller.Delete[*model.Contact])                  // delete
	router.API().DELETE("/contact/:id", controller.Delete[*model.Contact])              // delete
	router.API().PUT("/contact", controller.Update[*model.Contact])                     // update
	router.API().PUT("/contact/:id", controller.Update[*model.Contact])                 // update
	router.API().PATCH("/contact", controller.UpdatePartial[*model.Contact])            // update partial
	router.API().PATCH("/contact/:id", controller.UpdatePartial[*model.Contact])        // update partial
	router.API().GET("/contact", controller.List[*model.Contact])                       // list
	router.API().GET("/contact/:id", controller.Get[*model.Contact])                    // get
	router.API().POST("/contact/batch", controller.BatchCreate[*model.Contact])         // batch create
	router.API().DELETE("/contact/batch", controller.BatchDelete[*model.Contact])       // batch delete
	router.API().PUT("/contact/batch", controller.BatchUpdate[*model.Contact])          // batch update
	router.API().PATCH("/contact/batch", controller.BatchUpdatePartial[*model.Contact]) // batch update partial

	// Models are defined in separate packages:
	// - model.*          -> core models
	// - model_asset.*    -> asset management models
	// - model_instance.* -> instance management models
	//
	// Register core model routes.
	router.Register[*model.Order](router.API(), "order", consts.Most) // route: /api/order
	router.Register[*model.Log](router.API(), "log", consts.Most)     // route: /api/log
	// Register asset routes.
	asset := router.API().Group("/asset")
	router.Register[*model_asset.Computer](asset, "computer", consts.Most) // route: /api/asset/computer
	router.Register[*model_asset.Monitor](asset, "monitor")                // route: /api/asset/monitor
	router.Register[*model_asset.Software](asset, "software")              // route: /api/asset/software
	router.Register[*model_asset.Printer](asset, "printer")                // route: /api/asset/printer
	// Register instance routes.
	instance := router.API().Group("/instance")
	router.Register[*model_instance.Datacenter](instance, "datacenter")   // route: /api/instance/datacenter
	router.Register[*model_instance.Cluster](instance, "cluster")         // route: /api/instance/cluster
	router.Register[*model_instance.Database](instance, "database")       // route: /api/instance/database
	router.Register[*model_instance.Certificate](instance, "certificate") // route: /api/instance/certificate

	// With seperate mysql instance
	cfg := config.MySQL{}
	cfg.Host = "127.0.0.1"
	cfg.Port = 3306
	cfg.Database = "test"
	cfg.Username = "test"
	cfg.Password = "test"
	cfg.Charset = "utf8mb4"
	db, err := mysql.New(cfg)
	if err != nil {
		panic(err)
	}
	// NOTE: It is your responsibility to ensure the table that map to model already exists.
	external := router.API().Group("/external")
	router.RegisterWithConfig(external, "/user", &types.ControllerConfig[*model.User]{DB: db, TableName: "users"})
	router.RegisterWithConfig(external, "category", &types.ControllerConfig[*model.Category]{DB: db, TableName: "groups"})

	RunOrDie(bootstrap.Run)
}

type WechatConfig struct {
	AppId     string `json:"app_id" yaml:"app_id" mapstructure:"app_id" ini:"app_id" default:"myid"`
	AppSecret string `json:"app_secret" yaml:"app_secret" mapstructure:"app_secret" ini:"app_secret" default:"mysecret"`
}
type AlipayConfig struct {
	AppId     string `json:"app_id" yaml:"app_id" mapstructure:"app_id" ini:"app_id" default:"myid"`
	AppSecret string `json:"app_secret" yaml:"app_secret" mapstructure:"app_secret" ini:"app_secret" default:"mysecret"`
}

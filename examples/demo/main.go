package main

import (
	"demo/model"
	model_asset "demo/model/asset"
	model_instance "demo/model/instance"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/middleware"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
)

func main() {
	// Ensure config file `config.init` exists in the current directory, empty config is permitted.
	// Alternatively, specify a custom config file path:
	//     config.SetConfigFile("path/to/config.ini")

	util.RunOrDie(bootstrap.Bootstrap)

	// set base auth config.
	config.App.AuthConfig.BaseAuthUsername = "admin"
	config.App.AuthConfig.BaseAuthPassword = "admin"

	// use middleware.
	router.API.Use(
		middleware.RateLimiter(),

		// middleware.JwtAuth(),
		// middleware.Gzip(),

		// // setup config.App.AuthConfig.BaseAuthUsername and config.App.AuthConfig.BaseAuthPassword before use this middleware
		// middleware.BaseAuth(),
	)

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
	// router.Register[*model.User](router.API, "/user")
	// router.Register[*model.User](router.API, "user", types.Most)
	router.Register[*model.User](router.API, "user")

	// Only register `Create` operation for the model `Group`.
	// Generated endpoints:
	//   POST   /api/group     - create group
	//   POST   /api/group/:id - create group
	router.Register[*model.Group](router.API, "group", types.Create)
	// Only register `Delete` operation for the model `Group`.
	// Generated endpoints:
	//   DELETE /api/group     - delete group
	//   DELETE /api/group/:id - delete group
	router.Register[*model.Group](router.API, "group", types.Delete)
	// Only register `Update` operation for the model `Group`.
	// Generated endpoints:
	//   PUT    /api/group     - update group
	//   PUT    /api/group/:id - update group
	router.Register[*model.Group](router.API, "group", types.Update)
	// Only register `Update` operation for the model `Group`.
	// Generated endpoints:
	//   PATCH  /api/group     - update partial group
	//   PATCH  /api/group/:id - update partial group
	router.Register[*model.Group](router.API, "group", types.UpdatePartial)
	// Only register `List` operation for the model `Group`.
	// Generated endpoints:
	//   GET    /api/group/    - List users
	router.Register[*model.Group](router.API, "group", types.List)
	// Only register `Get` operation for the model `Group`.
	// Generated endpoints:
	//   GET    /api/group/:id - Get group by ID
	router.Register[*model.Group](router.API, "group", types.Get)
	// Generated endpoints:
	//   POST   /api/group/import - import groups.
	router.Register[*model.Group](router.API, "group", types.Import)
	// Generated endpoints:
	//   GET    /api/group/export - export groups.
	router.Register[*model.Group](router.API, "group", types.Export)

	// `All` equvalent to: `Most` + `Import` + `Export`
	router.Register[*model.Department](router.API, "department", types.All)

	// Manual RESTful API route configuration for model `Contact`.
	router.API.POST("/contact", controller.Create[*model.Contact])             // create
	router.API.DELETE("/contact", controller.Delete[*model.Contact])           // delete
	router.API.DELETE("/contact/:id", controller.Delete[*model.Contact])       // delete
	router.API.PUT("/contact", controller.Update[*model.Contact])              // update
	router.API.PUT("/contact/:id", controller.Update[*model.Contact])          // update
	router.API.PATCH("/contact", controller.UpdatePartial[*model.Contact])     // update partial
	router.API.PATCH("/contact/:id", controller.UpdatePartial[*model.Contact]) // update partial
	router.API.GET("/contact", controller.List[*model.Contact])                // list
	router.API.GET("/contact/:id", controller.Get[*model.Contact])             // get

	// Models are defined in separate packages:
	// - model.*          -> core models
	// - model_asset.*    -> asset management models
	// - model_instance.* -> instance management models
	//
	// Register core model routes.
	router.Register[*model.Order](router.API, "order") // route: /api/order
	router.Register[*model.Log](router.API, "log")     // route: /api/log
	// Register asset routes.
	asset := router.API.Group("/asset")
	router.Register[*model_asset.Computer](asset, "computer") // route: /api/asset/computer
	router.Register[*model_asset.Monitor](asset, "monitor")   // route: /api/asset/monitor
	router.Register[*model_asset.Software](asset, "software") // route: /api/asset/software
	router.Register[*model_asset.Printer](asset, "printer")   // route: /api/asset/printer
	// Register instance routes.
	instance := router.API.Group("/instance")
	router.Register[*model_instance.Datacenter](instance, "datacenter")   // route: /api/instance/datacenter
	router.Register[*model_instance.Cluster](instance, "cluster")         // route: /api/instance/cluster
	router.Register[*model_instance.Database](instance, "database")       // route: /api/instance/database
	router.Register[*model_instance.Certificate](instance, "certificate") // route: /api/instance/certificate

	util.RunOrDie(router.Run)
}

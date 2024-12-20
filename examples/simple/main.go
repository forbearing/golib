package main

import (
	"net/http"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
)

func main() {
	// Verify existence of `config.ini` in the current directory. Empty configuration is permitted.
	// To specify an alternative config file, use config.SetConfigPath, e.g.:
	//
	// config.SetConfigFile("path/to/config.ini")

	/*
		// more config usage:
		// SetConfigFile("config.ini") is equivalent to SetConfigName("config") + SetConfigType("ini")
		config.SetConfigName("config")      // set the config name.
		config.SetConfigType("ini")         // set the config type.
		config.AddPath("/etc", "/", "/tmp") // add config search path.
	*/

	// Bootstrap all initializers
	util.RunOrDie(bootstrap.Bootstrap)

	// Set up your routes here, the routes are not required.
	router.API.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	router.API.GET("/hello", func(c *gin.Context) { c.String(http.StatusOK, "hello world!") })

	// Any router panic, exit, or error will cause program termination.
	util.RunOrDie(router.Run)
}

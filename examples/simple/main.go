package main

import (
	"net/http"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
)

func main() {
	// Ensure config file `config.ini` exists in the current path.
	// To use a different config file, call config.SetConfigPath to specify the path, e.g.:
	//
	// config.SetConfigFile("config.ini")

	// Bootstrap all initializers
	util.RunOrDie(bootstrap.Bootstrap)

	// Set up your routes here, the routes are not required.
	router.API.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	router.API.GET("/hello", func(c *gin.Context) { c.String(http.StatusOK, "hello world!") })

	// Any router panic, exit, or error will cause program termination.
	util.RunOrDie(router.Run)
}

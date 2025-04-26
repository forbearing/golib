package main

import (
	"net/http"

	"demo/model"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
)

func main() {
	// Bootstrap all initializers
	util.RunOrDie(bootstrap.Bootstrap)

	// Set up your routes here, the routes are not required.
	router.API().GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	router.API().GET("/hello", func(c *gin.Context) { c.String(http.StatusOK, "hello world!") })
	router.Register[*model.User](router.API(), "user")

	util.RunOrDie(bootstrap.Run)
}

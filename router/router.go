package router

import (
	"net"
	"net/http"
	"strconv"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	router *gin.Engine
	API    *gin.RouterGroup
)

func Init() error {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	router.Use(
		middleware.RequestID(),
		middleware.Logger("./logs/api.log"),
		middleware.Recovery("./logs/recovery.log"),
		middleware.Cors(),
		middleware.RateLimiter(),
	)
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	API = router.Group("/api")
	API.Use(
		middleware.JwtAuth(),
		// middleware.Authz(),
		middleware.Gzip(),
	)
	return nil
}

func Run() error {
	addr := net.JoinHostPort(config.App.ServerConfig.Listen, strconv.Itoa(config.App.ServerConfig.Port))
	zap.S().Infow("starting server", "addr", addr, "mode", config.App.Mode, "domain", config.App.Domain)
	// for _, r := range router.Routes() {
	// 	// zap.S().Infof("%v %v", r.Method, r.Path)
	// 	zap.S().Infow("", "method", r.Method, "path", r.Path)
	// }
	return router.Run(addr)
}

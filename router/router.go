package router

import (
	"net"
	"net/http"
	"strconv"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	Base *gin.Engine
	API  *gin.RouterGroup
)

func Init() error {
	gin.SetMode(gin.ReleaseMode)
	Base = gin.New()

	Base.Use(
		middleware.RequestID(),
		middleware.Logger("api.log"),
		middleware.Recovery("recovery.log"),
		middleware.Cors(),
		// middleware.RateLimiter(),
	)
	Base.GET("/ping", func(ctx *gin.Context) { ctx.String(http.StatusOK, "pong") })
	Base.GET("/metrics", gin.WrapH(promhttp.Handler()))
	Base.GET("/-/healthz", controller.Probe.Healthz)
	Base.GET("/-/readyz", controller.Probe.Readyz)
	Base.GET("/-/pageid", controller.PageID)

	API = Base.Group("/api")
	API.Use(
		// middleware.JwtAuth(),
		// middleware.Authz(),
		middleware.Gzip(),
	)
	return nil
}

func Run() error {
	addr := net.JoinHostPort(config.App.ServerConfig.Listen, strconv.Itoa(config.App.ServerConfig.Port))
	zap.S().Infow("starting server", "addr", addr, "mode", config.App.Mode, "domain", config.App.Domain)
	for _, r := range Base.Routes() {
		zap.S().Debugw("", "method", r.Method, "path", r.Path)
	}
	return Base.Run(addr)
}

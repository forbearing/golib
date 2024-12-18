package router

import (
	"net"
	"strconv"
	"strings"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/middleware"
	"github.com/forbearing/golib/types"
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

// Register registers HTTP routes for a given model type with specified verbs.
// It supports common CRUD operations along with import/export functionality.
//
// Parameters:
//   - router: The Gin router instance to register routes on
//   - path: Base path for the resource (automatically handles '/api/' prefix)
//   - verbs: Optional list of HTTPVerb to register. If empty, defaults to Most (basic CRUD operations)
//
// Route patterns registered:
//   - POST   /{path}         -> Create
//   - DELETE /{path}         -> Delete (bulk)
//   - DELETE /{path}/:id     -> Delete (single)
//   - PUT    /{path}         -> Update (bulk)
//   - PUT    /{path}/:id     -> Update (single)
//   - PATCH  /{path}         -> UpdatePartial (bulk)
//   - PATCH  /{path}/:id     -> UpdatePartial (single)
//   - GET    /{path}         -> List
//   - GET    /{path}/:id     -> Get
//   - POST   /{path}/import  -> Import
//   - GET    /{path}/export  -> Export
func Register[M types.Model](router gin.IRouter, path string, verbs ...types.HTTPVerb) {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		zap.S().Warn("empty path, skip register routes")
		return
	}
	path = strings.TrimPrefix(path, `/api/`) // remove path prefix: '/api/'
	path = strings.TrimPrefix(path, "/")     // trim left "/"
	path = strings.TrimSuffix(path, "/")     // trim right "/"
	path = "/" + path

	verbMap := make(map[types.HTTPVerb]bool)

	if len(verbs) == 0 {
		verbMap[types.Most] = true
	} else {
		for _, verb := range verbs {
			verbMap[verb] = true
		}
	}
	if verbMap[types.All] {
		verbMap[types.Most] = true
		verbMap[types.Import] = true
		verbMap[types.Export] = true
	}
	if verbMap[types.Most] {
		verbMap[types.Create] = true
		verbMap[types.Delete] = true
		verbMap[types.Update] = true
		verbMap[types.UpdatePartial] = true
		verbMap[types.List] = true
		verbMap[types.Get] = true
	}

	if verbMap[types.Create] {
		router.POST(path, controller.Create[M])
	}
	if verbMap[types.Delete] {
		router.DELETE(path, controller.Delete[M])
		router.DELETE(path+"/:id", controller.Delete[M])
	}
	if verbMap[types.Update] {
		router.PUT(path, controller.Update[M])
		router.PUT(path+"/:id", controller.Update[M])
	}
	if verbMap[types.UpdatePartial] {
		router.PATCH(path, controller.UpdatePartial[M])
		router.PATCH(path+"/:id", controller.UpdatePartial[M])
	}
	if verbMap[types.List] {
		router.GET(path, controller.List[M])
	}
	if verbMap[types.Get] {
		router.GET(path+"/:id", controller.Get[M])
	}
	if verbMap[types.Import] {
		router.POST(path+"/import", controller.Import[M])
	}
	if verbMap[types.Export] {
		router.GET(path+"/export", controller.Export[M])
	}
}

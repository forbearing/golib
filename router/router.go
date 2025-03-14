package router

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/middleware"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	Base *gin.Engine
	API  *gin.RouterGroup

	server *http.Server
)

func Init() error {
	gin.SetMode(gin.ReleaseMode)
	Base = gin.New()

	Base.Use(
		middleware.TraceID(),
		middleware.Logger("api.log"),
		middleware.Recovery("recovery.log"),
		middleware.Cors(),
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
	log := zap.S()
	addr := net.JoinHostPort(config.App.Server.Listen, strconv.Itoa(config.App.Server.Port))
	log.Infow("backend server started", "addr", addr, "mode", config.App.Mode, "domain", config.App.Domain)
	for _, r := range Base.Routes() {
		log.Debugw("", "method", r.Method, "path", r.Path)
	}

	server = &http.Server{
		Addr:           addr,
		Handler:        Base,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Errorw("failed to start server", "err", err)
		return err
	}
	return nil
}

func Stop() {
	if server == nil {
		return
	}
	zap.S().Infow("backend server shutdown initiated")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zap.S().Errorw("backend server shutdown failed", "err", err)
	} else {
		zap.S().Infow("backend server shutdown completed")
	}
	server = nil
}

// Register registers HTTP routes for a given model type with specified verbs
// using default controller configuration.
// It supports common CRUD operations along with import/export functionality.
//
// Parameters:
//   - router: The Gin router instance to register routes on
//   - path: Base path for the resource (automatically handles '/api/' prefix)
//   - verbs: Optional list of HTTPVerb to register. If empty, defaults to Most (basic CRUD operations)
//
// Route patterns registered:
//   - POST   /{path}         -> Create
//   - DELETE /{path}         -> Delete
//   - DELETE /{path}/:id     -> Delete
//   - PUT    /{path}         -> Update
//   - PUT    /{path}/:id     -> Update
//   - PATCH  /{path}         -> UpdatePartial
//   - PATCH  /{path}/:id     -> UpdatePartial
//   - GET    /{path}         -> List
//   - GET    /{path}/:id     -> Get
//   - POST   /{path}/import  -> Import
//   - GET    /{path}/export  -> Export
//
// For custom controller configuration, use RegisterWithConfig instead.
func Register[M types.Model](router gin.IRouter, rawPath string, verbs ...consts.HTTPVerb) {
	rawPath = strings.TrimSpace(rawPath)
	if len(rawPath) == 0 {
		zap.S().Warn("empty path, skip register routes")
		return
	}
	path := buildPath(rawPath)
	verbMap := buildVerbMap(verbs...)

	register[M](router, path, verbMap)
}

// RegisterWithConfig is same as Register, but with custom controller configuration.
// The cfg parameter allow custom controller behavior.
// more details see Register.
func RegisterWithConfig[M types.Model](cfg *types.ControllerConfig[M], router gin.IRouter, rawPath string, verbs ...consts.HTTPVerb) {
	rawPath = strings.TrimSpace(rawPath)
	if len(rawPath) == 0 {
		zap.S().Warn("empty path, skip register routes")
		return
	}
	path := buildPath(rawPath)
	verbMap := buildVerbMap(verbs...)

	register(router, path, verbMap, cfg)
}

func register[M types.Model](router gin.IRouter, path string, verbMap map[consts.HTTPVerb]bool, cfg ...*types.ControllerConfig[M]) {
	if verbMap[consts.Create] {
		router.POST(path, controller.CreateFactory(cfg...))
	}
	if verbMap[consts.Delete] {
		router.DELETE(path, controller.DeleteFactory(cfg...))
		router.DELETE(path+"/:id", controller.DeleteFactory(cfg...))
	}
	if verbMap[consts.Update] {
		router.PUT(path, controller.UpdateFactory(cfg...))
		router.PUT(path+"/:id", controller.UpdateFactory(cfg...))
	}
	if verbMap[consts.UpdatePartial] {
		router.PATCH(path, controller.UpdatePartialFactory(cfg...))
		router.PATCH(path+"/:id", controller.UpdatePartialFactory(cfg...))
	}
	if verbMap[consts.List] {
		router.GET(path, controller.ListFactory(cfg...))
	}
	if verbMap[consts.Get] {
		router.GET(path+"/:id", controller.GetFactory(cfg...))
	}

	if verbMap[consts.BatchCreate] {
		router.POST(path+"/batch", controller.BatchCreateFactory(cfg...))
	}
	if verbMap[consts.BatchDelete] {
		router.DELETE(path+"/batch", controller.BatchDeleteFactory(cfg...))
	}
	if verbMap[consts.BatchUpdate] {
		router.PUT(path+"/batch", controller.BatchUpdateFactory(cfg...))
	}
	if verbMap[consts.BatchUpdatePartial] {
		router.PATCH(path+"/batch", controller.BatchUpdatePartialFactory(cfg...))
	}

	if verbMap[consts.Import] {
		router.POST(path+"/import", controller.ImportFactory(cfg...))
	}
	if verbMap[consts.Export] {
		router.GET(path+"/export", controller.ExportFactory(cfg...))
	}
}

// buildPath normalizes the API path.
func buildPath(path string) string {
	path = strings.TrimPrefix(path, `/api/`) // remove path prefix: '/api/'
	path = strings.TrimPrefix(path, "/")     // trim left "/"
	path = strings.TrimSuffix(path, "/")     // trim right "/"
	return "/" + path
}

// buildVerbMap creates a map of allowed HTTP verbs according to the specified verbs.
func buildVerbMap(verbs ...consts.HTTPVerb) map[consts.HTTPVerb]bool {
	verbMap := make(map[consts.HTTPVerb]bool)

	if len(verbs) == 0 {
		verbMap[consts.Most] = true
	} else {
		for _, verb := range verbs {
			verbMap[verb] = true
		}
	}
	if verbMap[consts.All] {
		verbMap[consts.Most] = true
		verbMap[consts.Import] = true
		verbMap[consts.Export] = true
	}
	if verbMap[consts.Most] {
		verbMap[consts.Create] = true
		verbMap[consts.Delete] = true
		verbMap[consts.Update] = true
		verbMap[consts.UpdatePartial] = true
		verbMap[consts.List] = true
		verbMap[consts.Get] = true
		verbMap[consts.BatchCreate] = true
		verbMap[consts.BatchDelete] = true
		verbMap[consts.BatchUpdate] = true
		verbMap[consts.BatchUpdatePartial] = true
	}
	return verbMap
}

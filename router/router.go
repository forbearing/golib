package router

import (
	"context"
	"net"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/internal/openapigen"
	"github.com/forbearing/golib/middleware"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

var (
	base *gin.Engine
	api  *gin.RouterGroup

	server *http.Server
)

var globalErrors = make([]error, 0)

func Init() error {
	gin.SetMode(gin.ReleaseMode)
	base = gin.New()

	base.Use(
		middleware.TraceID(),
		middleware.Logger("api.log"),
		middleware.Recovery("recovery.log"),
		middleware.Cors(),
		middleware.RouteParams(),
	)
	base.GET("/metrics", gin.WrapH(promhttp.Handler()))
	base.GET("/-/healthz", controller.Probe.Healthz)
	base.GET("/-/readyz", controller.Probe.Readyz)
	base.GET("/-/pageid", controller.PageID)
	base.GET("/api.json", middleware.BaseAuth(), gin.WrapH(openapigen.DocumentHandler()))
	base.GET("/api/docs/*any", middleware.BaseAuth(), ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api.json")))

	api = base.Group("/api")
	api.Use(
		middleware.Gzip(),
		// middleware.JwtAuth(),
		// middleware.Authz(),
	)
	return nil
}

func Run() error {
	log := zap.S()
	if err := multierr.Combine(globalErrors...); err != nil {
		log.Error(err)
		return err
	}

	addr := net.JoinHostPort(config.App.Server.Listen, strconv.Itoa(config.App.Server.Port))
	log.Infow("backend server started", "addr", addr, "mode", config.App.Mode, "domain", config.App.Domain)
	for _, r := range base.Routes() {
		log.Debugw("", "method", r.Method, "path", r.Path)
	}

	server = &http.Server{
		Addr:           addr,
		Handler:        base,
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

func API() *gin.RouterGroup { return api }
func Base() *gin.Engine     { return base }

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
//   - POST   /{path}/batch   -> BatchCreate
//   - DELETE /{path}/batch   -> BatchDelete
//   - PUT    /{path}/batch   -> BatchUpdate
//   - PATCH  /{path}/batch   -> BatchUpdatePartial
//
// For custom controller configuration, pass a ControllerConfig object.
func Register[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.Most), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.Most)
	}
}

func RegisterCreate[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.Create), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.Create)
	}
}

func RegisterDelete[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.Delete), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.Delete)
	}
}

func RegisterUpdate[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.Update), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.Update)
	}
}

func RegisterUpdatePartial[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.UpdatePartial), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.UpdatePartial)
	}
}

func RegisterList[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.List), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.List)
	}
}

func RegisterGet[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.Get), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.Get)
	}
}

func RegisterImport[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.Import), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.Import)
	}
}

func RegisterExport[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.Export), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.Export)
	}
}

func RegisterBatchCreate[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.BatchCreate), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.BatchCreate)
	}
}

func RegisterBatchDelete[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.BatchDelete), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.BatchDelete)
	}
}

func RegisterBatchUpdate[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.BatchUpdate), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.BatchUpdate)
	}
}

func RegisterBatchUpdatePartial[M types.Model](router gin.IRouter, rawPath string, cfg ...*types.ControllerConfig[M]) {
	if validPath(rawPath) {
		register(router, buildPath(rawPath), buildVerbMap(consts.BatchUpdatePartial), cfg...)
		openapigen.Set[M](buildPath(rawPath), consts.BatchUpdatePartial)
	}
}

func validPath(rawPath string) bool {
	rawPath = strings.TrimSpace(rawPath)
	if len(rawPath) == 0 {
		zap.S().Warn("empty path, skip register routes")
		return false
	}
	return true
}

func register[M types.Model](router gin.IRouter, path string, verbMap map[consts.HTTPVerb]bool, cfg ...*types.ControllerConfig[M]) {
	v := reflect.ValueOf(router).Elem()
	base := v.FieldByName("basePath").String()
	// AutoMigrate ensures the table structure exists in the database.
	// This automatically creates or updates the table based on the model structure.
	// Alternatively, you can use model.Register() to manually control table creation.
	m := reflect.New(reflect.TypeOf(*new(M)).Elem()).Interface().(M)
	if err := database.DB.Table(m.GetTableName()).AutoMigrate(m); err != nil {
		globalErrors = append(globalErrors, err)
	}

	if verbMap[consts.Create] {
		router.POST(path, controller.CreateFactory(cfg...))
		middleware.RouteManager.Add(filepath.Join(base, path))
	}
	if verbMap[consts.Delete] {
		router.DELETE(path, controller.DeleteFactory(cfg...))
		router.DELETE(path+"/:id", controller.DeleteFactory(cfg...))
		middleware.RouteManager.Add(filepath.Join(base, path))
		middleware.RouteManager.Add(filepath.Join(base, path+"/:id"))
	}
	if verbMap[consts.Update] {
		router.PUT(path, controller.UpdateFactory(cfg...))
		router.PUT(path+"/:id", controller.UpdateFactory(cfg...))
		middleware.RouteManager.Add(filepath.Join(base, path))
		middleware.RouteManager.Add(filepath.Join(base, path+"/:id"))
	}
	if verbMap[consts.UpdatePartial] {
		router.PATCH(path, controller.UpdatePartialFactory(cfg...))
		router.PATCH(path+"/:id", controller.UpdatePartialFactory(cfg...))
		middleware.RouteManager.Add(filepath.Join(base, path))
		middleware.RouteManager.Add(filepath.Join(base, path+"/:id"))
	}
	if verbMap[consts.List] {
		router.GET(path, controller.ListFactory(cfg...))
		middleware.RouteManager.Add(filepath.Join(base, path))
	}
	if verbMap[consts.Get] {
		router.GET(path+"/:id", controller.GetFactory(cfg...))
		middleware.RouteManager.Add(filepath.Join(base, path, "/:id"))
	}

	if verbMap[consts.BatchCreate] {
		router.POST(path+"/batch", controller.BatchCreateFactory(cfg...))
		middleware.RouteManager.Add(filepath.Join(base, path, "/batch"))
	}
	if verbMap[consts.BatchDelete] {
		router.DELETE(path+"/batch", controller.BatchDeleteFactory(cfg...))
		middleware.RouteManager.Add(path + "/batch")
	}
	if verbMap[consts.BatchUpdate] {
		router.PUT(path+"/batch", controller.BatchUpdateFactory(cfg...))
		middleware.RouteManager.Add(path + "/batch")
	}
	if verbMap[consts.BatchUpdatePartial] {
		router.PATCH(path+"/batch", controller.BatchUpdatePartialFactory(cfg...))
		middleware.RouteManager.Add(path + "/batch")
	}

	if verbMap[consts.Import] {
		router.POST(path+"/import", controller.ImportFactory(cfg...))
		middleware.RouteManager.Add(path + "/import")
	}
	if verbMap[consts.Export] {
		router.GET(path+"/export", controller.ExportFactory(cfg...))
		middleware.RouteManager.Add(path + "/export")
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

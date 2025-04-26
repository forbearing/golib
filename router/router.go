package router

import (
	"context"
	"net"
	"net/http"
	gopath "path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/controller"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/internal/openapigen"
	"github.com/forbearing/golib/middleware"
	"github.com/forbearing/golib/model"
	model_authz "github.com/forbearing/golib/model/authz"
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
	base.GET("/-/api.json", middleware.BaseAuth(), gin.WrapH(openapigen.DocumentHandler()))
	base.GET("/-/api/docs/*any", middleware.BaseAuth(), ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/-/api.json")))

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

	// Delete all permissions in database
	permissions := make([]*model_authz.Permission, 0)
	if err := database.Database[*model_authz.Permission]().WithLimit(-1).List(&permissions); err != nil {
		log.Error(err)
		return err
	}
	if err := database.Database[*model_authz.Permission]().WithLimit(-1).WithBatchSize(100).WithPurge().Delete(permissions...); err != nil {
		log.Error(err)
		return err
	}
	// Create permissions in database
	permissions = make([]*model_authz.Permission, 0)
	for endpoint, methods := range model.Routes {
		for _, method := range methods {
			permissions = append(permissions, &model_authz.Permission{
				Resource: convertGinPathToCasbinKeyMatch3(endpoint),
				Action:   method,
			})
		}
	}
	if err := database.Database[*model_authz.Permission]().WithLimit(-1).WithBatchSize(100).Create(permissions...); err != nil {
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
		ReadTimeout:    config.App.Server.ReadTimeout,
		WriteTimeout:   config.App.Server.WriteTimeout,
		IdleTimeout:    config.App.IdleTimeout,
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
	// v := reflect.ValueOf(router).Elem()
	// base := v.FieldByName("basePath").String()
	var base string
	if group, ok := router.(*gin.RouterGroup); ok {
		base = group.BasePath()
	} else {
		panic("unknown router type")
	}
	// AutoMigrate ensures the table structure exists in the database.
	// This automatically creates or updates the table based on the model structure.
	// Alternatively, you can use model.Register() to manually control table creation.
	if model.IsValid[M]() {
		m := reflect.New(reflect.TypeOf(*new(M)).Elem()).Interface().(M)
		if err := database.DB.Table(m.GetTableName()).AutoMigrate(m); err != nil {
			globalErrors = append(globalErrors, err)
		}
	}

	if verbMap[consts.Create] {
		endpoint := gopath.Join(base, path)
		router.POST(path, controller.CreateFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodPost)
		middleware.RouteManager.Add(endpoint)
	}
	if verbMap[consts.Delete] {
		endpoint := gopath.Join(base, path)
		endpoint2 := gopath.Join(base, path, "/:id")
		router.DELETE(path, controller.DeleteFactory(cfg...))
		router.DELETE(path+"/:id", controller.DeleteFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodDelete)
		model.Routes[endpoint2] = append(model.Routes[endpoint2], http.MethodDelete)
		middleware.RouteManager.Add(endpoint)
		middleware.RouteManager.Add(endpoint2)

	}
	if verbMap[consts.Update] {
		endpoint := gopath.Join(base, path)
		endpoint2 := gopath.Join(base, path, "/:id")
		router.PUT(path, controller.UpdateFactory(cfg...))
		router.PUT(path+"/:id", controller.UpdateFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodPut)
		model.Routes[endpoint2] = append(model.Routes[endpoint2], http.MethodPut)
		middleware.RouteManager.Add(endpoint)
		middleware.RouteManager.Add(endpoint2)

	}
	if verbMap[consts.UpdatePartial] {
		endpoint := gopath.Join(base, path)
		endpoint2 := gopath.Join(base, path, "/:id")
		router.PATCH(path, controller.UpdatePartialFactory(cfg...))
		router.PATCH(path+"/:id", controller.UpdatePartialFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodPatch)
		model.Routes[endpoint2] = append(model.Routes[endpoint2], http.MethodPatch)
		middleware.RouteManager.Add(endpoint)
		middleware.RouteManager.Add(endpoint2)
	}
	if verbMap[consts.List] {
		endpoint := gopath.Join(base, path)
		router.GET(path, controller.ListFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodGet)
		middleware.RouteManager.Add(endpoint)
	}
	if verbMap[consts.Get] {
		endpoint := gopath.Join(base, path, "/:id")
		router.GET(path+"/:id", controller.GetFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodGet)
		middleware.RouteManager.Add(endpoint)
	}

	if verbMap[consts.BatchCreate] {
		endpoint := gopath.Join(base, path, "/batch")
		router.POST(path+"/batch", controller.BatchCreateFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodPost)
		middleware.RouteManager.Add(gopath.Join(base, path, "/batch"))
	}
	if verbMap[consts.BatchDelete] {
		endpoint := gopath.Join(base, path, "/batch")
		router.DELETE(path+"/batch", controller.BatchDeleteFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodDelete)
		middleware.RouteManager.Add(endpoint)
	}
	if verbMap[consts.BatchUpdate] {
		endpoint := gopath.Join(base, path, "/batch")
		router.PUT(path+"/batch", controller.BatchUpdateFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodPut)
		middleware.RouteManager.Add(endpoint)
	}
	if verbMap[consts.BatchUpdatePartial] {
		endpoint := gopath.Join(base, path, "/batch")
		router.PATCH(path+"/batch", controller.BatchUpdatePartialFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodPatch)
		middleware.RouteManager.Add(endpoint)
	}

	if verbMap[consts.Import] {
		endpoint := gopath.Join(base, path, "/import")
		router.POST(path+"/import", controller.ImportFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodPost)
		middleware.RouteManager.Add(endpoint)
	}
	if verbMap[consts.Export] {
		endpoint := gopath.Join(base, path, "/export")
		router.GET(path+"/export", controller.ExportFactory(cfg...))
		model.Routes[endpoint] = append(model.Routes[endpoint], http.MethodGet)
		middleware.RouteManager.Add(endpoint)
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

func convertGinPathToCasbinKeyMatch3(ginPath string) string {
	// Match :param style and replace with {param}
	re := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
	return re.ReplaceAllString(ginPath, `{$1}`)
}

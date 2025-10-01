package types

import (
	"context"
	"net/http"
	"net/url"

	"github.com/forbearing/gst/types/consts"
	"github.com/gin-gonic/gin"
)

type ControllerContext struct {
	Username string // currrent login user.
	UserID   string // currrent login user id
	Route    string
	Params   map[string]string
	Query    map[string][]string

	RequestID string
	TraceID   string
	PSpanID   string
	SpanID    string
	Seq       int
}

// NewControllerContext creates a ControllerContext from gin.Context
func NewControllerContext(c *gin.Context) *ControllerContext {
	if c == nil {
		return new(ControllerContext)
	}

	params := make(map[string]string)
	for _, key := range c.GetStringSlice(consts.PARAMS) {
		params[key] = c.Param(key)
	}

	return &ControllerContext{
		Route:     c.GetString(consts.CTX_ROUTE),
		Username:  c.GetString(consts.CTX_USERNAME),
		UserID:    c.GetString(consts.CTX_USER_ID),
		RequestID: c.GetString(consts.REQUEST_ID),
		TraceID:   c.GetString(consts.TRACE_ID),
		Params:    params,
		Query:     c.Request.URL.Query(),
	}
}

type DatabaseContext struct {
	Username string // currrent login user.
	UserID   string // currrent login user id
	Route    string
	Params   map[string]string
	Query    map[string][]string

	context   context.Context
	RequestID string
	TraceID   string
	PSpanID   string
	SpanID    string
	Seq       int
}

// NewDatabaseContext creates a DatabaseContext from gin.Context
//
// You can pass the custom context.Context to propagate span tracing,
// otherwise use the c.Request.Context().
func NewDatabaseContext(c *gin.Context, ctxs ...context.Context) *DatabaseContext {
	if c == nil {
		return new(DatabaseContext)
	}

	ctx := c.Request.Context()
	if len(ctxs) > 0 && ctxs[0] != nil {
		ctx = ctxs[0]
	}

	params := make(map[string]string)
	for _, key := range c.GetStringSlice(consts.PARAMS) {
		params[key] = c.Param(key)
	}

	return &DatabaseContext{
		context:   ctx,
		Route:     c.GetString(consts.CTX_ROUTE),
		Username:  c.GetString(consts.CTX_USERNAME),
		UserID:    c.GetString(consts.CTX_USER_ID),
		RequestID: c.GetString(consts.REQUEST_ID),
		TraceID:   c.GetString(consts.TRACE_ID),
		Params:    params,
		Query:     c.Request.URL.Query(),
	}
}

// Context converts *DatabaseContext to context.Context.
// It starts from the underlying ctx.context and conditionally injects extra metadata.
func (dc *DatabaseContext) Context() context.Context {
	if dc == nil || dc.context == nil {
		return context.Background()
	}

	c := dc.context
	if len(dc.Username) != 0 {
		c = context.WithValue(c, consts.CTX_USERNAME, dc.Username)
	}
	if len(dc.UserID) != 0 {
		c = context.WithValue(c, consts.CTX_USER_ID, dc.UserID)
	}
	if len(dc.RequestID) != 0 {
		c = context.WithValue(c, consts.REQUEST_ID, dc.RequestID)
	}
	if len(dc.TraceID) != 0 {
		c = context.WithValue(c, consts.TRACE_ID, dc.TraceID)
	}

	return c
}

type ServiceContext struct {
	Method       string        // http method
	Request      *http.Request // http request
	URL          *url.URL      // request url
	Header       http.Header   // http request header
	WriterHeader http.Header   // http writer header
	ClientIP     string        // client ip
	UserAgent    string        // user agent

	context context.Context
	Writer  http.ResponseWriter
	// Body    []byte

	// route parameters,
	//
	// eg: PUT /api/gists/:id/star
	// Params: map[string]string{"id": "xxxxx-mygistid-xxxxx"}
	//
	// eg: DELETE /api/user/:userid/shelf/shelfid/book
	// Params: map[string]string{"userid": "xxxxx-myuserid-xxxxx", "shelfid": "xxxxx-myshelfid-xxxxx"}
	Params map[string]string
	Query  map[string][]string

	SessionID string // session id
	Username  string // currrent login user.
	UserID    string // currrent login user id
	Route     string

	RequestID string
	TraceID   string
	PSpanID   string
	SpanID    string
	Seq       int

	ginCtx *gin.Context
	phase  consts.Phase
}

// NewServiceContext creates ServiceContext from gin.Context.
// Including request details, headers and user information.
//
// You can pass the custom context.Context to propagate span tracing,
// otherwise use the c.Request.Context().
func NewServiceContext(c *gin.Context, ctxs ...context.Context) *ServiceContext {
	if c == nil {
		return new(ServiceContext)
	}

	params := make(map[string]string)
	for _, key := range c.GetStringSlice(consts.PARAMS) {
		params[key] = c.Param(key)
	}

	ctx := c.Request.Context()
	if len(ctxs) > 0 && ctxs[0] != nil {
		ctx = ctxs[0]
	}

	return &ServiceContext{
		Request: c.Request,

		Method:       c.Request.Method,
		URL:          c.Request.URL,
		Header:       c.Request.Header,
		WriterHeader: c.Writer.Header(),
		ClientIP:     c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		Params:       params,
		Query:        c.Request.URL.Query(),

		Route:     c.GetString(consts.CTX_ROUTE),
		Username:  c.GetString(consts.CTX_USERNAME),
		UserID:    c.GetString(consts.CTX_USER_ID),
		SessionID: c.GetString(consts.CTX_SESSION_ID),

		RequestID: c.GetString(consts.REQUEST_ID),
		TraceID:   c.GetString(consts.TRACE_ID),

		ginCtx:  c,
		context: ctx,
		Writer:  c.Writer,
	}
}

// Context converts *ServiceContex to context.Context.
// It starts from the underlying ctx.context and conditionally injects extra metadata.
func (sc *ServiceContext) Context() context.Context {
	if sc == nil || sc.context == nil {
		return context.Background()
	}

	c := sc.context
	if len(sc.Username) != 0 {
		c = context.WithValue(c, consts.CTX_USERNAME, sc.Username)
	}
	if len(sc.UserID) != 0 {
		c = context.WithValue(c, consts.CTX_USER_ID, sc.UserID)
	}
	if len(sc.RequestID) != 0 {
		c = context.WithValue(c, consts.REQUEST_ID, sc.RequestID)
	}
	if len(sc.TraceID) != 0 {
		c = context.WithValue(c, consts.TRACE_ID, sc.TraceID)
	}

	return c
}

func (sc *ServiceContext) DatabaseContext() *DatabaseContext {
	return NewDatabaseContext(sc.ginCtx, sc.context)
}

func (sc *ServiceContext) Data(code int, contentType string, data []byte) {
	sc.ginCtx.Data(code, contentType, data)
}

func (sc *ServiceContext) HTML(code int, name string, obj any) {
	sc.ginCtx.HTML(code, name, obj)
}

func (sc *ServiceContext) Redirect(code int, location string) {
	sc.ginCtx.Redirect(code, location)
}

func (sc *ServiceContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	sc.ginCtx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (sc *ServiceContext) Cookie(name string) (string, error) {
	return sc.ginCtx.Cookie(name)
}

func (sc *ServiceContext) SetPhase(phase consts.Phase) { sc.phase = phase }
func (sc *ServiceContext) GetPhase() consts.Phase      { return sc.phase }
func (sc *ServiceContext) WithPhase(phase consts.Phase) *ServiceContext {
	sc.phase = phase
	return sc
}

type ModelContext struct {
	dbctx *DatabaseContext
}

func NewModelContext(dbctx *DatabaseContext, ctx context.Context) *ModelContext {
	if ctx == nil {
		ctx = context.Background()
	}

	return &ModelContext{
		dbctx: &DatabaseContext{
			context:  ctx,
			Username: dbctx.Username,
			UserID:   dbctx.UserID,
			Route:    dbctx.Route,
			Params:   dbctx.Params,
			Query:    dbctx.Query,

			RequestID: dbctx.RequestID,
			TraceID:   dbctx.TraceID,
			PSpanID:   dbctx.PSpanID,
			SpanID:    dbctx.SpanID,
			Seq:       dbctx.Seq,
		},
	}
}

func (mc *ModelContext) Context() context.Context {
	return mc.dbctx.context
}

func (mc *ModelContext) DatabaseContext() *DatabaseContext {
	return mc.dbctx
}

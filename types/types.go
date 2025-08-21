package types

import (
	"context"
	"net/http"
	"net/url"

	"github.com/forbearing/golib/types/consts"
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

	RequestID string
	TraceID   string
	PSpanID   string
	SpanID    string
	Seq       int
}

// NewDatabaseContext creates a DatabaseContext from gin.Context
func NewDatabaseContext(c *gin.Context) *DatabaseContext {
	if c == nil {
		return new(DatabaseContext)
	}

	params := make(map[string]string)
	for _, key := range c.GetStringSlice(consts.PARAMS) {
		params[key] = c.Param(key)
	}

	return &DatabaseContext{
		Route:     c.GetString(consts.CTX_ROUTE),
		Username:  c.GetString(consts.CTX_USERNAME),
		UserID:    c.GetString(consts.CTX_USER_ID),
		RequestID: c.GetString(consts.REQUEST_ID),
		TraceID:   c.GetString(consts.TRACE_ID),
		Params:    params,
		Query:     c.Request.URL.Query(),
	}
}

// NewGormContext converts *types.DatabaseContext to context.Context for use with gorm custom logger.
func NewGormContext(ctx *DatabaseContext) context.Context {
	c := context.Background()
	if ctx == nil {
		return c
	}
	c = context.WithValue(c, consts.CTX_USERNAME, ctx.Username)
	c = context.WithValue(c, consts.CTX_USER_ID, ctx.UserID)
	c = context.WithValue(c, consts.REQUEST_ID, ctx.RequestID)
	c = context.WithValue(c, consts.TRACE_ID, ctx.TraceID)
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

	Context context.Context
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
func NewServiceContext(c *gin.Context) *ServiceContext {
	if c == nil {
		return new(ServiceContext)
	}

	params := make(map[string]string)
	for _, key := range c.GetStringSlice(consts.PARAMS) {
		params[key] = c.Param(key)
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
		Context: c.Request.Context(),
		Writer:  c.Writer,
	}
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

func (sc *ServiceContext) SetPhase(phase consts.Phase) { sc.phase = phase }
func (sc *ServiceContext) GetPhase() consts.Phase      { return sc.phase }
func (sc *ServiceContext) WithPhase(phase consts.Phase) *ServiceContext {
	sc.phase = phase
	return sc
}

type ControllerConfig[M Model] struct {
	DB        any // only support *gorm.DB
	TableName string
}

package helper

import (
	"context"

	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
)

// NewControllerContext creates a ControllerContext from gin.Context
func NewControllerContext(c *gin.Context) *types.ControllerContext {
	if c == nil {
		return new(types.ControllerContext)
	}
	return &types.ControllerContext{
		Route:     c.GetString(consts.CTX_ROUTE),
		Username:  c.GetString(consts.CTX_USERNAME),
		UserId:    c.GetString(consts.CTX_USER_ID),
		RequestId: c.GetString(consts.REQUEST_ID),
		TraceId:   c.GetString(consts.TRACE_ID),
	}
}

// NewDatabaseContext creates a DatabaseContext from gin.Context
func NewDatabaseContext(c *gin.Context) *types.DatabaseContext {
	if c == nil {
		return new(types.DatabaseContext)
	}
	return &types.DatabaseContext{
		Route:     c.GetString(consts.CTX_ROUTE),
		Username:  c.GetString(consts.CTX_USERNAME),
		UserId:    c.GetString(consts.CTX_USER_ID),
		RequestId: c.GetString(consts.REQUEST_ID),
		TraceId:   c.GetString(consts.TRACE_ID),
	}
}

// NewServiceContext creates ServiceContext from gin.Context.
// Including request details, headers and user information.
func NewServiceContext(c *gin.Context) *types.ServiceContext {
	return &types.ServiceContext{
		Request: c.Request,

		Method:       c.Request.Method,
		URL:          c.Request.URL,
		Header:       c.Request.Header,
		WriterHeader: c.Writer.Header(),
		ClientIP:     c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),

		Route:     c.GetString(consts.CTX_ROUTE),
		Username:  c.GetString(consts.CTX_USERNAME),
		UserId:    c.GetString(consts.CTX_USER_ID),
		SessionId: c.GetString(consts.CTX_SESSION_ID),

		RequestId: c.GetString(consts.REQUEST_ID),
		TraceId:   c.GetString(consts.TRACE_ID),
	}
}

// NewGormContext converts *types.DatabaseContext to context.Context for use with gorm custom logger.
func NewGormContext(ctx *types.DatabaseContext) context.Context {
	c := context.Background()
	if ctx == nil {
		return c
	}
	c = context.WithValue(c, consts.CTX_USERNAME, ctx.Username)
	c = context.WithValue(c, consts.CTX_USER_ID, ctx.UserId)
	c = context.WithValue(c, consts.REQUEST_ID, ctx.RequestId)
	c = context.WithValue(c, consts.TRACE_ID, ctx.TraceId)
	return c
}

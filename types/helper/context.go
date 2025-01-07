package helper

import (
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
		Username:  c.GetString(consts.CTX_USERNAME),
		UserId:    c.GetString(consts.CTX_USER_ID),
		RequestId: c.GetString(consts.REQUEST_ID),
	}
}

// NewServiceContext build ServiceContext from gin.Context.
// Including request details, headers and user information.
func NewServiceContext(c *gin.Context) *types.ServiceContext {
	var requestId string
	val, _ := c.Get(consts.REQUEST_ID)
	switch v := val.(type) {
	case string:
		requestId = v
	}

	return &types.ServiceContext{
		Request: c.Request,

		Method:       c.Request.Method,
		URL:          c.Request.URL,
		Header:       c.Request.Header,
		WriterHeader: c.Writer.Header(),
		ClientIP:     c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),

		Username:  c.GetString(consts.CTX_USERNAME),
		UserId:    c.GetString(consts.CTX_USER_ID),
		SessionId: c.GetString(consts.CTX_SESSION_ID),
		RequestId: requestId,
	}
}

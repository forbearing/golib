package middleware

import (
	"github.com/forbearing/golib/authz/rbac"
	"github.com/forbearing/golib/logger"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Authz 中间件开启时必须也开启 JwtAuth 中间件, 并且位于 JwtAuth 之后.
// 要不然 username 会为空则则 isAllow 总是为 false.
func Authz() gin.HandlerFunc {
	return func(c *gin.Context) {
		var allow bool
		var err error
		sub := c.GetString(consts.CTX_USERNAME)
		obj := c.Request.URL.Path
		act := c.Request.Method

		if sub != consts.ROOT && sub != consts.ADMIN {
			sub = c.GetString(consts.CTX_USER_ID)
		}
		if allow, err = rbac.Enforcer.Enforce(sub, obj, act); err != nil {
			zap.S().Error(err)
			ResponseJSON(c, CodeFailure)
			c.Abort()
			return
		}
		logger.Authz.Infoz("",
			zap.String("sub", sub),
			zap.String("obj", obj),
			zap.String("act", act),
			zap.Bool("res", allow),
		)
		if allow {
			c.Next()
		} else {
			ResponseJSON(c, CodeForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

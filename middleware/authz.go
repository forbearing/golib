package middleware

import (
	"github.com/forbearing/golib/authz/rbac/basic"
	"github.com/forbearing/golib/logger"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Authz 中间件开启时必须也开启 JwtAuth 中间件, 要不然 username 会为空则则 isAllow 总是为 false.
func Authz() gin.HandlerFunc {
	return func(c *gin.Context) {
		// isAllow, err := rbac.RBAC.Enforcer(c.GetString(consts.CTX_USERNAME), c.Request.URL.Path, c.Request.Method)
		var isAllow bool
		var err error
		if c.GetString(consts.CTX_USERNAME) == consts.ROOT || c.GetString(consts.CTX_USERNAME) == consts.ADMIN {
			isAllow, err = basic.RBAC.Enforcer(c.GetString(consts.CTX_USERNAME), c.Request.URL.Path, c.Request.Method)
		} else {
			isAllow, err = basic.RBAC.Enforcer(c.GetString(consts.CTX_USER_ID), c.Request.URL.Path, c.Request.Method)
		}
		logger.Authz.Infoz("", zap.String("username", c.GetString(consts.CTX_USERNAME)),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Bool("is_allow", isAllow),
		)
		if err != nil {
			zap.S().Error(err)
			ResponseJSON(c, CodeFailure)
			c.Abort()
			return
		}
		if isAllow {
			c.Next()
		} else {
			ResponseJSON(c, CodeForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

package middleware

import (
	"github.com/forbearing/golib/rbac"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAllow, err := rbac.RBAC.Enforcer(c.GetString(types.CTX_USERNAME), c.Request.URL.Path, c.Request.Method)
		if err != nil {
			zap.S().Error(err)
			ResponseJSON(c, CodeFailure)
			c.Abort()
			return
		}
		if isAllow {
			c.Next()
		} else {
			ResponseJSON(c, CodeNoPermission)
			c.Abort()
			return
		}
	}
}

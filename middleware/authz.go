package middleware

import (
	"github.com/forbearing/golib/rbac"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Authz() gin.HandlerFunc {
	return func(c *gin.Context) {
		// isAllow, err := rbac.RBAC.Enforcer(c.GetString(types.CTX_USERNAME), c.Request.URL.Path, c.Request.Method)
		var isAllow bool
		var err error
		if c.GetString(types.CTX_USERNAME) == "root" || c.GetString(types.CTX_USERNAME) == "admin" {
			isAllow, err = rbac.RBAC.Enforcer(c.GetString(types.CTX_USERNAME), c.Request.URL.Path, c.Request.Method)
		} else {
			isAllow, err = rbac.RBAC.Enforcer(c.GetString(types.CTX_USER_ID), c.Request.URL.Path, c.Request.Method)
		}
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

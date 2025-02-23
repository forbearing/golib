package middleware

import (
	"github.com/forbearing/golib/jwt"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.ParseTokenFromHeader(c.Request.Header)
		if err != nil {
			ResponseJSON(c, CodeInvalidToken)
			c.Abort()
			return
		}

		// 将当前请求的 username 信息保存到请求的上线 *gin.Context 中
		// 后续的处理函数可以通过 c.Get("username") 来获取当前请求的用户信息
		// TODO: 将 user id 和 username 定义成变量/常量
		c.Set(consts.CTX_USER_ID, claims.UserId)
		c.Set(consts.CTX_USERNAME, claims.Username)
		c.Set(consts.CTX_SESSION_ID, c.GetHeader("X-Session-Id"))
		c.Next()
	}
}

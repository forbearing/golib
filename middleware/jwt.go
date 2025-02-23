package middleware

import (
	"net/http"

	"github.com/forbearing/golib/jwt"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
)

// JwtAuth 效果如下:
// 1.重复登录之后，会刷新 accessToken, refreshToken, 之后老的 accessToken 是失效
// 2.换浏览器、换操作系统都需要重新登录，重新登录之后会挤掉其他设备、浏览器的登录
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, claims, err := jwt.ParseTokenFromHeader(c.Request.Header)
		if err != nil {
			ResponseJSON(c, NewCode(http.StatusUnauthorized, err.Error()))
			c.Abort()
			return
		}
		if err := jwt.Verify(claims, accessToken, c.Request.UserAgent()); err != nil {
			ResponseJSON(c, NewCode(http.StatusUnauthorized, err.Error()))
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

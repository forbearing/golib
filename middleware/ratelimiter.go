package middleware

import (
	"time"

	. "github.com/forbearing/golib/response"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map/v2"
	"golang.org/x/time/rate"
)

var ratelimiterMap = cmap.New[*rate.Limiter]()

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter, found := ratelimiterMap.Get(c.ClientIP())
		if !found {
			// 桶大小为100, 每秒最多5个请求.
			limiter = rate.NewLimiter(rate.Every(200*time.Millisecond), 100)
			ratelimiterMap.Set(c.ClientIP(), limiter)
		}
		// if err := limiter.Wait(context.TODO()); err != nil {
		// 	zap.S().Error(err)
		// 	ResponseJSON(c, CodeFailure)
		// 	c.Abort()
		// 	return
		// }
		if !limiter.Allow() {
			ResponseJSON(c, CodeTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}

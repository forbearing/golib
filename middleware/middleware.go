package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/jwt"
	"github.com/forbearing/golib/logger"
	pkgzap "github.com/forbearing/golib/logger/zap"
	"github.com/forbearing/golib/model"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

var (
	ratelimitMap   = cmap.New[*limiter.Limiter]()
	ratelimiterMap = cmap.New[*rate.Limiter]()
)

func Logger(filename ...string) gin.HandlerFunc {
	// return ginzap.Ginzap(pkgzap.NewGinLogger(filename...), time.RFC3339, true)
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		fields := []zapcore.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("username", c.GetString(types.CTX_USERNAME)),
			zap.String("user_id", c.GetString(types.CTX_USER_ID)),
			zap.String("log_id", c.GetString(types.REQUEST_ID)),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("latency", time.Since(start).String()),
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Gin.Error(e, fields...)
			}
		} else {
			logger.Gin.Info(path, fields...)
		}
	}
}

func Recovery(filename ...string) gin.HandlerFunc {
	// TODO: replace it using custom logger.
	return ginzap.RecoveryWithZap(pkgzap.NewGin(filename...), true)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 加一个 X-Session-ID 用来记录来自哪个 session
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Session-Id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if len(header) == 0 {
			// zap.S().Error("not found authorization header")
			ResponseJSON(c, CodeNeedLogin)
			c.Abort()
			return
		}

		// 按空格分割
		items := strings.SplitN(header, " ", 2)
		if !(len(items) == 2 && items[0] == "Bearer") {
			// zap.S().Error("authorization header is invalid")
			ResponseJSON(c, CodeInvalidToken)
			c.Abort()
			return
		}

		// items[1] 是获取到的 tokenString, 我们使用之前定义好的解析 jwt 的函数来解析它
		claims, err := jwt.ParseToken(items[1])
		if err != nil {
			// zap.S().Error(err)
			ResponseJSON(c, CodeInvalidToken)
			c.Abort()
			return
		}

		// 将当前请求的 username 信息保存到请求的上线 *gin.Context 中
		// 后续的处理函数可以通过 c.Get("username") 来获取当前请求的用户信息
		// TODO: 将 user id 和 username 定义成变量/常量
		c.Set(types.CTX_USER_ID, claims.UserId)
		c.Set(types.CTX_USERNAME, claims.Username)
		c.Set(types.CTX_SESSION_ID, c.GetHeader("X-Session-Id"))
		c.Next()
	}
}

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

func BaseAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		config.BaseAuthUsername: config.BaseAuthPassword,
	})
}

// OperationLogger 中间件必须放在最后一个.
func OperationLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "" {
			c.Next() // next()方法的作用是跳过该调用链去直接后面的中间件以及api路由
		}
		info := func() (string, string) {
			username := c.GetString(types.CTX_USERNAME)
			var table string
			items := strings.Split(c.Request.URL.Path, `/`)
			if len(items) > 0 {
				table = items[len(items)-1]
			}
			return username, table
		}
		switch c.Request.Method {
		case http.MethodGet:
		case http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch:
			username, table := info()
			if err := database.Database[*model.OperationLog]().Create(&model.OperationLog{
				IP:        c.ClientIP(),
				User:      username,
				Table:     table,
				Model:     table,
				Method:    c.Request.Method,
				URI:       c.Request.RequestURI,
				UserAgent: c.Request.UserAgent(),
			}); err != nil {
				logger.Global.Error(err)
				return
			}
		}
	}
}

func Gzip() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/download"}), gzip.WithExcludedExtensions([]string{".pdf", ".mp4"}))
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		lmt, exists := ratelimitMap.Get(c.ClientIP())
		if !exists {
			// This setting means:
			// create a 1 request/second limiter and
			// every token bucket in it will expire 1 hour after it was initially set.
			lmt = tollbooth.NewLimiter(2, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
			lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).
				SetMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete})
		}

		if err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request); err != nil {
			ResponseJSON(c, CodeTooManyRequests)
			logger.Global.Debug(err)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

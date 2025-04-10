package middleware

import (
	"strconv"
	"time"

	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/metrics"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger(filename ...string) gin.HandlerFunc {
	// return ginzap.Ginzap(pkgzap.NewGinLogger(filename...), time.RFC3339, true)
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Set(consts.CTX_ROUTE, path)
		c.Next()

		metrics.HttpRequestsTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(c.Writer.Status())).Inc()
		metrics.HttpRequestDuration.WithLabelValues(c.Request.Method, path, strconv.Itoa(c.Writer.Status())).Observe(time.Since(start).Seconds())

		fields := []zapcore.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String(consts.CTX_USERNAME, c.GetString(consts.CTX_USERNAME)),
			zap.String(consts.CTX_USER_ID, c.GetString(consts.CTX_USER_ID)),
			zap.String(consts.REQUEST_ID, c.GetString(consts.REQUEST_ID)),
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

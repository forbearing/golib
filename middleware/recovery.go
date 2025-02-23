package middleware

import (
	pkgzap "github.com/forbearing/golib/logger/zap"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func Recovery(filename ...string) gin.HandlerFunc {
	// TODO: replace it using custom logger.
	return ginzap.RecoveryWithZap(pkgzap.NewGin(filename...), true)
}

package middleware

import (
	"github.com/forbearing/gst/types/consts"
	"github.com/gin-gonic/gin"
)

// AuthMarker is a middleware that marks the current route as requiring authentication.
// This middleware sets a flag in gin.Context to indicate that the current route requires authentication.
func AuthMarker() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(consts.CTX_REQUIRES_AUTH, true)
		c.Next()
	}
}

package middleware

import (
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := util.UUID()
		c.Set(types.REQUEST_ID, id)
		c.Header("X-Request-Id", id)
	}
}

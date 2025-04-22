package middleware

import (
	"net/http"
	"strings"

	"github.com/forbearing/golib/database"
	model_log "github.com/forbearing/golib/model/log"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationLogger 中间件必须放在最后一个.
func OperationLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "" {
			c.Next() // next()方法的作用是跳过该调用链去直接后面的中间件以及api路由
		}
		info := func() (string, string) {
			username := c.GetString(consts.CTX_USERNAME)
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
			if err := database.Database[*model_log.OperationLog]().Create(&model_log.OperationLog{
				IP:        c.ClientIP(),
				User:      username,
				Table:     table,
				Model:     table,
				Method:    c.Request.Method,
				URI:       c.Request.RequestURI,
				UserAgent: c.Request.UserAgent(),
			}); err != nil {
				zap.S().Error(err)
				return
			}
		}
	}
}

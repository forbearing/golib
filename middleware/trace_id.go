package middleware

import (
	"github.com/forbearing/golib/types/consts"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get(consts.TRACEID)
		pspanId := c.Request.Header.Get(consts.SPANID)
		spanId := util.SpanID()
		requestId := util.RequestID()
		if len(traceId) == 0 {
			// If traceid is empty, it means that it is the first request.
			traceId = spanId
		}
		c.Set(consts.REQUEST_ID, requestId)
		c.Set(consts.TRACEID, traceId)
		c.Set(consts.SPANID, spanId)
		c.Set(consts.PSPANID, pspanId)
		c.Header(consts.HEADER_REQUEST_ID, requestId)
		c.Header(consts.HEADER_TRACE_ID, traceId)
		c.Header(consts.HEADER_SPAN_ID, spanId)
		c.Header(consts.HEADER_PSPAN_ID, pspanId)
		c.Next()
	}
}

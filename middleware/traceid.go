package middleware

import (
	"github.com/forbearing/gst/types/consts"
	"github.com/forbearing/gst/util"
	"github.com/gin-gonic/gin"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get(consts.TRACE_ID)
		pspanId := c.Request.Header.Get(consts.SPAN_ID)
		spanId := util.SpanID()
		if len(traceId) == 0 {
			// If traceid is empty, it means that it is the first request.
			traceId = spanId
		}
		requestId := traceId
		c.Set(consts.REQUEST_ID, requestId)
		c.Set(consts.TRACE_ID, traceId)
		c.Set(consts.PSPAN_ID, pspanId)
		c.Set(consts.SPAN_ID, spanId)
		c.Set(consts.SEQ, 0)
		c.Header(consts.HEADER_REQUEST_ID, requestId)
		c.Header(consts.HEADER_TRACE_ID, traceId)
		c.Header(consts.HEADER_SPAN_ID, spanId)
		c.Header(consts.HEADER_PSPAN_ID, pspanId)
		c.Next()
	}
}

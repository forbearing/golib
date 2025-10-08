package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/forbearing/gst/config"
	"github.com/forbearing/gst/provider/otel"
	"github.com/forbearing/gst/types/consts"
	"github.com/forbearing/gst/util"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Tracing returns a middleware that handles both trace ID generation and OpenTelemetry tracing
// This middleware combines the functionality of TraceID() and Tracing() middlewares
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		var traceId, spanId string
		var span trace.Span
		var ctx context.Context

		// If OTEL is enabled, create OpenTelemetry span and use its trace ID
		if otel.IsEnabled() {
			// Create span name from HTTP method and route
			spanName := c.Request.Method + " " + c.FullPath()
			if c.FullPath() == "" {
				spanName = c.Request.Method + " " + c.Request.URL.Path
			}

			// Start new span
			ctx, span = otel.StartSpan(c.Request.Context(), spanName)

			// Extract OTEL trace ID and span ID
			spanContext := span.SpanContext()
			if spanContext.HasTraceID() {
				traceId = spanContext.TraceID().String()
				spanId = spanContext.SpanID().String()
			}

			// Set span attributes
			span.SetAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.scheme", c.Request.URL.Scheme),
				attribute.String("http.host", c.Request.Host),
				attribute.String("http.target", c.Request.URL.Path),
				attribute.String("http.route", c.FullPath()),
				attribute.String("http.user_agent", c.Request.UserAgent()),
				attribute.String("http.remote_addr", c.ClientIP()),
			)

			// Add request headers as attributes (selective)
			if contentType := c.Request.Header.Get("Content-Type"); contentType != "" {
				span.SetAttributes(attribute.String("http.request.content_type", contentType))
			}
			if contentLength := c.Request.Header.Get("Content-Length"); contentLength != "" {
				span.SetAttributes(attribute.String("http.request.content_length", contentLength))
			}

			// Store span in context for use in handlers
			c.Set("otel_span", span)
			c.Request = c.Request.WithContext(ctx)

			// Defer span completion
			defer func() {
				// Record response attributes
				span.SetAttributes(
					attribute.Int("http.status_code", c.Writer.Status()),
					attribute.Int("http.response.size", c.Writer.Size()),
					attribute.String("http.response.content_type", c.Writer.Header().Get("Content-Type")),
				)

				// Set span status based on HTTP status code
				statusCode := c.Writer.Status()
				if statusCode >= 400 {
					span.SetStatus(codes.Error, "HTTP "+strconv.Itoa(statusCode))
					span.SetAttributes(attribute.Bool("error", true))
				} else {
					span.SetStatus(codes.Ok, "")
				}

				// Record any errors from the request context
				if len(c.Errors) > 0 {
					span.SetStatus(codes.Error, c.Errors.String())
					span.SetAttributes(attribute.Bool("error", true))
					for i, err := range c.Errors {
						span.SetAttributes(attribute.String("error."+strconv.Itoa(i), err.Error()))
					}
				}

				span.End()
			}()
		} else {
			// Fallback to custom ID generation if OTEL is not enabled
			customTraceId := c.Request.Header.Get(consts.TRACE_ID)
			customSpanId := util.SpanID()
			if len(customTraceId) == 0 {
				customTraceId = customSpanId
			}
			traceId = customTraceId
			spanId = customSpanId
		}

		// Set unified trace ID and request ID in gin context
		c.Set(consts.REQUEST_ID, traceId) // Use traceId as requestId
		c.Set(consts.TRACE_ID, traceId)
		c.Set(consts.SPAN_ID, spanId)
		c.Set(consts.SEQ, 0)

		// Set X-Trace-ID header for frontend
		c.Header("X-Trace-ID", traceId)

		// Add gst trace IDs as span attributes if OTEL is enabled
		if otel.IsEnabled() && span != nil {
			span.SetAttributes(
				attribute.String(fmt.Sprintf("%s.trace_id", config.App.OTEL.ServiceName), traceId),
				attribute.String(fmt.Sprintf("%s.span_id", config.App.OTEL.ServiceName), spanId),
				attribute.String(fmt.Sprintf("%s.request_id", config.App.OTEL.ServiceName), traceId),
			)

			// Record start time for duration calculation
			start := time.Now()
			c.Set("request_start_time", start)

			// Process request
			c.Next()

			// Record duration
			duration := time.Since(start)
			span.SetAttributes(attribute.Int64("http.duration_ms", duration.Milliseconds()))
		} else {
			// Process request without tracing
			c.Next()
		}
	}
}

// GetSpanFromContext retrieves the OpenTelemetry span from Gin context
func GetSpanFromContext(c *gin.Context) trace.Span {
	if span, exists := c.Get("otel_span"); exists {
		if otelSpan, ok := span.(trace.Span); ok {
			return otelSpan
		}
	}
	return trace.SpanFromContext(c.Request.Context())
}

// AddSpanTags adds custom tags to the current span
func AddSpanTags(c *gin.Context, tags map[string]any) {
	span := GetSpanFromContext(c)
	if span != nil && span.IsRecording() {
		otel.AddSpanTags(span, tags)
	}
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(c *gin.Context, name string, attrs ...attribute.KeyValue) {
	span := GetSpanFromContext(c)
	if span != nil && span.IsRecording() {
		otel.AddSpanEvent(span, name, attrs...)
	}
}

// RecordError records an error in the current span
func RecordError(c *gin.Context, err error) {
	span := GetSpanFromContext(c)
	if span != nil && span.IsRecording() {
		otel.RecordError(span, err)
	}
}

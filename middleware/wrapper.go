package middleware

import (
	"fmt"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/provider/jaeger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// middlewareWrapper wraps any gin middleware with jaeger tracing capabilities.
// It creates a span for the middleware execution and records performance metrics.
//
// Parameters:
//   - name: The name of the middleware for tracing identification
//   - middleware: The gin.HandlerFunc to be wrapped
//
// Returns:
//   - A new gin.HandlerFunc with tracing capabilities
//
// Example:
//
//	wrappedLogger := middlewareWrapper("logger", Logger())
//	wrappedAuth := middlewareWrapper("jwt-auth", JWT())
//	router.Use(wrappedLogger, wrappedAuth)
func middlewareWrapper(name string, middleware gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip tracing if Jaeger is not enabled
		if !jaeger.IsEnabled() {
			middleware(c)
			return
		}

		// Create span name with middleware prefix
		spanName := fmt.Sprintf("middleware.%s", name)

		// Start new span for middleware execution
		ctx, span := jaeger.StartSpan(c.Request.Context(), spanName)
		defer span.End()

		// Update request context with the new span context
		c.Request = c.Request.WithContext(ctx)

		// Set span attributes
		span.SetAttributes(
			attribute.String("middleware.name", name),
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.path", c.Request.URL.Path),
			attribute.String("http.route", c.FullPath()),
		)

		// Record start time
		start := time.Now()

		// Execute the wrapped middleware
		middleware(c)

		// Record execution duration
		duration := time.Since(start)
		span.SetAttributes(
			attribute.Int64("middleware.duration_ms", duration.Milliseconds()),
			attribute.Int64("middleware.duration_ns", duration.Nanoseconds()),
		)

		// Check if middleware caused any errors (based on response status)
		if c.Writer.Status() >= 400 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", c.Writer.Status()))
			span.SetAttributes(
				attribute.Int("http.status_code", c.Writer.Status()),
				attribute.Bool("middleware.error", true),
			)
		} else {
			span.SetStatus(codes.Ok, "")
			span.SetAttributes(
				attribute.Int("http.status_code", c.Writer.Status()),
				attribute.Bool("middleware.error", false),
			)
		}

		// Add service name as attribute
		if config.App.Jaeger.ServiceName != "" {
			span.SetAttributes(
				attribute.String("service.name", config.App.Jaeger.ServiceName),
			)
		}
	}
}

// wrapMiddlewares wraps multiple middlewares with tracing capabilities.
// This is a convenience function for wrapping multiple middlewares at once.
//
// Parameters:
//   - middlewares: A map where key is the middleware name and value is the gin.HandlerFunc
//
// Returns:
//   - A slice of wrapped gin.HandlerFunc with tracing capabilities
//
// Example:
//
//	wrapped := wrapMiddlewares(map[string]gin.HandlerFunc{
//	    "logger": Logger(),
//	    "cors": Cors(),
//	    "recovery": Recovery(),
//	})
//	router.Use(wrapped...)
func wrapMiddlewares(middlewares map[string]gin.HandlerFunc) []gin.HandlerFunc {
	wrapped := make([]gin.HandlerFunc, 0, len(middlewares))
	for name, middleware := range middlewares {
		wrapped = append(wrapped, middlewareWrapper(name, middleware))
	}
	return wrapped
}

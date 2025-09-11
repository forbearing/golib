package tracing

import (
	"context"
	"fmt"
	"time"

	"github.com/forbearing/golib/provider/jaeger"
	"github.com/forbearing/golib/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingWrapper wraps a Cache implementation with distributed tracing capabilities
type TracingWrapper[T any] struct {
	cache     types.Cache[T]
	ctx       context.Context
	cacheType string
}

// NewTracingWrapper creates a new tracing wrapper for the given cache
func NewTracingWrapper[T any](cache types.Cache[T], cacheType string) *TracingWrapper[T] {
	return &TracingWrapper[T]{
		cache:     cache,
		ctx:       context.Background(),
		cacheType: cacheType,
	}
}

// WithContext sets the context for tracing
func (tw *TracingWrapper[T]) WithContext(ctx context.Context) types.Cache[T] {
	return &TracingWrapper[T]{
		cache:     tw.cache,
		ctx:       ctx,
		cacheType: tw.cacheType,
	}
}

// Set stores a key-value pair with tracing
func (tw *TracingWrapper[T]) Set(key string, value T, ttl time.Duration) error {
	_, span := tw.startSpan("Cache.Set")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("cache.operation", "set"),
		attribute.String("cache.key", key),
		attribute.String("cache.ttl", ttl.String()),
		attribute.String("cache.type", tw.cacheType),
	)

	// Record start time for duration measurement
	start := time.Now()

	// Call the underlying cache implementation
	err := tw.cache.Set(key, value, ttl)

	// Record operation duration
	duration := time.Since(start)
	span.SetAttributes(
		attribute.String("cache.duration", duration.String()),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("Failed to set cache key: %v", err))
		return err
	}

	span.SetStatus(codes.Ok, "Cache key set successfully")
	return nil
}

// Get retrieves a value by key with tracing
func (tw *TracingWrapper[T]) Get(key string) (T, error) {
	_, span := tw.startSpan("Cache.Get")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("cache.operation", "get"),
		attribute.String("cache.key", key),
		attribute.String("cache.type", tw.cacheType),
	)

	// Record start time for duration measurement
	start := time.Now()

	// Call the underlying cache implementation
	value, err := tw.cache.Get(key)

	// Record operation duration
	duration := time.Since(start)
	span.SetAttributes(
		attribute.String("cache.duration", duration.String()),
	)

	if err != nil {
		span.SetAttributes(
			attribute.Bool("cache.hit", false),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("Failed to get cache key: %v", err))
		return value, err
	}

	span.SetAttributes(
		attribute.Bool("cache.hit", true),
	)
	span.SetStatus(codes.Ok, "Cache key retrieved successfully")
	return value, nil
}

// Peek retrieves a value by key without affecting its position with tracing
func (tw *TracingWrapper[T]) Peek(key string) (T, error) {
	_, span := tw.startSpan("Cache.Peek")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("cache.operation", "peek"),
		attribute.String("cache.key", key),
		attribute.String("cache.type", tw.cacheType),
	)

	// Record start time for duration measurement
	start := time.Now()

	// Call the underlying cache implementation
	value, err := tw.cache.Peek(key)

	// Record operation duration
	duration := time.Since(start)
	span.SetAttributes(
		attribute.String("cache.duration", duration.String()),
	)

	if err != nil {
		span.SetAttributes(
			attribute.Bool("cache.hit", false),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("Failed to peek cache key: %v", err))
		return value, err
	}

	span.SetAttributes(
		attribute.Bool("cache.hit", true),
	)
	span.SetStatus(codes.Ok, "Cache key peeked successfully")
	return value, nil
}

// Delete removes a key from the cache with tracing
func (tw *TracingWrapper[T]) Delete(key string) error {
	_, span := tw.startSpan("Cache.Delete")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("cache.operation", "delete"),
		attribute.String("cache.key", key),
		attribute.String("cache.type", tw.cacheType),
	)

	// Record start time for duration measurement
	start := time.Now()

	// Call the underlying cache implementation
	err := tw.cache.Delete(key)

	// Record operation duration
	duration := time.Since(start)
	span.SetAttributes(
		attribute.String("cache.duration", duration.String()),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("Failed to delete cache key: %v", err))
		return err
	}

	span.SetStatus(codes.Ok, "Cache key deleted successfully")
	return nil
}

// Exists checks if a key exists in the cache with tracing
func (tw *TracingWrapper[T]) Exists(key string) bool {
	_, span := tw.startSpan("Cache.Exists")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("cache.operation", "exists"),
		attribute.String("cache.key", key),
		attribute.String("cache.type", tw.cacheType),
	)

	// Record start time for duration measurement
	start := time.Now()

	// Call the underlying cache implementation
	exists := tw.cache.Exists(key)

	// Record operation duration
	duration := time.Since(start)
	span.SetAttributes(
		attribute.String("cache.duration", duration.String()),
		attribute.Bool("cache.exists", exists),
	)

	span.SetStatus(codes.Ok, "Cache key existence checked successfully")
	return exists
}

// Len returns the number of items in the cache with tracing
func (tw *TracingWrapper[T]) Len() int {
	_, span := tw.startSpan("Cache.Len")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("cache.operation", "len"),
		attribute.String("cache.type", tw.cacheType),
	)

	// Record start time for duration measurement
	start := time.Now()

	// Call the underlying cache implementation
	length := tw.cache.Len()

	// Record operation duration
	duration := time.Since(start)
	span.SetAttributes(
		attribute.String("cache.duration", duration.String()),
		attribute.Int("cache.length", length),
	)

	span.SetStatus(codes.Ok, "Cache length retrieved successfully")
	return length
}

// Clear removes all items from the cache with tracing
func (tw *TracingWrapper[T]) Clear() {
	_, span := tw.startSpan("Cache.Clear")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("cache.operation", "clear"),
		attribute.String("cache.type", tw.cacheType),
	)

	// Record start time for duration measurement
	start := time.Now()

	// Call the underlying cache implementation
	tw.cache.Clear()

	// Record operation duration
	duration := time.Since(start)
	span.SetAttributes(
		attribute.String("cache.duration", duration.String()),
	)

	span.SetStatus(codes.Ok, "Cache cleared successfully")
}

// startSpan creates a new span for the given operation
func (tw *TracingWrapper[T]) startSpan(operationName string) (context.Context, trace.Span) {
	tracer := jaeger.GetTracer()
	ctx, span := tracer.Start(tw.ctx, operationName)
	return ctx, span
}

package cache

import (
	"context"
	"testing"

	"github.com/forbearing/golib/cache/lrue"
	tracing "github.com/forbearing/golib/cache/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTracingWrapper(t *testing.T) {
	// Initialize cache
	baseCache := lrue.Cache[string]()
	require.NotNil(t, baseCache)

	// Create tracing wrapper
	ctx := context.Background()
	tracingCache := baseCache.WithContext(ctx)
	require.NotNil(t, tracingCache)

	// Test Set operation
	err := tracingCache.Set("test-key", "test-value", 0)
	assert.NoError(t, err)

	// Test Get operation
	value, err := tracingCache.Get("test-key")
	assert.NoError(t, err)
	assert.Equal(t, "test-value", value)

	// Test Peek operation
	peekValue, err := tracingCache.Peek("test-key")
	assert.NoError(t, err)
	assert.Equal(t, "test-value", peekValue)

	// Test Exists operation
	exists := tracingCache.Exists("test-key")
	assert.True(t, exists)

	// Test Len operation
	length := tracingCache.Len()
	assert.Equal(t, 1, length)

	// Test Delete operation
	err = tracingCache.Delete("test-key")
	assert.NoError(t, err)
	_, err = tracingCache.Get("test-key")
	assert.Error(t, err)

	// Test Clear operation
	err = tracingCache.Set("key1", "value1", 0)
	assert.NoError(t, err)
	err = tracingCache.Set("key2", "value2", 0)
	assert.NoError(t, err)
	length = tracingCache.Len()
	assert.Equal(t, 2, length)

	tracingCache.Clear()
	length = tracingCache.Len()
	assert.Equal(t, 0, length)
}

func TestNewTracingWrapper(t *testing.T) {
	baseCache := lrue.Cache[int]()
	wrapper := tracing.NewTracingWrapper(baseCache)

	require.NotNil(t, wrapper)
	// Test that the wrapper works by calling a method
	err := wrapper.Set("test", 123, 0)
	assert.NoError(t, err)
}

func TestWithTracing(t *testing.T) {
	baseCache := lrue.Cache[string]()
	wrapper := WithTracing(baseCache)

	require.NotNil(t, wrapper)
	// Test that the wrapper works by calling a method
	err := wrapper.Set("test", "value", 0)
	assert.NoError(t, err)
}

func TestCacheWithTracing(t *testing.T) {
	ctx := context.Background()
	cache := CacheWithTracing[string](ctx)

	require.NotNil(t, cache)

	// Test basic operations
	err := cache.Set("key", "value", 0)
	assert.NoError(t, err)
	value, err := cache.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", value)
}

func TestExpirableCacheWithTracing(t *testing.T) {
	ctx := context.Background()
	cache := ExpirableCacheWithTracing[string](ctx)

	require.NotNil(t, cache)

	// Test basic operations
	err := cache.Set("key", "value", 0)
	assert.NoError(t, err)
	value, err := cache.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", value)
}

func TestTracingWrapperWithContext(t *testing.T) {
	baseCache := lrue.Cache[string]()
	ctx1 := context.Background()
	tracingCache1 := baseCache.WithContext(ctx1)

	// Create new context and wrapper
	ctx2 := context.WithValue(context.Background(), "test", "value")
	tracingCache2 := tracingCache1.WithContext(ctx2)

	require.NotNil(t, tracingCache2)

	// Both should work independently
	err := tracingCache1.Set("key1", "value1", 0)
	assert.NoError(t, err)
	err = tracingCache2.Set("key2", "value2", 0)
	assert.NoError(t, err)

	value1, err := tracingCache1.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value1)

	value2, err := tracingCache2.Get("key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", value2)
}

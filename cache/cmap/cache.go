package cmap

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/forbearing/golib/cache/tracing"
	"github.com/forbearing/golib/types"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	cacheMap = cmap.New[any]()
	mu       sync.Mutex
)

func Init() error {
	return nil
}

type cache[T any] struct {
	c cmap.ConcurrentMap[string, T]
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	if exists {
		return val.(*cache[T])
	}

	mu.Lock()
	defer mu.Unlock()

	val, exists = cacheMap.Get(key)
	if !exists {
		val = &cache[T]{c: cmap.New[T]()}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) error {
	c.c.Set(key, value)
	return nil
}

func (c *cache[T]) Get(key string) (T, error) {
	value, exists := c.c.Get(key)
	if !exists {
		var zero T
		return zero, types.ErrEntryNotFound
	}
	return value, nil
}

func (c *cache[T]) Peek(key string) (T, error) {
	return c.Get(key)
}

func (c *cache[T]) Delete(key string) error {
	c.c.Remove(key)
	return nil
}

func (c *cache[T]) Exists(key string) bool {
	return c.c.Has(key)
}

func (c *cache[T]) Len() int {
	return c.c.Count()
}

func (c *cache[T]) Clear() {
	c.c.Clear()
}

// WithContext returns a new Cache instance with the given context for tracing
func (c *cache[T]) WithContext(ctx context.Context) types.Cache[T] {
	return tracing.NewTracingWrapper(c, "cmap").WithContext(ctx)
}

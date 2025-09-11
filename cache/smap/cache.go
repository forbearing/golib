package smap

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/forbearing/golib/cache/tracing"
	"github.com/forbearing/golib/types"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var cacheMap = cmap.New[any]()

func Init() error {
	return nil
}

type cache[T any] struct {
	m sync.Map
	n int64
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	if !exists {
		val = &cache[T]{m: sync.Map{}}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) error {
	_, loaded := c.m.LoadOrStore(key, value)
	if loaded {
		c.m.Store(key, value)
	} else {
		atomic.AddInt64(&c.n, 1)
	}
	return nil
}

func (c *cache[T]) Get(key string) (T, error) {
	v, ok1 := c.m.Load(key)
	if !ok1 {
		var zero T
		return zero, types.ErrEntryNotFound
	}
	_v, ok2 := v.(T)
	if !ok2 {
		var zero T
		return zero, types.ErrEntryNotFound
	}
	return _v, nil
}

func (c *cache[T]) Peek(key string) (T, error) {
	return c.Get(key)
}

func (c *cache[T]) Delete(key string) error {
	_, exists := c.m.LoadAndDelete(key)
	if exists {
		atomic.AddInt64(&c.n, -1)
	}
	return nil
}

func (c *cache[T]) Exists(key string) bool {
	_, exists := c.m.Load(key)
	return exists
}

func (c *cache[T]) Len() int {
	return int(c.n)
}

func (c *cache[T]) Clear() {
	c.m.Range(func(key, _ any) bool {
		c.m.Delete(key)
		return true
	})
	atomic.StoreInt64(&c.n, 0)
}

// WithContext returns a new Cache instance with the given context for tracing
func (c *cache[T]) WithContext(ctx context.Context) types.Cache[T] {
	return tracing.NewTracingWrapper(c, "smap").WithContext(ctx)
}

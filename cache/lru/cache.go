package lru

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/forbearing/gst/cache/tracing"
	"github.com/forbearing/gst/config"
	"github.com/forbearing/gst/types"
	lru "github.com/hashicorp/golang-lru/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	cacheMap = cmap.New[any]()
	tmp      *lru.Cache[string, any] // tmp is a temporary cache used to check the config is correct.
	mu       sync.Mutex
)

func Init() (err error) {
	if tmp, err = lru.New[string, any](config.App.Cache.Capacity); err != nil {
		return err
	}
	tmp.Purge()
	return nil
}

type cache[T any] struct {
	c   *lru.Cache[string, T]
	ctx context.Context
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeFor[T]()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	if exists {
		//nolint:errcheck
		return val.(types.Cache[T])
	}

	mu.Lock()
	defer mu.Unlock()

	val, exists = cacheMap.Get(key)
	if !exists {
		// lru.New() only error on negative size.
		_lru, _ := lru.New[string, T](config.App.Cache.Capacity)
		val = tracing.NewTracingWrapper(&cache[T]{c: _lru, ctx: context.Background()}, "lru")
		cacheMap.Set(key, val)
	}
	//nolint:errcheck
	return val.(types.Cache[T])
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) error {
	c.c.Add(key, value)
	return nil
}

func (c *cache[T]) Get(key string) (T, error) {
	value, ok := c.c.Get(key)
	if !ok {
		var zero T
		return zero, types.ErrEntryNotFound
	}
	return value, nil
}

func (c *cache[T]) Peek(key string) (T, error) {
	value, ok := c.c.Get(key)
	if !ok {
		var zero T
		return zero, types.ErrEntryNotFound
	}
	return value, nil
}

func (c *cache[T]) Delete(key string) error {
	c.c.Remove(key)
	return nil
}

func (c *cache[T]) Exists(key string) bool {
	return c.c.Contains(key)
}

func (c *cache[T]) Len() int {
	return c.c.Len()
}

func (c *cache[T]) Clear() {
	c.c.Purge()
}

func (c *cache[T]) WithContext(ctx context.Context) types.Cache[T] {
	c.ctx = ctx
	return c
}

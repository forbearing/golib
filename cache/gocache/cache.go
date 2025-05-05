package gocache

import (
	"reflect"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	cmap "github.com/orcaman/concurrent-map/v2"
	pkgcache "github.com/patrickmn/go-cache"
)

var cacheMap = cmap.New[any]()

func Init() error {
	return nil
}

type cache[T any] struct {
	c *pkgcache.Cache
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	if !exists {
		val = &cache[T]{c: pkgcache.New(config.App.Cache.Expiration, 2*config.App.Cache.Expiration)}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) {
	c.c.Set(key, value, ttl)
}

func (c *cache[T]) Get(key string) (T, bool) {
	var zero T
	val, ok := c.c.Get(key)
	if !ok {
		return zero, false
	}
	if val == nil {
		return zero, false
	}
	return val.(T), ok
}

func (c *cache[T]) Peek(key string) (T, bool) {
	return c.Get(key)
}

func (c *cache[T]) Exists(key string) bool {
	_, exists := c.c.Get(key)
	return exists
}

func (c *cache[T]) Delete(key string) {
	c.c.Delete(key)
}

func (c *cache[T]) Len() int {
	return c.c.ItemCount()
}

func (c *cache[T]) Clear() {
	c.c.Flush()
}

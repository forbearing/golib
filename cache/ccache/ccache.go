package ccache

import (
	"reflect"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	"github.com/karlseguin/ccache/v3"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var cacheMap = cmap.New[any]()

func Init() (err error) {
	return nil
}

type cache[T any] struct {
	c *ccache.Cache[T]
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	if !exists {
		val = &cache[T]{c: ccache.New(ccache.Configure[T]().MaxSize(int64(config.App.Cache.Capacity)))}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) {
	c.c.Set(key, value, ttl)
}

func (c *cache[T]) Get(key string) (T, bool) {
	var zero T
	val := c.c.Get(key)
	if val == nil {
		return zero, false
	}
	if val.Expired() {
		return zero, false
	}
	return val.Value(), true
}

func (c *cache[T]) Peek(key string) (T, bool) {
	return c.Get(key)
}

func (c *cache[T]) Exists(key string) bool {
	val := c.c.Get(key)
	if val == nil {
		return false
	}
	if val.Expired() {
		return false
	}
	return true
}

func (c *cache[T]) Delete(key string) {
	c.c.Delete(key)
}

func (c *cache[T]) Keys() []string {
	return []string{}
}

func (c *cache[T]) Len() int {
	return c.c.ItemCount()
}

func (c *cache[T]) Flush() {
	c.c.Clear()
}


package ristretto

import (
	"reflect"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	cacheMap = cmap.New[any]()
	tmp      *ristretto.Cache[string, any]
)

func Init() (err error) {
	if tmp, err = ristretto.NewCache(buildConf[any]()); err != nil {
		return err
	}
	tmp.Close()
	return nil
}

type cache[T any] struct {
	c *ristretto.Cache[string, T]
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	if !exists {
		c, _ := ristretto.NewCache(buildConf[T]())
		val = &cache[T]{c: c}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) {
	c.c.Set(key, value, 1)
}

func (c *cache[T]) Get(key string) (T, bool) {
	return c.c.Get(key)
}

func (c *cache[T]) Peek(key string) (T, bool) {
	return c.c.Get(key)
}

func (c *cache[T]) Exists(key string) bool {
	_, exists := c.c.Get(key)
	return exists
}

func (c *cache[T]) Delete(key string) {
	c.c.Del(key)
}

func (c *cache[T]) Len() int {
	return -1
}

func (c *cache[T]) Clear() {
	c.c.Clear()
}

func buildConf[T any]() *ristretto.Config[string, T] {
	return &ristretto.Config[string, T]{
		NumCounters: int64(config.App.Cache.Capacity),
		MaxCost:     1 << 30,
		BufferItems: 64,
	}
}

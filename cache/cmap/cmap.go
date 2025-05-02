package cmap

import (
	"reflect"
	"time"

	"github.com/forbearing/golib/types"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var cacheMap = cmap.New[any]()

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
	if !exists {
		val = &cache[T]{c: cmap.New[T]()}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, _ time.Duration) { c.c.Set(key, value) }
func (c *cache[T]) Get(key string) (T, bool)                 { return c.c.Get(key) }
func (c *cache[T]) Peek(key string) (T, bool)                { return c.c.Get(key) }
func (c *cache[T]) Delete(key string)                        { c.c.Remove(key) }
func (c *cache[T]) Exists(key string) bool                   { return c.c.Has(key) }
func (c *cache[T]) Keys() []string                           { return c.c.Keys() }
func (c *cache[T]) Len() int                                 { return c.c.Count() }
func (c *cache[T]) Flush()                                   { c.c.Clear() }

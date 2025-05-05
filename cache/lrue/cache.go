// Package lrue is a expirable lru cache.
package lrue

import (
	"reflect"
	"sync"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	"github.com/hashicorp/golang-lru/v2/expirable"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	cacheMap = cmap.New[any]()
	mu       sync.Mutex
)

func Init() error { return nil }

type cache[T any] struct {
	c *expirable.LRU[string, T]
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
		val = &cache[T]{c: expirable.NewLRU[string, T](config.App.Cache.Capacity, nil, config.App.Cache.Expiration)}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}
func (c *cache[T]) Set(key string, value T, _ time.Duration) { c.c.Add(key, value) }
func (c *cache[T]) Get(key string) (T, bool)                 { return c.c.Get(key) }
func (c *cache[T]) Peek(key string) (T, bool)                { return c.c.Get(key) }
func (c *cache[T]) Delete(key string)                        { c.c.Remove(key) }
func (c *cache[T]) Exists(key string) bool                   { return c.c.Contains(key) }
func (c *cache[T]) Len() int                                 { return c.c.Len() }
func (c *cache[T]) Clear()                                   { c.c.Purge() }

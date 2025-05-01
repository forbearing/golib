package lru

import (
	"reflect"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	lru "github.com/hashicorp/golang-lru/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	cacheMap = cmap.New[any]()
	tmp      *lru.Cache[string, any] // tmp is a temporary cache used to check the config is correct.
)

func Init() (err error) {
	if tmp, err = lru.New[string, any](config.App.Cache.Capacity); err != nil {
		return err
	}
	tmp.Purge()
	return nil
}

type cache[T any] struct {
	c *lru.Cache[string, T]
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	// lru.New() only error on negative size.
	if !exists {
		_lru, _ := lru.New[string, T](config.App.Cache.Capacity)
		val = &cache[T]{c: _lru}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T)   { c.c.Add(key, value) }
func (c *cache[T]) Get(key string) (T, bool)  { return c.c.Get(key) }
func (c *cache[T]) Peek(key string) (T, bool) { return c.c.Get(key) }
func (c *cache[T]) Remove(key string)         { c.c.Remove(key) }
func (c *cache[T]) Exists(key string) bool    { return c.c.Contains(key) }
func (c *cache[T]) Keys() []string            { return c.c.Keys() }
func (c *cache[T]) Count() int                { return c.c.Len() }
func (c *cache[T]) Flush()                    { c.c.Purge() }

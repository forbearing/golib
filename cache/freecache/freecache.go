package freecache

import (
	"reflect"
	"time"

	"github.com/coocood/freecache"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.uber.org/zap"
)

var cacheMap = cmap.New[any]()

func Init() error {
	return nil
}

type cache[T any] struct {
	c *freecache.Cache
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	if !exists {
		val = &cache[T]{c: freecache.NewCache(config.App.Cache.Capacity)}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) {
	val, err := util.Marshal(value)
	if err != nil {
		zap.S().Error(err)
	} else {
		if err := c.c.Set([]byte(key), val, int(ttl.Seconds())); err != nil {
			zap.S().Error(err)
		}
	}
}

func (c *cache[T]) Get(key string) (T, bool) {
	var zero T
	val, err := c.c.Get([]byte(key))
	if err != nil {
		// if not found, not log error
		return zero, false
	}
	var result T
	err = util.Unmarshal(val, &result)
	if err != nil {
		zap.S().Error(err)
		return zero, false
	}
	return result, true
}

func (c *cache[T]) Peek(key string) (T, bool) {
	return c.Get(key)
}

func (c *cache[T]) Delete(key string) { c.c.Del([]byte(key)) }
func (c *cache[T]) Exists(key string) bool {
	_, err := c.c.Get([]byte(key))
	return err == nil
}

func (c *cache[T]) Keys() []string {
	keys := make([]string, 0)
	it := c.c.NewIterator()
	for entry := it.Next(); entry != nil; entry = it.Next() {
		keys = append(keys, string(entry.Key))
	}
	return keys
}

func (c *cache[T]) Len() int { return int(c.c.EntryCount()) }
func (c *cache[T]) Flush()   { c.c.Clear() }

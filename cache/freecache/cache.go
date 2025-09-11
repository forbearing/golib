package freecache

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/coocood/freecache"
	"github.com/forbearing/golib/cache/tracing"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
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
	c *freecache.Cache
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
		val = &cache[T]{c: freecache.NewCache(config.App.Cache.Capacity)}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) error {
	val, err := util.Marshal(value)
	if err != nil {
		return err
	}
	return c.c.Set([]byte(key), val, int(ttl.Seconds()))
}

func (c *cache[T]) Get(key string) (T, error) {
	var zero T
	val, err := c.c.Get([]byte(key))
	if err != nil {
		return zero, types.ErrEntryNotFound
	}
	var result T
	err = util.Unmarshal(val, &result)
	if err != nil {
		return zero, err
	}
	return result, nil
}

func (c *cache[T]) Peek(key string) (T, error) {
	return c.Get(key)
}

func (c *cache[T]) Delete(key string) error {
	c.c.Del([]byte(key))
	return nil
}

func (c *cache[T]) Exists(key string) bool {
	_, err := c.c.Get([]byte(key))
	return err == nil
}

func (c *cache[T]) Len() int {
	return int(c.c.EntryCount())
}

func (c *cache[T]) Clear() {
	c.c.Clear()
}

// WithContext returns a new Cache instance with the given context for tracing
func (c *cache[T]) WithContext(ctx context.Context) types.Cache[T] {
	return tracing.NewTracingWrapper(c, "freecache").WithContext(ctx)
}

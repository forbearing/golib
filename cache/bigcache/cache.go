package bigcache

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/forbearing/golib/cache/tracing"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	cacheMap = cmap.New[any]()

	tmp *bigcache.BigCache // tmp is a temporary cache used to check the config is correct.
	mu  sync.Mutex

	maxEntrySize     = 1024 * 64 // 64KB
	hardMaxCacheSize = 0
	verbose          = false
)

func Init() (err error) {
	if tmp, err = bigcache.NewBigCache(buildConfig()); err != nil {
		return err
	}
	tmp.Close()

	return nil
}

type cache[T any] struct {
	c *bigcache.BigCache
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
		val = &cache[T]{c: newBigCache()}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, _ time.Duration) error {
	val, err := util.Marshal(value)
	if err != nil {
		return err
	}
	return c.c.Set(key, val)
}

func (c *cache[T]) Get(key string) (T, error) {
	var zero T
	val, err := c.c.Get(key)
	if err != nil {
		// Return ErrEntryNotFound for not found cases
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
	c.c.Delete(key)
	return nil
}

func (c *cache[T]) Exists(key string) bool {
	_, err := c.c.Get(key)
	return err == nil
}

func (c *cache[T]) Len() int {
	return c.c.Len()
}

func (c *cache[T]) Clear() {
	c.c.Reset()
}

// WithContext returns a new Cache instance with the given context for tracing
func (c *cache[T]) WithContext(ctx context.Context) types.Cache[T] {
	return tracing.NewTracingWrapper(c, "bigcache").WithContext(ctx)
}

func newBigCache() *bigcache.BigCache {
	cache, _ := bigcache.NewBigCache(buildConfig())
	return cache
}

func buildConfig() bigcache.Config {
	return bigcache.Config{
		Shards:           config.App.Shards,
		LifeWindow:       config.App.Cache.LifeWindow,
		CleanWindow:      config.App.Cache.CleanWindow,
		MaxEntrySize:     maxEntrySize,
		HardMaxCacheSize: hardMaxCacheSize,
		Verbose:          verbose,
	}
}

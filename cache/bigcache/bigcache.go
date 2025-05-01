package bigcache

import (
	"reflect"
	"time"

	"github.com/allegro/bigcache"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	jsoniter "github.com/json-iterator/go"
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.uber.org/zap"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	cacheMap = cmap.New[any]()

	tmp *bigcache.BigCache // tmp is a temporary cache used to check the config is correct.

	defaultExpire    = time.Hour
	maxEntrySize     = 1024 // 1KB
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
	if !exists {
		val = &cache[T]{c: newBigCache()}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T) {
	val, err := json.Marshal(value)
	if err != nil {
		zap.S().Error(err)
	} else {
		if err := c.c.Set(key, val); err != nil {
			zap.S().Error(err)
		}
	}
}

func (c *cache[T]) Get(key string) (T, bool) {
	var zero T
	val, err := c.c.Get(key)
	if err != nil {
		// if not found, not log error
		return zero, false
	}
	var result T
	err = json.Unmarshal(val, &result)
	if err != nil {
		zap.S().Error(err)
		return zero, false
	}
	return result, true
}

func (c *cache[T]) Peek(key string) (T, bool) {
	return c.Get(key)
}
func (c *cache[T]) Remove(key string) { c.c.Delete(key) }
func (c *cache[T]) Exists(key string) bool {
	_, err := c.c.Get(key)
	return err == nil
}

func (c *cache[T]) Keys() []string {
	keys := make([]string, 0)
	iterator := c.c.Iterator()
	for iterator.SetNext() {
		entry, err := iterator.Value()
		if err != nil {
			continue
		}
		keys = append(keys, entry.Key())
	}
	return keys
}

func (c *cache[T]) Count() int {
	count := 0
	iterator := c.c.Iterator()
	for iterator.SetNext() {
		_, err := iterator.Value()
		if err == nil {
			count++
		}
	}
	return count
}
func (c *cache[T]) Flush() { c.c.Reset() }

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

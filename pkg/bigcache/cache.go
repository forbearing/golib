package bigcache

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	pkgzap "github.com/forbearing/golib/logger/zap"
)

var (
	globalCache  *bigcache.BigCache // globalCache is a global cache used by the entire application.
	sessionCache *bigcache.BigCache // sessionCache is a cache used by session store.
	once         sync.Once

	defaultExpire = time.Hour
)

func Init() (err error) {
	if globalCache == nil {
		if globalCache, err = bigcache.NewBigCache(bigcache.Config{
			Shards:           1024,
			LifeWindow:       0, // globalCache no expiration
			CleanWindow:      0,
			MaxEntrySize:     1024, // 1KB
			HardMaxCacheSize: 0,
			Verbose:          true,
			OnRemove:         nil,
			Logger:           pkgzap.NewStdLog(),
		}); err != nil {
			return err
		}
	}

	if sessionCache == nil {
		if sessionCache, err = bigcache.NewBigCache(bigcache.Config{
			Shards:             1024,
			LifeWindow:         defaultExpire,
			CleanWindow:        time.Second,
			MaxEntriesInWindow: 1000 * 10 * 60,
			MaxEntrySize:       500,
			Verbose:            true,
			HardMaxCacheSize:   0,
			Logger:             pkgzap.NewStdLog(),
		}); err != nil {
			return err
		}
	}
	return nil
}

type Cache struct {
	cache *bigcache.BigCache
	model reflect.Type
}

// New creates a globalCache object.
// mode could be a pointer or not pointer.
// eg: New(&service{}) or New(service{})
//
// more usage see test case.
func New(model any) *Cache {
	if reflect.ValueOf(model).Kind() == reflect.Pointer {
		return &Cache{
			cache: globalCache,
			model: reflect.TypeOf(model).Elem(),
		}
	}
	return &Cache{
		cache: globalCache,
		model: reflect.TypeOf(model),
	}
}

func (c *Cache) Len() int                { return c.cache.Len() }
func (c *Cache) Cap() int                { return c.cache.Capacity() }
func (c *Cache) Reset() error            { return c.cache.Reset() }
func (c *Cache) Delete(key string) error { return c.cache.Delete(key) }

// Get get entry from cache.
// NOTE: Get always returns pointer and always not nil, eg: &service{}, new(string).
// more detail see test cases.
func (c *Cache) Get(key string) (any, error) {
	entry, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	i := (reflect.New(c.model)).Interface()
	if err := json.Unmarshal(entry, i); err != nil {
		return nil, err
	}
	return i, nil
}

// Set add entry to cache.
// entry could be struct or pointer, eg: Set(key, &Service) or Set(key, Service{}).
// more detail se test cases.
func (c *Cache) Set(key string, entry any) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return c.cache.Set(key, data)
}

func (c *Cache) Close() (err error) {
	once.Do(func() {
		err = globalCache.Close()
	})
	return
}

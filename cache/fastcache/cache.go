package fastcache

import (
	"reflect"
	"sync"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.uber.org/zap"
)

var (
	cacheMap = cmap.New[any]()
	mu       sync.Mutex
)

func Init() (err error) {
	return nil
}

type cache[T any] struct {
	c *fastcache.Cache
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
		val = &cache[T]{c: fastcache.New(config.App.Cache.Capacity)}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, _ time.Duration) {
	val, err := util.Marshal(value)
	if err != nil {
		zap.S().Error(err)
	} else {
		c.c.Set([]byte(key), val)
	}
}

func (c *cache[T]) Get(key string) (T, bool) {
	var zero T
	value, ok := c.c.HasGet(nil, []byte(key))
	if !ok {
		return zero, false
	}
	var result T
	if err := util.Unmarshal(value, &result); err != nil {
		zap.S().Error(err)
		return zero, false
	}
	return result, true
}

func (c *cache[T]) Peek(key string) (T, bool) { return c.Get(key) }
func (c *cache[T]) Delete(key string)         { c.c.Del([]byte(key)) }

func (c *cache[T]) Exists(key string) bool { return c.c.Has([]byte(key)) }
func (c *cache[T]) Len() int               { return 0 }
func (c *cache[T]) Clear()                 { c.c.Reset() }

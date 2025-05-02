package smap

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/forbearing/golib/types"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var cacheMap = cmap.New[any]()

func Init() error {
	return nil
}

type cache[T any] struct {
	m sync.Map
	n int64
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	if !exists {
		val = &cache[T]{m: sync.Map{}}
		cacheMap.Set(key, val)
	}
	return val.(*cache[T])
}

func (c *cache[T]) Set(key string, value T, _ time.Duration) {
	_, loaded := c.m.LoadOrStore(key, value)
	if loaded {
		c.m.Store(key, value)
	} else {
		atomic.AddInt64(&c.n, 1)
	}
}

func (c *cache[T]) Get(key string) (T, bool) {
	v, ok1 := c.m.Load(key)
	_v, ok2 := v.(T)
	return _v, ok1 && ok2
}

func (c *cache[T]) Peek(key string) (T, bool) {
	v, ok1 := c.m.Load(key)
	_v, ok2 := v.(T)
	return _v, ok1 && ok2
}

func (c *cache[T]) Delete(key string) {
	_, exists := c.m.LoadAndDelete(key)
	if exists {
		atomic.AddInt64(&c.n, -1)
	}
}

func (c *cache[T]) Exists(key string) bool {
	_, exists := c.m.Load(key)
	return exists
}

func (c *cache[T]) Keys() []string {
	keys := make([]string, 0, 1024)
	c.m.Range(func(key, _ any) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}
func (c *cache[T]) Len() int { return int(c.n) }
func (c *cache[T]) Flush() {
	c.m.Range(func(key, _ any) bool {
		c.m.Delete(key)
		return true
	})
	atomic.StoreInt64(&c.n, 0)
}

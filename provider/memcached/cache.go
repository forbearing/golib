package memcached

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/forbearing/golib/types"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var _ types.Cache[any] = (*cache[any])(nil)

type cache[T any] struct{}

func Cache[T any]() types.Cache[T] {
	return new(cache[T])
}

func (c *cache[T]) Set(key string, value T, ttl time.Duration) {
	if !initialized {
		zap.S().Warn("memcached not initialized")
		return
	}
	val, err := json.Marshal(value)
	if err != nil {
		zap.S().Error(err)
		return
	}
	var exp int32
	if ttl <= 0 {
		exp = 0
	} else if ttl < 30*24*time.Hour {
		exp = int32(ttl.Seconds())
	} else {
		exp = int32(time.Now().Add(ttl).Unix())
	}
	if err := client.Set(&memcache.Item{
		Key:        key,
		Value:      val,
		Expiration: exp,
	}); err != nil {
		zap.S().Error(err)
	}
}

func (c *cache[T]) Get(key string) (T, bool) {
	var zero T
	if !initialized {
		zap.S().Warn("memcached not initialized")
		return zero, false
	}
	item, err := client.Get(key)
	if err != nil {
		return zero, false
	}
	var result T
	err = json.Unmarshal(item.Value, &result)
	if err != nil {
		zap.S().Error(err)
		return zero, false
	}
	return result, true
}
func (c *cache[T]) Peek(key string) (T, bool) { return c.Get(key) }
func (c *cache[T]) Delete(key string) {
	if !initialized {
		zap.S().Warn("memcached not initialized")
		return
	}
	if err := client.Delete(key); err != nil {
		zap.S().Error(err)
	}
}

func (c *cache[T]) Exists(key string) bool {
	if !initialized {
		zap.S().Warn("memcached not initialized")
		return false
	}
	_, err := client.Get(key)
	return err == nil
}

func (c *cache[T]) Len() int {
	if !initialized {
		zap.S().Warn("memcached not initialized")
		return 0
	}
	// NOTE: memcached don't support.
	return 0
}

func (c *cache[T]) Clear() {
	if !initialized {
		zap.S().Warn("memcached not initialized")
		return
	}
	if err := client.FlushAll(); err != nil {
		zap.S().Error(err)
	}
}

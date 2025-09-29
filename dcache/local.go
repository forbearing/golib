package dcache

import (
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.uber.org/zap/zapcore"
)

var (
	// 为什么选择 cmap v2
	//  1. sync.Map 不支持泛型, 在大量使用泛型的缓存库里面不使用泛型很突兀/麻烦
	//  2. cmap v2 比 sync.Map 性能要高很多
	localCacheMap = cmap.New[any]()
	localCacheMu  sync.Mutex
	localMaxItems = 1 << 24
)
var _ cacheMetricsProvider = (*localCache[any])(nil)

// localCache implements interface Cache use *ristretto as the backend memory localCache.
type localCache[T any] struct {
	c *ristretto.Cache[string, T]
}

// NewLocalCache 创建的缓存不具备分布式的能力, 需要分布式缓存请使用 NewDistributedCache
func NewLocalCache[T any]() (Cache[T], error) {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()

	// Fast path: check if cache already exists
	val, exists := localCacheMap.Get(key)
	if exists {
		return val.(*localCache[T]), nil
	}

	localCacheMu.Lock()
	defer localCacheMu.Unlock()

	// Double-check after acquiring lock
	val, exists = localCacheMap.Get(key)
	if !exists {
		c, err := ristretto.NewCache(buildConf[T]())
		if err != nil {
			return nil, err
		}
		val = &localCache[T]{c: c}
		localCacheMap.Set(key, val)
	}
	return val.(*localCache[T]), nil
}

func (c *localCache[T]) Set(key string, value T, ttl time.Duration) error {
	if success := c.c.SetWithTTL(key, value, 1, ttl); !success {
		return errors.New("cache rejected the set operation")
	}
	// Block here until value to be set.
	c.c.Wait()
	return nil
}

func (c *localCache[T]) Get(key string) (T, error) {
	val, ok := c.c.Get(key)
	if !ok {
		var zero T
		return zero, ErrEntryNotFound
	}
	return val, nil
}

// Delete removes the item with the provided key from the cache.
// It always returns nil as the underlying cache implementation doesn't
// provide information about whether the key existed or the deletion succeeded.
func (c *localCache[T]) Delete(key string) error {
	c.c.Del(key)
	return nil
}

func (c *localCache[T]) Exists(key string) bool {
	_, exists := c.c.Get(key)
	return exists
}

func (c *localCache[T]) Metrics() *localMetrics {
	m := c.c.Metrics
	return &localMetrics{
		Misses:       m.Misses(),
		KeysAdded:    m.KeysAdded(),
		KeysUpdated:  m.KeysUpdated(),
		KeysEvicted:  m.KeysEvicted(),
		CostAdded:    m.CostAdded(),
		CostEvicted:  m.CostEvicted(),
		SetsDropped:  m.SetsDropped(),
		SetsRejected: m.SetsRejected(),
		GetsDropped:  m.GetsDropped(),
		GetsKept:     m.GetsKept(),
		Ratio:        m.Ratio(),
	}
}

type localMetrics struct {
	Misses       uint64
	KeysAdded    uint64
	KeysUpdated  uint64
	KeysEvicted  uint64
	CostAdded    uint64
	CostEvicted  uint64
	SetsDropped  uint64
	SetsRejected uint64
	GetsDropped  uint64
	GetsKept     uint64
	Ratio        float64
}

func (m *localMetrics) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if m == nil {
		return nil
	}

	enc.AddUint64("misses", m.Misses)
	enc.AddUint64("keys_added", m.KeysAdded)
	enc.AddUint64("keys_updated", m.KeysUpdated)
	enc.AddUint64("keys_evicted", m.KeysEvicted)
	enc.AddUint64("cost_added", m.CostAdded)
	enc.AddUint64("cost_evicted", m.CostEvicted)
	enc.AddUint64("sets_dropped", m.SetsDropped)
	enc.AddUint64("sets_rejected", m.SetsRejected)
	enc.AddUint64("gets_dropped", m.GetsDropped)
	enc.AddUint64("gets_kept", m.GetsKept)
	enc.AddFloat64("ratio", m.Ratio)

	return nil
}

func buildConf[T any]() *ristretto.Config[string, T] {
	return &ristretto.Config[string, T]{
		// NumCounters 应该是你预期缓存项最大数量的大约 10 倍
		// 这个值影响内部布隆过滤器的准确性
		NumCounters: int64(localMaxItems) * 10,

		// MaxCost 就是你想要缓存的最大项数
		// 因为每个项的 cost 都是 1
		MaxCost: int64(localMaxItems),

		// BufferItems 控制写缓冲区大小
		// 我把缓存的个数设置为 1千多万, 所以这里设置的大一些
		BufferItems: 256,

		// 开启指标收集
		Metrics: true,
	}
}

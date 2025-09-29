package dcache

import (
	"time"
)

type Cache[T any] interface {
	Set(key string, value T, ttl time.Duration) error
	Get(key string) (T, error)
	Delete(key string) error
	Exists(key string) bool
}

type DistributedCache[T any] interface {
	// SetWithSync 原理
	//	1.设置本地缓存, 过期时间为 localTTL
	//	2.发送 `Set` 事件到状态节点
	//	3.状态节点设置 redis 缓存, 过期时间为 remoteTTL (Cache.Set 方法不会设置 redis 缓存)
	//	4.状态节点发送 `SetDone` 事件
	// 	5.当前节点更新本地缓存
	SetWithSync(key string, value T, localTTL time.Duration, remoteTTL time.Duration) error

	// GetWithSync 原理
	//	1.从本地缓存获取
	//	2.如果本地缓存不存在, 则从 redis 中获取
	//	3.如果从 redis 中获取到了则回填到本地缓存,过期时间为 localTTL.
	//	  回填本地缓存并不会发送 `Set` 事件到状态节点
	GetWithSync(key string, localTTL time.Duration) (T, error)

	// DeleteWithSync 原理
	//	1.从本地缓存删除
	//	2.发送 `Del` 事件到状态节点
	//	3.状态节点删除 redis 缓存 (Cache.Delete 方法不会删除 redis 缓存)
	//	4.状态节点发送 `DelDone` 事件
	//	5.当前节点删除本地缓存
	DeleteWithSync(key string) error

	Cache[T]
}

type cacheMetricsProvider interface {
	Metrics() *localMetrics
}

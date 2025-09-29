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
	// SetWithSync stores a value in both local and distributed cache with synchronization.
	//
	// Operation flow:
	//	1. Set value in local cache with localTTL expiration
	//	2. Send 'Set' event to state node
	//	3. State node sets Redis cache with remoteTTL expiration (Cache.Set method does not set Redis cache)
	//	4. State node sends 'SetDone' event
	//	5. Current node updates local cache
	SetWithSync(key string, value T, localTTL time.Duration, remoteTTL time.Duration) error

	// GetWithSync retrieves a value from local cache first, then from distributed cache if not found.
	//
	// Operation flow:
	//	1. Retrieve from local cache
	//	2. If not found in local cache, retrieve from Redis
	//	3. If found in Redis, backfill to local cache with localTTL expiration
	//	   Note: Backfilling local cache does not send 'Set' event to state node
	GetWithSync(key string, localTTL time.Duration) (T, error)

	// DeleteWithSync removes a value from both local and distributed cache with synchronization.
	//
	// Operation flow:
	//	1. Delete from local cache
	//	2. Send 'Del' event to state node
	//	3. State node deletes Redis cache (Cache.Delete method does not delete Redis cache)
	//	4. State node sends 'DelDone' event
	//	5. Current node deletes from local cache
	DeleteWithSync(key string) error

	Cache[T]
}

type cacheMetricsProvider interface {
	Metrics() *localMetrics
}

package dcache

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

const (
	MIN_GOROUTINES = 10000
	// TOPIC_REDIS_SET_DEL is the topic name to publish the entry associated with the key should update/delete event.
	TOPIC_REDIS_SET_DEL = "core-distributed-cache-set-del"
	// TOPIC_REDIS_DONE is the topic name to receive the entry associated with the key was update/delete event,
	// We should update/delete the local cache now.
	TOPIC_REDIS_DONE = "core-distributed-cache-done"

	GROUP_REDIS_SET_DEL = "core-distributed-cache-set-del"
	GROUP_REDIS_DONE    = "core-distributed-cache-done"
)

var (
	DistributedRedisCli   redis.UniversalClient // 被分布式缓存的协调节点使用到
	DistributedRedisCache Cache[any]            // 被分布式缓存的 GetWithSync 使用到
	redisPoolSize         = 500                 // 被 newRedis 使用, 测试用例需要用到 newRedis
	redisMaxIdleConns     = 100                 // 被 newRedis 使用, 测试用例需要用到 newRedis
)

var ErrEntryNotFound = errors.New("entry not found")

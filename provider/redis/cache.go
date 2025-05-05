package redis

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/types"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var _ types.Cache[any] = (*cache[any])(nil)

type cache[T any] struct{}

func Cache[T any]() types.Cache[T] {
	return new(cache[T])
}

func (*cache[T]) Set(key string, data T, expiration time.Duration) {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return
	}
	val, err := json.Marshal(data)
	if err != nil {
		zap.S().Error(err)
		return
	}
	if err := cli.Set(context.Background(), key, val, expiration).Err(); err != nil {
		zap.S().Error(err)
	}
}

func (*cache[T]) Get(key string) (T, bool) {
	var zero T
	if !initialized {
		zap.S().Warn("redis not initialized")
		return zero, false
	}
	val, err := cli.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return zero, false
		} else {
			zap.S().Error(err)
		}
		return zero, false
	}
	var result T
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		zap.S().Error(err)
		return zero, false
	}
	return result, true
}

func (c *cache[T]) Peek(key string) (T, bool) {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return *new(T), false
	}
	return c.Get(key)
}

func (*cache[T]) Delete(key string) {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return
	}
	if _, err := cli.Del(context.Background(), key).Result(); err != nil {
		zap.S().Error(err)
	}
}

func (*cache[T]) Exists(key string) bool {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return false
	}
	res, err := cli.Exists(context.Background(), key).Result()
	return err == nil && res > 0
}

func (*cache[T]) Len() int {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return 0
	}
	// In Redis Cluster, this only counts the selected node.
	count, err := cli.DBSize(context.Background()).Result()
	if err != nil {
		zap.S().Error(err)
		return 0
	}
	return int(count)
}

func (*cache[T]) Clear() {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return
	}
	if _, err := cli.FlushAll(context.Background()).Result(); err != nil {
		zap.S().Error(err)
	}
}

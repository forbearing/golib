package redis

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/cache/tracing"
	"github.com/forbearing/golib/types"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var _ types.Cache[any] = (*cache[any])(nil)

type cache[T any] struct {
	ctx context.Context
}

func Cache[T any]() types.Cache[T] {
	return tracing.NewTracingWrapper(&cache[T]{ctx: context.Background()}, "redis")
}

func (c *cache[T]) Set(key string, data T, ttl time.Duration) error {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return errors.New("redis not initialized")
	}
	val, err := json.Marshal(data)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	if err := cli.Set(c.ctx, key, val, ttl).Err(); err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}

func (c *cache[T]) Get(key string) (T, error) {
	var zero T
	if !initialized {
		zap.S().Warn("redis not initialized")
		return zero, errors.New("redis not initialized")
	}
	val, err := cli.Get(c.ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return zero, types.ErrEntryNotFound
		} else {
			zap.S().Error(err)
		}
		return zero, err
	}
	var result T
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		zap.S().Error(err)
		return zero, err
	}
	return result, nil
}

func (c *cache[T]) Peek(key string) (T, error) {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return *new(T), errors.New("redis not initialized")
	}
	return c.Get(key)
}

func (c *cache[T]) Delete(key string) error {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return errors.New("redis not initialized")
	}
	if _, err := cli.Del(c.ctx, key).Result(); err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}

func (c *cache[T]) Exists(key string) bool {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return false
	}
	res, err := cli.Exists(c.ctx, key).Result()
	if err != nil {
		return false
	}
	return res > 0
}

func (c *cache[T]) Len() int {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return 0
	}
	// In Redis Cluster, this only counts the selected node.
	count, err := cli.DBSize(c.ctx).Result()
	if err != nil {
		zap.S().Error(err)
		return 0
	}
	return int(count)
}

func (c *cache[T]) Clear() {
	if !initialized {
		zap.S().Warn("redis not initialized")
		return
	}
	if _, err := cli.FlushAll(c.ctx).Result(); err != nil {
		zap.S().Error(err)
	}
}

func (c *cache[T]) WithContext(ctx context.Context) types.Cache[T] { return &cache[T]{ctx: ctx} }

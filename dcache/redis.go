package dcache

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache implements CacheManager interface use Redis as the backend storage.
type redisCache[T any] struct {
	cli redis.UniversalClient // cli is Redis client.
	ctx context.Context       // ctx is global context used by the client.

	prefix string
}

// NewRedisCache creates CacheManager implementation that uses Redis as backend.
// It is your responsibility to ensure the redis client is valid.
func NewRedisCache[T any](ctx context.Context, cli redis.UniversalClient, opts ...RedisCacheOption[T]) (Cache[T], error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if cli == nil {
		return nil, errors.New("redis client is nil")
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := cli.Ping(pingCtx).Err(); err != nil {
		return nil, err
	}

	rc := &redisCache[T]{
		cli: cli,
		ctx: ctx,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(rc); err != nil {
			return nil, err
		}
	}
	return rc, nil
}

func (rc *redisCache[T]) Set(key string, value T, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if len(val) == 0 {
		return errors.New("cannot store empty value in redis")
	}
	return rc.cli.Set(rc.ctx, rc.prefix+key, val, ttl).Err()
}

func (rc *redisCache[T]) Get(key string) (T, error) {
	var zero T
	data, err := rc.cli.Get(rc.ctx, rc.prefix+key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return zero, ErrEntryNotFound
		}
		return zero, err
	}
	if len(data) == 0 {
		return zero, ErrEntryNotFound
	}
	var result T
	if err = json.Unmarshal(data, &result); err != nil {
		return zero, err
	}
	return result, nil
}

func (rc *redisCache[T]) Delete(key string) error {
	err := rc.cli.Del(rc.ctx, rc.prefix+key).Err()
	if errors.Is(err, redis.Nil) {
		return ErrEntryNotFound
	}
	return err
}

func (rc *redisCache[T]) Exists(key string) bool {
	res, err := rc.cli.Exists(rc.ctx, rc.prefix+key).Result()
	if err != nil {
		return false
	}
	return res > 0
}

// Range is not recommended for Redis implementation as it would require scanning all keys
func (rc *redisCache[T]) Range(f func(key string, value T) bool) {}

// RedisCacheOption is used to configure RedisCache.
type RedisCacheOption[T any] func(*redisCache[T]) error

func WithRedisKeyPrefix[T any](prefix string) RedisCacheOption[T] {
	return func(rc *redisCache[T]) error {
		rc.prefix = strings.TrimSpace(prefix)
		return nil
	}
}

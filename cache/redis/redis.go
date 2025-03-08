// 如果使用 Redis 版本小于等于 6，安装 v8 版本
// 如果使用 Redis 版本大于等于 7，安装 v9 版本
package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	// redis "github.com/redis/go-redis/v9"
	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/forbearing/golib/util"
	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	rdb         *redis.Client
	mu          sync.Mutex
	initialized bool

	ctx = context.Background()

	ErrKeyNotExists    = errors.New("key no loger exists, may be expired")
	ErrRedisIsDisabled = errors.New("redis is disabled")
)

// sonic library is about 2 times faster than standard library encoding/json.
// var json = sonic.ConfigStd

// func Init() (err error) {
// 	if !config.App.RedisConfig.Enable {
// 		return
// 	}
// 	once.Do(func() {
// 		rdb = redis.NewClient(&redis.Options{
// 			Addr:     net.JoinHostPort(config.App.RedisConfig.Host, strconv.Itoa(int(config.App.RedisConfig.Port))),
// 			Password: config.App.RedisConfig.Password,
// 			DB:       config.App.RedisConfig.DB,
// 		})
// 		// 确保 redis 没有问题
// 		err = rdb.Set(context.TODO(), config.App.RedisConfig.Namespace+"_"+"now", time.Now().Format(consts.DATE_TIME_LAYOUT), config.App.RedisConfig.Expiration).Err()
// 	})
// 	if err != nil {
// 		return errors.Wrap(err, "failed to connect to redis")
// 	}
// 	zap.S().Infow("successfully connect to redis", "host", config.App.RedisConfig.Host, "port", config.App.RedisConfig.Port, "db", config.App.RedisConfig.DB)
// 	return nil
// }

func Init() (err error) {
	cfg := config.App.RedisConfig
	if !cfg.Enable {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return nil
	}

	if rdb, err = New(cfg); err != nil {
		return errors.Wrap(err, "failed to connect to redis")
	}
	if err = rdb.Set(context.TODO(), cfg.Namespace+"_"+"now", time.Now().Format(consts.DATE_TIME_LAYOUT), cfg.Expiration).Err(); err != nil {
		return errors.Wrap(err, "failed to ping redis")
	}

	zap.S().Infow("successfully connect to redis", "host", cfg.Host, "port", cfg.Port, "db", cfg.DB)
	initialized = true
	return nil
}

func New(cfg config.RedisConfig) (*redis.Client, error) {
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(int(cfg.Port)))
	opts := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	if cfg.PoolSize > 0 {
		opts.PoolSize = cfg.PoolSize
	}
	if cfg.DialTimeout > 0 {
		opts.DialTimeout = cfg.DialTimeout
	}
	if cfg.ReadTimeout > 0 {
		opts.ReadTimeout = cfg.ReadTimeout
	}
	if cfg.WriteTimeout > 0 {
		opts.WriteTimeout = cfg.WriteTimeout
	}
	if cfg.MinIdleConns > 0 {
		opts.MinIdleConns = cfg.MinIdleConns
	}
	if cfg.MaxRetries > 0 {
		opts.MaxRetries = cfg.MaxRetries
	}
	if cfg.MinRetryBackoff > 0 {
		opts.MinRetryBackoff = cfg.MinRetryBackoff
	}
	if cfg.MaxRetryBackoff > 0 {
		opts.MaxRetryBackoff = cfg.MaxRetryBackoff
	}
	if cfg.IdleTimeout > 0 {
		opts.IdleTimeout = cfg.IdleTimeout
	}
	if cfg.MaxConnAge > 0 {
		opts.MaxConnAge = cfg.MaxConnAge
	}
	if cfg.EnableTLS {
		var tlsConfig *tls.Config
		var err error
		if tlsConfig, err = util.BuildTLSConfig(cfg.CertFile, cfg.KeyFile, cfg.CAFile, cfg.InsecureSkipVerify); err != nil {
			return nil, errors.Wrap(err, "failed to build TLS config")
		}
		opts.TLSConfig = tlsConfig
	}

	return redis.NewClient(opts), nil
}

// Set set any data into redis with specific key.
// If the data type is custom type or structure, you must implement the interface encoding.BinaryMarshaler.
func Set(key string, data any, expiration ...time.Duration) error {
	if !config.App.RedisConfig.Enable {
		zap.S().Warn(ErrRedisIsDisabled.Error())
		return nil
	}
	_expiration := config.App.RedisConfig.Expiration
	if len(expiration) > 0 {
		_expiration = expiration[0]
	}
	return rdb.Set(ctx, key, data, _expiration).Err()
}

// SetM set types.Model into redis with specific key.
func SetM[M types.Model](key string, m M, expiration ...time.Duration) error {
	if !config.App.RedisConfig.Enable {
		zap.S().Warn(ErrRedisIsDisabled.Error())
		return nil
	}
	_expiration := config.App.RedisConfig.Expiration
	if len(expiration) > 0 {
		_expiration = expiration[0]
	}
	return rdb.Set(ctx, key, modelMarshaler[M]{Model: m}, _expiration).Err()
}

// SetML set one or multiple types.Model into redis with specific key.
func SetML[M types.Model](key string, ml []M, expiration ...time.Duration) error {
	if !config.App.RedisConfig.Enable {
		zap.S().Warn(ErrRedisIsDisabled.Error())
		return nil
	}
	_expiration := config.App.RedisConfig.Expiration
	if len(expiration) > 0 {
		_expiration = expiration[0]
	}
	bl := make([]modelMarshaler[M], 0)
	for i := range ml {
		bl = append(bl, modelMarshaler[M]{Model: ml[i]})
	}
	return rdb.Set(ctx, key, modelMarshalerList[M](bl), _expiration).Err()
}

// SetSession wrapped Set function to set session data to redis, only for session.
func SetSession(sessionId string, data []byte) error {
	return Set(fmt.Sprintf("%s:session:%s", config.App.RedisConfig.Namespace, sessionId), data, config.App.AuthConfig.AccessTokenExpireDuration)
}

// Get will get raw cache([]byte) from redis.
func Get(key string) ([]byte, error) {
	if !config.App.RedisConfig.Enable {
		zap.S().Warn(ErrRedisIsDisabled.Error())
		return make([]byte, 0), nil
	}
	cache, err := rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotExists
	}
	return cache, nil
}

// GetInt get cache from redis and decode into integer.
func GetInt(key string) (int64, error) {
	if !config.App.RedisConfig.Enable {
		zap.S().Warn(ErrRedisIsDisabled.Error())
		return 0, nil
	}
	cache, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, ErrKeyNotExists
	}
	val, err := strconv.Atoi(cache)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

// GetM will get cache from redis and decode into types.Model.
func GetM[M types.Model](key string) (M, error) {
	if !config.App.RedisConfig.Enable {
		zap.S().Warn(ErrRedisIsDisabled.Error())
		return *new(M), nil
	}
	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return *new(M), ErrKeyNotExists
		}
		zap.S().Error(err)
		return *new(M), err
	}
	typ := reflect.TypeOf(*new(M)).Elem()
	val := reflect.New(typ).Interface().(M)
	if err := json.Unmarshal(data, val); err != nil {
		zap.S().Error(err)
		return *new(M), err
	}
	return val, nil
}

// GetML will get cache from redis and decode into []types.Model.
func GetML[M types.Model](key string) ([]M, error) {
	if !config.App.RedisConfig.Enable {
		zap.S().Warn(ErrRedisIsDisabled.Error())
		return make([]M, 0), nil
	}
	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrKeyNotExists
		}
		zap.S().Error(err)
		return nil, err
	}
	// typ := reflect.TypeOf(*new(M)).Elem()

	dest := make([]modelMarshaler[M], 0)
	if err := json.Unmarshal(data, &dest); err != nil {
		zap.S().Error(err)
		return nil, err
	}
	ml := make([]M, 0)
	for i := range dest {
		ml = append(ml, dest[i].Model)
	}
	return ml, nil
}

// Getsession wrapped Get function to get session data from cache
func GetSession[M types.Model](sessionId string) ([]byte, error) {
	key := fmt.Sprintf("%s:session:%s", config.App.RedisConfig.Namespace, sessionId)
	return Get(key)
}

func Remove(keys ...string) error {
	if !config.App.RedisConfig.Enable {
		zap.S().Warn(ErrRedisIsDisabled.Error())
		return nil
	}
	return rdb.Del(ctx, keys...).Err()
}

// RemovePrefix will scan and delete all redis key that matchs the `prefix`.
// for example: myprefix*
func RemovePrefix(prefix string) error {
	if !config.App.RedisConfig.Enable {
		zap.S().Warn(ErrRedisIsDisabled.Error())
		return nil
	}
	if !strings.HasSuffix(prefix, "*") {
		prefix = prefix + "*"
	}
	iter := rdb.Scan(ctx, 0, prefix, 0).Iterator()
	for iter.Next(ctx) {
		if err := rdb.Del(ctx, iter.Val()).Err(); err != nil {
			zap.S().Error(err)
			return err
		}
	}
	if err := iter.Err(); err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}

// modelMarshaler
// MarshalBinary, UnmarshalBinary 的 receiver 不能是指针, 否则 redis 会报错:
// redis: can't marshal redis.modelMarshaler[*myproject/model.FeishuUser] (implement encoding.BinaryMarshaler)
//
// MarshalJSON, UnmarshalJSON 的 receiver 必须是指针, 否则 panic
type modelMarshaler[M types.Model] struct {
	Model M
}

func (b modelMarshaler[M]) MarshalBinary() ([]byte, error) {
	return json.Marshal(b.Model)
	// buf := new(bytes.Buffer)
	// if err := gob.NewEncoder(buf).Encode(b.Model); err != nil {
	// 	zap.S().Error(err)
	// 	return nil, err
	// }
	// return buf.Bytes(), nil
}

// func (b modelMarshaler[M]) UnmarshalBinary(data []byte) error {
// 	return json.Unmarshal(data, b.Model)
// }

// func (b *modelMarshaler[M]) MarshalJSON() ([]byte, error) {
// 	data, err := json.Marshal(b.Model)
// 	if err != nil {
// 		zap.S().Error(err)
// 		return nil, err
// 	}
// 	return data, err
// }

func (b *modelMarshaler[M]) UnmarshalJSON(data []byte) error {
	if reflect.DeepEqual(b.Model, *new(M)) {
		b.Model = reflect.New(reflect.TypeOf(*new(M)).Elem()).Interface().(M)
	}
	if err := json.Unmarshal(data, &b.Model); err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}

// modelMarshalerList
// MarshalBinary 的 receiver一定不能是指针
type modelMarshalerList[M types.Model] []modelMarshaler[M]

func (bl modelMarshalerList[M]) MarshalBinary() ([]byte, error) {
	// ml := make([]types.Model, 0)
	// for i := range bl {
	// 	ml = append(ml, bl[i].Model)
	// }
	// return json.Marshal(ml)

	ml := make([]types.Model, len(bl))
	for i := range bl {
		ml[i] = bl[i].Model
	}
	return json.Marshal(ml)
}

// func (bl modelMarshalerList[M]) MarshalJSON() ([]byte, error) {
// 	ml := make([]types.Model, 0)
// 	for i := range bl {
// 		ml = append(ml, bl[i].Model)
// 	}
// 	return json.Marshal(ml)
// }
// func (bl *modelMarshalerList[M]) UnmarshalJSON(data []byte) error {
// 	bs := make([]modelMarshaler[M], 0)
// 	if err := json.Unmarshal(data, &bs); err != nil {
// 		zap.S().Error(err)
// 		return err
// 	}
// 	*bl = modelMarshalerList[M](bs)
// 	return nil
// }

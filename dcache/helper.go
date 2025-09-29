package dcache

import (
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap/zapcore"
)

type mapMarshaler map[string]int64

func (m mapMarshaler) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if m == nil {
		return nil
	}

	for k, v := range m {
		enc.AddInt64(k, v)
	}

	return nil
}

func calculateHitRatio(hits, misses int64) int64 {
	if hits+misses == 0 {
		return 0
	}
	return hits * 100 / (hits + misses)
}

func newProducer(brokers []string, topic string) (*kgo.Client, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.ClientID(fmt.Sprintf("producer-%s-%s", topic, hostname)),

		// 低延迟优化
		kgo.ProducerLinger(1*time.Millisecond), // 极短的批处理等待时间
		// kgo.ProducerBatchMaxBytes(n),           // 较小的批处理大小
		// kgo.MaxBufferedRecords(n),              // 大缓冲区以处理突发流量

		// 可靠性降级以换取更低延迟
		// 不需要消息幂等性, 状态节点会自动去重复和记录最大时间戳来保证 最终状态一致性
		// 本地环境下发现如下配置可以在每批次 operator 中减少100-200ms的延迟
		// kgo.RequiredAcks(kgo.NoAck()),
		// kgo.DisableIdempotentWrite(),           // 禁用幂等性以减少开销
		kgo.RetryTimeout(300*time.Millisecond), // 快速失败而不是长时间重试

		// TCP连接优化
		kgo.DialTimeout(300*time.Millisecond),     // 快速连接超时
		kgo.RequestTimeoutOverhead(1*time.Second), // 最小1s,否则kgo.NewClient 会报错
	)
}

// newConsumer 创建 kafka 消费者, 会有多个消费者
func newConsumer(brokers []string, topic string, group string) (*kgo.Client, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.ConsumeTopics(topic),
		kgo.ClientID(fmt.Sprintf("consumer-%s-%s", topic, hostname)),

		// 不需要自动提交, 也不需要手动提交, 系统每次重启之后使用最新的 offset
		kgo.DisableAutoCommit(),
		// 每次启动时,都是新的 group id
		kgo.ConsumerGroup(fmt.Sprintf("%s-%d", group, time.Now().UnixNano())),
		// 系统启动时,总是消费最新的消息
		kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()),

		// 低延迟消费优化
		kgo.FetchMaxWait(10*time.Millisecond), // 非常短的拉取等待时间
		kgo.FetchMinBytes(1),                  // 任何数据都立即返回
		// kgo.FetchMaxBytes(n),           // 较大的最大获取大小 (10MB)

		// TCP连接优化
		kgo.DialTimeout(300*time.Millisecond),
	)
}

// newRedis 仅在 benchmark 中使用
func newRedis() (redis.UniversalClient, error) {
	var client *redis.Client
	// TODO: 优先读取 redis 配置 ClusterUrl
	host := os.Getenv("LOCAL_REDIS_HOST")
	port := os.Getenv("LOCAL_REDIS_PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	redisUrl := fmt.Sprintf("redis://%s:%s/0", host, port)

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}
	opt.PoolSize = redisPoolSize
	opt.MaxIdleConns = redisMaxIdleConns
	client = redis.NewClient(opt)

	return client, nil
}

// newRedisOrDie 仅在 benchmark 中使用
func newRedisOrDie() redis.UniversalClient {
	client, err := newRedis()
	if err != nil {
		panic(err)
	}
	return client
}

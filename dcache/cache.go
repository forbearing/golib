package dcache

import (
	"context"
	"encoding/json"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/forbearing/golib/util"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/panjf2000/ants/v2"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

var once sync.Once

// Init 负责 Set/Delete redis key, 完成后再向 kafka 中推送消息
// 告诉多个分布式 core 节点删除本地二级缓存
//
// 关键规则实现:
//  1. 已记录事件的最大时间戳, 新收到的事件时间戳如果小于记录的最大时间戳, 则丢弃
//  2. 按时间戳去重：只保留每个键的最新操作
//     例如: 如果 Set(11:14) Delete(11:10), 则 Delete 并不会执行, 只会执行 Set(11:14)
//  3. 对不同的 key 进行时间戳排序, 严格按照时间戳顺序执行 redis 缓存操作 Set/Delete.
//  4. 记录当前批次事件中的最大时间戳
func Init() {
	once.Do(func() {
		const compKey = "comp"
		const compVal = "[DistributedCache.Setup]"

		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		zlog := zap.Must(zap.NewProduction())
		logger := zlog.With(zap.String("hostname", hostname), zap.String(compKey, compVal))
		logger.Info("distributed cache setup")

		redisCli := DistributedRedisCli
		if redisCli == nil {
			panic("DistributedRedisCli is nil")
		}

		var wg sync.WaitGroup

		// 手动通过线程池控制 kafka 并发量
		gopool, err := ants.NewPool(runtime.NumCPU()*2000, ants.WithPreAlloc(false))
		if err != nil {
			panic(err)
		}

		// 获取 Kafka 集群配置
		seeds := os.Getenv("KAFKA_BROKERS")
		var seedsArr []string
		for seed := range strings.SplitSeq(seeds, ",") {
			seed = strings.TrimSpace(seed)
			if len(seed) > 0 {
				seedsArr = append(seedsArr, seed)
			}
		}

		// 初始化 Kafka 消费者和生产者
		consumer, err := newConsumer(seedsArr, TOPIC_REDIS_SET_DEL, GROUP_REDIS_SET_DEL)
		if err != nil {
			panic(err)
		}
		producer, err := newProducer(seedsArr, TOPIC_REDIS_DONE)
		if err != nil {
			panic(err)
		}

		// 为每个 key 维护独立的最大时间戳
		keyMaxTimestamps := cmap.New[int64]()

		util.SafeGo(func() {
			for {
				// 基础上下文，用于操作超时控制
				baseCtx := context.Background()
				fetches := consumer.PollFetches(context.Background())
				if fetches.IsClientClosed() {
					logger.Error("fetches.IsClientClosed", zap.Error(err))
					continue
				}
				fetches.EachError(func(s string, i int32, err error) {
					logger.Error("failed to fetch from kafka",
						zap.Error(err),
						zap.String("topic", TOPIC_REDIS_SET_DEL),
						zap.String("s", s),
						zap.Int32("i", i),
					)
				})

				// 重置批次计数器
				var totalRecords int = 0     // 总消息数
				var successRecords int64 = 0 // 成功处理的消息数
				var failedRecords int64 = 0  // 处理失败的消息数
				var skippedRecords int = 0   // 跳过的无效的消息数

				// 用于跟踪本批次处理的消息的偏移量
				offsets := make(map[string]map[int32]kgo.EpochOffset)

				// ---------------------------------------------------------------------
				// 第一阶段：收集所有事件并按时间戳去重，保留每个键的最新操作
				// ---------------------------------------------------------------------

				// 存储每个键的最新操作，实现规则1和规则3
				keyEvents := make(map[string]*event)

				begin := time.Now()
				// 遍历所有分区的消息
				fetches.EachPartition(func(p kgo.FetchTopicPartition) {
					if len(p.Records) == 0 {
						return // 静默跳过空分区
					}

					totalRecords += len(p.Records)

					// 确保为每个主题初始化偏移量映射
					if _, exists := offsets[p.Topic]; !exists {
						offsets[p.Topic] = make(map[int32]kgo.EpochOffset)
					}

					var lastOffset int64 = -1
					for _, record := range p.Records {
						lastOffset = record.Offset // 记录最后一条消息的偏移量

						// 解析事件
						event := new(event)
						if err := json.Unmarshal(record.Value, event); err != nil {
							logger.Error("failed to unmarshal event from kafka record",
								zap.Error(err),
								zap.Int64("offset", record.Offset),
							)
							failedRecords++
							continue
						}

						// 获取该 key 的历史最大时间戳
						keyMaxTs, _ := keyMaxTimestamps.Get(event.Key)

						// 规则一：过滤掉时间戳小于该 key 历史最大时间戳的事件
						if event.TS <= keyMaxTs {
							logger.Warn("skipping outdated event for key",
								zap.String("key", event.Key),
								zap.Int64("event_ts", event.TS),
								zap.Int64("key_max_ts", keyMaxTs),
								zap.String("op", event.Op.String()),
							)
							skippedRecords++
							continue
						}

						// 规则二: 按时间戳去重：只保留每个键的最新操作
						existingEvent, exists := keyEvents[event.Key]
						if !exists || event.TS > existingEvent.TS {
							keyEvents[event.Key] = event
						}

					}

					// 更新分区偏移量，用于后续可能的手动提交偏移量(可能用不到了)
					if lastOffset >= 0 {
						offsets[p.Topic][p.Partition] = kgo.EpochOffset{
							Offset: lastOffset + 1,
							Epoch:  -1,
						}
					}
				})

				// 如果没有消息需要处理，则继续等待下一批
				if len(keyEvents) == 0 {
					logger.Debug("no events to process in this batch",
						zap.Int("total_records", totalRecords),
						zap.Int("skipped_records", skippedRecords),
						zap.Int64("failed_records", failedRecords),
					)
					continue
				}

				// 将map转换为切片，按照时间戳排序
				eventSlice := make([]*event, 0, len(keyEvents))
				for _, event := range keyEvents {
					eventSlice = append(eventSlice, event)
				}

				// 规则三: 严格按照时间戳排序 (从早到晚)
				sort.Slice(eventSlice, func(i, j int) bool {
					return eventSlice[i].TS < eventSlice[j].TS
				})

				// ---------------------------------------------------------------------
				// 第二阶段：按照时间戳顺序执行Redis操作, 操作完后推送 kafka 消息
				// ---------------------------------------------------------------------

				// 记录本批次处理的每个 key 的最大时间戳，用于批处理结束后更新
				batchKeyMaxTs := make(map[string]int64)

				// 批次操作 redis 和 kafka 超时控制
				wg.Add(len(eventSlice))
				for i := range eventSlice {
					evt := eventSlice[i]
					// 更新该 key 在本批次中的最大时间戳
					if ts, exists := batchKeyMaxTs[evt.Key]; !exists || evt.TS > ts {
						batchKeyMaxTs[evt.Key] = evt.TS
					}

					// TODO: 生产环境设置成 Debug 级别
					logger.Info("process event", zap.Object("event", evt))

					gopool.Submit(func() {
						defer wg.Done()
						switch evt.Op {
						case opSet:
							if evt.SyncToRedis {
								// logger.Info("redis set", zap.Int64("event_ts", evt.TS), zap.String("key", evt.Key), zap.Any("value", evt.Val), zap.Duration("redis_ttl", evt.RedisTTL))
								if err := redisCli.Set(baseCtx, evt.Key, []byte(evt.Val), evt.RedisTTL).Err(); err != nil {
									atomic.AddInt64(&failedRecords, 1)
									logger.Error("failed to set redis key",
										zap.Error(err),
										zap.String("key", evt.Key),
										zap.Object("event", evt),
									)
									return
								}
							}
							// 无论是否同步到Redis，都发送完成事件到Kafka
							evtDone := &event{
								CacheId:     evt.CacheId,
								Typ:         evt.Typ,
								Op:          opSetDone,
								Key:         evt.Key,
								Val:         evt.Val,
								TTL:         evt.TTL,
								TS:          time.Now().UnixNano(),
								Hostname:    evt.Hostname,
								SyncToRedis: evt.SyncToRedis,
								RedisTTL:    evt.RedisTTL,
							}
							data, err := json.Marshal(evtDone)
							if err != nil {
								logger.Error("failed to marshal event in redis set",
									zap.Error(err),
									zap.Object("event", evtDone),
								)
								atomic.AddInt64(&failedRecords, 1)
							} else {
								atomic.AddInt64(&successRecords, 1)
								// 同步推送 kafka 消息
								produceRecord := &kgo.Record{Topic: TOPIC_REDIS_DONE, Value: data}
								if err := producer.ProduceSync(baseCtx, produceRecord).FirstErr(); err != nil {
									logger.Error("failed to produce redis set done event",
										zap.Error(err),
										zap.Object("event", evtDone),
									)
								}
							}
						case opDel:
							if evt.SyncToRedis {
								if err := redisCli.Del(baseCtx, evt.Key).Err(); err != nil {
									logger.Error("failed to del redis key",
										zap.Error(err),
										zap.String("key", evt.Key),
										zap.Object("event", evt),
									)
									atomic.AddInt64(&failedRecords, 1)
									return
								}
							}
							// 无论是否同步到Redis，都发送完成事件到Kafka
							evtDone := &event{
								CacheId:     evt.CacheId,
								Typ:         evt.Typ,
								Op:          opDelDone,
								Key:         evt.Key,
								TS:          time.Now().UnixNano(),
								Hostname:    evt.Hostname,
								SyncToRedis: evt.SyncToRedis,
								RedisTTL:    evt.RedisTTL,
							}
							data, err := json.Marshal(evtDone)
							if err != nil {
								logger.Error("failed to marshal event in redis del",
									zap.Error(err),
									zap.Object("event", evtDone),
								)
								atomic.AddInt64(&failedRecords, 1)
							} else {
								atomic.AddInt64(&successRecords, 1)
								// 同步推送 kafka 消息
								produceRecord := &kgo.Record{Topic: TOPIC_REDIS_DONE, Value: data}
								if err := producer.ProduceSync(baseCtx, produceRecord).FirstErr(); err != nil {
									logger.Error("failed to produce redis del done event",
										zap.Error(err),
										zap.Object("event", evtDone),
									)
								}
							}
						default:
							logger.Warn("unknown operation type", zap.String("op", evt.Op.String()))
						}
					})
				}
				wg.Wait()

				// 批处理完成后，更新每个 key 的最大时间戳
				for key, ts := range batchKeyMaxTs {
					keyMaxTimestamps.Set(key, ts)
				}

				// 记录处理统计信息
				if totalRecords > 0 {
					logger.Info("successfully consumed events",
						zap.Int("total", totalRecords),
						zap.Int("deduplicated", len(eventSlice)),
						zap.Int64("success", successRecords),
						zap.Int64("failed", failedRecords),
						zap.Int("skipped", skippedRecords),
						zap.String("costed", util.FormatDurationSmart(time.Since(begin), 2)),
					)
				}

				// 清空 map 和 slice，帮助 GC 自动回收内存
				keyEvents = nil
				eventSlice = nil
				batchKeyMaxTs = nil

				// // 系统每次重启时，都会从最新的偏移量开始消费, 所以不需要保存偏移量
				// if len(offsets) > 0 {
				// 	consumer.CommitOffsets(ctx, offsets, func(c *kgo.Client, ocr1 *kmsg.OffsetCommitRequest, ocr2 *kmsg.OffsetCommitResponse, err error) {
				// 		if err != nil {
				// 			fmt.Println("failed to commit offsets:", err)
				// 		} else {
				// 			fmt.Printf("successfully committed offsets: total(%d), success(%d), failed(%d), offset(%v), costed(%s)\n",
				// 				totalRecords, successRecords, failedRecords, offsets, time.Since(begin).String())
				// 		}
				// 	})
				// }

			}
		}, "DistributedCache.Setup")

		// Output
		// 2:47PM INF successfully consumed records costed=977.84ms deduplicated=10000 failed=0 last_max_ts=1747810035544255000 new_max_ts=1747810035684441000 success=10000 total=56139
		// 2:48PM INF successfully consumed records costed=991.52ms deduplicated=10000 failed=0 last_max_ts=1747810035684441000 new_max_ts=1747810035790642000 success=10000 total=56846
		// 2:48PM INF successfully consumed records costed=1.03s deduplicated=10000 failed=0 last_max_ts=1747810035790642000 new_max_ts=1747810035948520000 success=10000 total=58564
		// 2:48PM INF successfully consumed records costed=1.03s deduplicated=10000 failed=0 last_max_ts=1747810035948520000 new_max_ts=1747810036242784000 success=10000 total=47102
		// 2:48PM INF successfully consumed records costed=1.07s deduplicated=10000 failed=0 last_max_ts=1747810036242784000 new_max_ts=1747810036628210000 success=10000 total=48791
		// 2:48PM INF successfully consumed records costed=1.02s deduplicated=9995 failed=0 last_max_ts=1747810036628210000 new_max_ts=1747810036808622000 success=9995 total=47263
		// 2:48PM INF successfully consumed records costed=1.03s deduplicated=9998 failed=0 last_max_ts=1747810036808622000 new_max_ts=1747810036995573000 success=9998 total=56503
		// 2:48PM INF successfully consumed records costed=1.03s deduplicated=9997 failed=0 last_max_ts=1747810036995573000 new_max_ts=1747810036995580000 success=9997 total=52768
		// 2:48PM INF successfully consumed records costed=305.42ms deduplicated=2827 failed=0 last_max_ts=1747810036995580000 new_max_ts=1747810055606577000 success=2827 total=40809
		// 2:48PM INF successfully consumed records costed=340.00ms deduplicated=3210 failed=0 last_max_ts=1747810055606577000 new_max_ts=1747810055715589000 success=3210 total=42127
		// 2:48PM INF successfully consumed records costed=346.26ms deduplicated=3245 failed=0 last_max_ts=1747810055715589000 new_max_ts=1747810055839879000 success=3245 total=40058
		// 2:48PM INF successfully consumed records costed=332.77ms deduplicated=3222 failed=0 last_max_ts=1747810055839879000 new_max_ts=1747810055941536000 success=3222 total=39228
		// 2:48PM INF successfully consumed records costed=342.09ms deduplicated=3146 failed=0 last_max_ts=1747810055941536000 new_max_ts=1747810056061429000 success=3146 total=39506
		// 2:48PM INF successfully consumed records costed=362.84ms deduplicated=3404 failed=0 last_max_ts=1747810056061429000 new_max_ts=1747810056170084000 success=3404 total=42165
		// 2:48PM INF successfully consumed records costed=353.88ms deduplicated=3348 failed=0 last_max_ts=1747810056170084000 new_max_ts=1747810056288620000 success=3348 total=38961
		// 2:48PM INF successfully consumed records costed=356.25ms deduplicated=3296 failed=0 last_max_ts=1747810056288620000 new_max_ts=1747810056388999000 success=3296 total=38570
		// 2:48PM INF successfully consumed records costed=371.59ms deduplicated=3321 failed=0 last_max_ts=1747810056388999000 new_max_ts=1747810056510415000 success=3321 total=40051
		// 2:48PM INF successfully consumed records costed=387.33ms deduplicated=3620 failed=0 last_max_ts=1747810056510415000 new_max_ts=1747810056629931000 success=3620 total=40695
		// 2:48PM INF successfully consumed records costed=392.28ms deduplicated=3648 failed=0 last_max_ts=1747810056629931000 new_max_ts=1747810056735707000 success=3648 total=39888
		// 2:48PM INF successfully consumed records costed=994.08ms deduplicated=10000 failed=0 last_max_ts=1747810056735707000 new_max_ts=1747810056882385000 success=10000 total=42071
		// 2:48PM INF successfully consumed records costed=990.28ms deduplicated=10000 failed=0 last_max_ts=1747810056882385000 new_max_ts=1747810058456780000 success=10000 total=47014
		// 2:48PM INF successfully consumed records costed=516.79ms deduplicated=5000 failed=0 last_max_ts=1747810058456780000 new_max_ts=1747810058568688000 success=5000 total=38861
		// 2:48PM INF successfully consumed records costed=332.00ms deduplicated=2873 failed=0 last_max_ts=1747810058568688000 new_max_ts=1747810058678092000 success=2873 total=37761
		// 2:48PM INF successfully consumed records costed=319.91ms deduplicated=2839 failed=0 last_max_ts=1747810058678092000 new_max_ts=1747810058766316000 success=2839 total=36631
		// 2:48PM INF successfully consumed records costed=315.16ms deduplicated=2967 failed=0 last_max_ts=1747810058766316000 new_max_ts=1747810058856516000 success=2967 total=36839
		// 2:48PM INF successfully consumed records costed=315.48ms deduplicated=2952 failed=0 last_max_ts=1747810058856516000 new_max_ts=1747810058964285000 success=2952 total=37050
		// 2:48PM INF successfully consumed records costed=339.21ms deduplicated=3100 failed=0 last_max_ts=1747810058964285000 new_max_ts=1747810059058019000 success=3100 total=37925
		// 2:48PM INF successfully consumed records costed=331.08ms deduplicated=3040 failed=0 last_max_ts=1747810059058019000 new_max_ts=1747810059170378000 success=3040 total=37750
		// 2:48PM INF successfully consumed records costed=336.70ms deduplicated=3132 failed=0 last_max_ts=1747810059170378000 new_max_ts=1747810059283774000 success=3132 total=38984
		// 2:48PM INF successfully consumed records costed=326.25ms deduplicated=3069 failed=0 last_max_ts=1747810059283774000 new_max_ts=1747810059377478000 success=3069 total=39299
		// 2:48PM INF successfully consumed records costed=289.65ms deduplicated=2737 failed=0 last_max_ts=1747810059377478000 new_max_ts=1747810059479912000 success=2737 total=35135
		// 2:48PM INF successfully consumed records costed=316.64ms deduplicated=2886 failed=0 last_max_ts=1747810059479912000 new_max_ts=1747810059572082000 success=2886 total=36912

		// 总结: 经过优化, 分布式缓存并发量极高的情况下也不会对 redis 造成太大的压力, 并且确保了分布式缓存的最终一致性
		// 更多测试结果请看 common/cache/ benchmark 结果
	})
}

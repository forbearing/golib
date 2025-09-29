package dcache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDistributedCacheBasicOperations 测试基本操作
func TestDistributedCacheBasicOperations(t *testing.T) {
	// 为了测试，我们需要替换一些依赖组件
	// 这里我们创建一个方法来获取测试用的distributedCache
	dc := setupTestDistributedCache[string](t)

	// 测试Set操作
	err := dc.Set("test-key", "test-value", 1*time.Minute)
	assert.NoError(t, err)

	// 本地缓存应该被设置
	val, err := dc.Get("test-key")
	assert.NoError(t, err)
	assert.Equal(t, "test-value", val)

	// 测试Delete操作
	err = dc.Delete("test-key")
	assert.NoError(t, err)

	// 应该不存在了
	assert.False(t, dc.Exists("test-key"))

	// 测试不存在的键
	_, err = dc.Get("non-existent")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrEntryNotFound))
}

// TestDistributedCacheWithSync 测试带同步的操作
func TestDistributedCacheWithSync(t *testing.T) {
	dc := setupTestDistributedCache[string](t)
	key, value := "test-key", "test-value"
	localTTL, remoteTTL := 500*time.Millisecond, 1*time.Minute

	// 测试SetWithSync
	err := dc.SetWithSync(key, value, localTTL, remoteTTL)
	assert.NoError(t, err)

	// 测试GetWithSync (从本地缓存获取)
	val, err := dc.Get("test-key")
	assert.NoError(t, err)
	assert.Equal(t, value, val)

	// 自动过期拿不到
	time.Sleep(localTTL)
	val, err = dc.Get("test-key")
	assert.Error(t, err, ErrEntryNotFound)
	assert.Equal(t, "", val)

	val, err = dc.GetWithSync(key, localTTL)
	assert.NoError(t, err)
	assert.Equal(t, value, val)

	// 主动删除 Delete
	dc.Delete(key)
	assert.NoError(t, err)
	val, err = dc.Get(key)
	assert.Error(t, err, ErrEntryNotFound)
	assert.Equal(t, "", val)

	val, err = dc.GetWithSync(key, localTTL)
	assert.NoError(t, err)
	assert.Equal(t, value, val)

	// 主动删除 DeleteWithSync
	dc.DeleteWithSync(key)
	assert.NoError(t, err)
	val, err = dc.Get(key)
	assert.Error(t, err, ErrEntryNotFound)
	assert.Equal(t, "", val)

	// 等状态节点删除 redis 中的 key
	time.Sleep(500 * time.Millisecond)
	val, err = dc.GetWithSync(key, localTTL)
	assert.Error(t, err, ErrEntryNotFound)
	assert.Equal(t, "", val)
}

// TestDistributedCacheTTL 测试TTL功能
func TestDistributedCacheTTL(t *testing.T) {
	dc := setupTestDistributedCache[string](t)

	// 设置非常短的TTL
	err := dc.Set("ttl-key", "ttl-value", 100*time.Millisecond)
	assert.NoError(t, err)

	// 立即应该能获取
	val, err := dc.Get("ttl-key")
	assert.NoError(t, err)
	assert.Equal(t, "ttl-value", val)

	// 等待TTL过期
	time.Sleep(200 * time.Millisecond)

	// 现在应该获取不到了
	_, err = dc.Get("ttl-key")
	assert.Error(t, err)
}

// TestDistributedCacheRemoteTTLValidation 测试RemoteTTL验证
func TestDistributedCacheRemoteTTLValidation(t *testing.T) {
	dc := setupTestDistributedCache[string](t)

	// 设置错误的TTL (remoteTTL < localTTL)
	err := dc.SetWithSync("invalid-ttl", "value", 2*time.Hour, 1*time.Hour)
	assert.Error(t, err)
}

// TestDistributedCacheConcurrency 测试并发操作
func TestDistributedCacheConcurrency(t *testing.T) {
	dc := setupTestDistributedCache[string](t)

	// 创建等待组来同步goroutines
	var wg sync.WaitGroup
	const numGoroutines = 100

	// 启动多个goroutines同时进行读写操作
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			key := fmt.Sprintf("concurrent-key-%d", idx)
			value := fmt.Sprintf("value-%d", idx)

			// 设置值
			err := dc.Set(key, value, 1*time.Minute)
			assert.NoError(t, err)

			// 读取值
			val, err := dc.Get(key)
			assert.NoError(t, err)
			assert.Equal(t, value, val)

			// 删除值
			err = dc.Delete(key)
			assert.NoError(t, err)
		}(i)
	}

	// 等待所有goroutines完成
	wg.Wait()
}

// TestDistributedCacheDifferentTypes 测试不同类型的缓存
func TestDistributedCacheDifferentTypes(t *testing.T) {
	// 字符串缓存
	strCache := setupTestDistributedCache[string](t)

	// 整数缓存
	intCache := setupTestDistributedCache[int](t)

	// 自定义结构体缓存
	type Person struct {
		Name string
		Age  int
	}
	personCache := setupTestDistributedCache[Person](t)

	// 测试各种类型操作
	err := strCache.Set("str", "string-value", 1*time.Minute)
	assert.NoError(t, err)

	err = intCache.Set("int", 42, 1*time.Minute)
	assert.NoError(t, err)

	err = personCache.Set("person", Person{Name: "Alice", Age: 30}, 1*time.Minute)
	assert.NoError(t, err)

	// 检查值
	strVal, err := strCache.Get("str")
	assert.NoError(t, err)
	assert.Equal(t, "string-value", strVal)

	intVal, err := intCache.Get("int")
	assert.NoError(t, err)
	assert.Equal(t, 42, intVal)

	personVal, err := personCache.Get("person")
	assert.NoError(t, err)
	assert.Equal(t, Person{Name: "Alice", Age: 30}, personVal)
}

// TestDistributedCacheLargeValues 测试大型值
func TestDistributedCacheLargeValues(t *testing.T) {
	dc := setupTestDistributedCache[string](t)

	// 创建一个大字符串
	largeValue := make([]byte, 1<<20) // 1MB
	for i := range largeValue {
		largeValue[i] = byte(i % 256)
	}
	largeString := string(largeValue)

	// 设置大值
	err := dc.Set("large", largeString, 1*time.Hour)
	assert.NoError(t, err)

	// 获取并验证
	val, err := dc.Get("large")
	assert.NoError(t, err)
	assert.Equal(t, largeString, val)
}

// TestDistributedCacheEdgeCases 测试边缘情况
func TestDistributedCacheEdgeCases(t *testing.T) {
	dc := setupTestDistributedCache[string](t)

	// 测试空键
	err := dc.Set("", "empty-key", 1*time.Hour)
	assert.NoError(t, err)
	val, err := dc.Get("")
	assert.NoError(t, err)
	assert.Equal(t, "empty-key", val)

	// 测试零TTL
	err = dc.Set("zero-ttl", "forever", 0)
	assert.NoError(t, err)

	// 测试极小TTL
	err = dc.Set("tiny-ttl", "quick", 1*time.Nanosecond)
	assert.NoError(t, err)
	time.Sleep(10 * time.Millisecond)
	_, err = dc.Get("tiny-ttl")
	assert.Error(t, err) // 应该已经过期

	// 测试极大TTL
	err = dc.Set("huge-ttl", "longterm", 100*365*24*time.Hour) // ~100年
	assert.NoError(t, err)
	val, err = dc.Get("huge-ttl")
	assert.NoError(t, err)
	assert.Equal(t, "longterm", val)
}

func setupTestDistributedCache[T any](t *testing.T) DistributedCache[T] {
	redisCache, err := NewRedisCache[any](context.TODO(), newRedisOrDie())
	if err != nil {
		t.Fatal(err)
	}
	distributed, err := NewDistributedCache(
		WithKafkaBrokers[T]([]string{"127.0.0.1:9092"}),
		WithRedisCache[T](redisCache),
	)
	if err != nil {
		t.Fatal(err)
	}

	return distributed
}

// func TestDistributedCache(t *testing.T) {
// 	redisCache, err := NewRedisCache[any](context.TODO(), newRedisOrDie())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	distributed, err := NewDistributedCache(
// 		WithKafkaBrokers[string]([]string{"127.0.0.1:9092"}),
// 		WithRedisCache[string](redisCache),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ttl := 1 * time.Minute
//
// 	// 只会在 redis 中设置一次,最终是 value5
// 	distributed.Set("key1", "value1", ttl)
// 	distributed.Set("key1", "value2", ttl)
// 	distributed.Set("key1", "value3", ttl)
// 	distributed.Set("key1", "value4", ttl)
// 	distributed.Set("key1", "value5", ttl)
//
// 	val, _ := distributed.Get("key1")
// 	fmt.Println(val)
// 	time.Sleep(1 * time.Second)
// 	val, _ = distributed.Get("key1")
// 	fmt.Println(val)
// 	distributed.Delete("key1")
//
// 	fmt.Println(distributed.GetWithSync("key1", ttl))
// 	fmt.Println(distributed.GetWithSync("key2", ttl))
//
// 	// 等待分布式缓存发送完事件
// 	time.Sleep(3 * time.Second)
// }

package cache

import (
	"context"

	"github.com/forbearing/golib/cache/bigcache"
	"github.com/forbearing/golib/cache/ccache"
	"github.com/forbearing/golib/cache/cmap"
	"github.com/forbearing/golib/cache/fastcache"
	"github.com/forbearing/golib/cache/freecache"
	"github.com/forbearing/golib/cache/gocache"
	"github.com/forbearing/golib/cache/lru"
	"github.com/forbearing/golib/cache/lrue"
	"github.com/forbearing/golib/cache/ristretto"
	"github.com/forbearing/golib/cache/smap"
	"github.com/forbearing/golib/cache/tracing"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
)

// Init initialize all cache implementations.
//
// # Cache Implementations Overview
//
// | Package     | Expiration Strategy       |
// |-------------|---------------------------|
// | lru         | No expiration             |
// | cmap        | No expiration             |
// | smap        | No expiration             |
// | fastcache   | No expiration             |
// | lrue        | Global expiration         |
// | bigcache    | Global expiration         |
// | ristretto   | Per-entry expiration      |
// | freecache   | Per-entry expiration      |
// | ccache      | Per-entry expiration      |
// | gocache     | Per-entry expiration      |
func Init() error {
	return util.CombineError(
		// ---- No expiration (eviction only by capacity or usage) ----
		lru.Init,
		cmap.Init,
		smap.Init,
		fastcache.Init,

		// ---- Global expiration (single TTL for all entries) ----
		lrue.Init,
		bigcache.Init,

		// ---- Per-entry expiration (each entry can have its own TTL) ----
		ristretto.Init,
		ccache.Init,
		gocache.Init,
		freecache.Init,
	)
}

func Cache[T any]() types.Cache[T]          { return lrue.Cache[T]() }
func ExpirableCache[T any]() types.Cache[T] { return ristretto.Cache[T]() }

// WithTracing wraps a cache with distributed tracing capabilities
func WithTracing[T any](cache types.Cache[T], cacheType string) *tracing.TracingWrapper[T] {
	return tracing.NewTracingWrapper(cache, cacheType)
}

// CacheWithTracing returns a default cache with tracing enabled
func CacheWithTracing[T any](ctx context.Context) types.Cache[T] {
	return Cache[T]().WithContext(ctx)
}

// ExpirableCacheWithTracing returns an expirable cache with tracing enabled
func ExpirableCacheWithTracing[T any](ctx context.Context) types.Cache[T] {
	return ExpirableCache[T]().WithContext(ctx)
}

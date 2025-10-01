package cache

import (
	"github.com/forbearing/gst/cache/bigcache"
	"github.com/forbearing/gst/cache/ccache"
	"github.com/forbearing/gst/cache/cmap"
	"github.com/forbearing/gst/cache/fastcache"
	"github.com/forbearing/gst/cache/freecache"
	"github.com/forbearing/gst/cache/gocache"
	"github.com/forbearing/gst/cache/lru"
	"github.com/forbearing/gst/cache/lrue"
	"github.com/forbearing/gst/cache/ristretto"
	"github.com/forbearing/gst/cache/smap"
	"github.com/forbearing/gst/types"
	"github.com/forbearing/gst/util"
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

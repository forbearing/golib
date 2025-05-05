package cache

import (
	"github.com/forbearing/golib/cache/bigcache"
	"github.com/forbearing/golib/cache/ccache"
	"github.com/forbearing/golib/cache/cmap"
	"github.com/forbearing/golib/cache/fastcache"
	"github.com/forbearing/golib/cache/freecache"
	"github.com/forbearing/golib/cache/lru"
	"github.com/forbearing/golib/cache/lrue"
	"github.com/forbearing/golib/cache/smap"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
)

func Init() error {
	return util.CombineError(
		lru.Init,
		lrue.Init,
		cmap.Init,
		smap.Init,
		fastcache.Init,
		bigcache.Init,
		freecache.Init,
		ccache.Init,
	)
}

func Cache[T any]() types.Cache[T] { return lrue.Cache[T]() }

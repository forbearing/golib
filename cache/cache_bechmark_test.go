package cache_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/forbearing/golib/cache/bigcache"
	"github.com/forbearing/golib/cache/cmap"
	"github.com/forbearing/golib/cache/freecache"
	"github.com/forbearing/golib/cache/lru"
	"github.com/forbearing/golib/cache/lrue"
	"github.com/forbearing/golib/cache/smap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/types"
)

type User struct {
	Name string `json:"name,omitempty"`
	model.Base
}

func benchInt(b *testing.B, cache types.Cache[int]) {
	b.Run("Get", func(b *testing.B) {
		for i := range b.N {
			cache.Set(fmt.Sprintf("key%d", i), i)
		}
		b.ResetTimer()
		for i := range b.N {
			cache.Get(fmt.Sprintf("key%d", i))
		}
	})
	b.Run("Set", func(b *testing.B) {
		for i := range b.N {
			cache.Set(fmt.Sprintf("key%d", i), i)
		}
	})
}

func benchString(b *testing.B, cache types.Cache[string]) {
	b.Run("Get", func(b *testing.B) {
		for i := range b.N {
			cache.Set(fmt.Sprintf("key%d", i), strconv.Itoa(i))
		}
		b.ResetTimer()
		for i := range b.N {
			cache.Get(fmt.Sprintf("key%d", i))
		}
	})
	b.Run("Set", func(b *testing.B) {
		for i := range b.N {
			cache.Set(fmt.Sprintf("key%d", i), strconv.Itoa(i))
		}
	})
}

func benchUser(b *testing.B, cache types.Cache[User]) {
	b.Run("Get", func(b *testing.B) {
		for i := range b.N {
			cache.Set(fmt.Sprintf("key%d", i), User{Name: fmt.Sprintf("user%d", i)})
		}
		b.ResetTimer()
		for i := range b.N {
			cache.Get(fmt.Sprintf("key%d", i))
		}
	})
	b.Run("Set", func(b *testing.B) {
		for i := range b.N {
			cache.Set(fmt.Sprintf("key%d", i), User{Name: fmt.Sprintf("user%d", i)})
		}
	})
}

func BenchmarkCache(b *testing.B) {
	if err := config.Init(); err != nil {
		b.Fatal(err)
	}
	if err := lru.Init(); err != nil {
		b.Fatal(err)
	}
	if err := cmap.Init(); err != nil {
		b.Fatal(err)
	}
	if err := smap.Init(); err != nil {
		b.Fatal(err)
	}
	if err := bigcache.Init(); err != nil {
		b.Fatal(err)
	}
	if err := freecache.Init(); err != nil {
		b.Fatal(err)
	}

	b.Run("int", func(b *testing.B) {
		b.Run("lru", func(b *testing.B) {
			benchInt(b, lru.Cache[int]())
		})
		b.Run("lrue", func(b *testing.B) {
			benchInt(b, lrue.Cache[int]())
		})
		b.Run("cmap", func(b *testing.B) {
			benchInt(b, cmap.Cache[int]())
		})
		b.Run("smap", func(b *testing.B) {
			benchInt(b, smap.Cache[int]())
		})
		b.Run("bigcache", func(b *testing.B) {
			benchInt(b, bigcache.Cache[int]())
		})
		b.Run("freecache", func(b *testing.B) {
			benchInt(b, freecache.Cache[int]())
		})
	})
	b.Run("string", func(b *testing.B) {
		b.Run("lru", func(b *testing.B) {
			benchString(b, lru.Cache[string]())
		})
		b.Run("lrue", func(b *testing.B) {
			benchString(b, lrue.Cache[string]())
		})
		b.Run("cmap", func(b *testing.B) {
			benchString(b, cmap.Cache[string]())
		})
		b.Run("smap", func(b *testing.B) {
			benchString(b, smap.Cache[string]())
		})
		b.Run("bigcache", func(b *testing.B) {
			benchString(b, bigcache.Cache[string]())
		})
		b.Run("freecache", func(b *testing.B) {
			benchString(b, freecache.Cache[string]())
		})
	})
	b.Run("model", func(b *testing.B) {
		b.Run("lru", func(b *testing.B) {
			benchUser(b, lru.Cache[User]())
		})
		b.Run("lrue", func(b *testing.B) {
			benchUser(b, lrue.Cache[User]())
		})
		b.Run("cmap", func(b *testing.B) {
			benchUser(b, cmap.Cache[User]())
		})
		b.Run("smap", func(b *testing.B) {
			benchUser(b, smap.Cache[User]())
		})
		b.Run("bigcache", func(b *testing.B) {
			benchUser(b, bigcache.Cache[User]())
		})
		b.Run("freecache", func(b *testing.B) {
			benchUser(b, freecache.Cache[User]())
		})
	})
}

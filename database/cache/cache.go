package cache

import (
	"reflect"

	"github.com/forbearing/golib/types"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var _ types.Cache[types.Model] = (*cache[types.Model])(nil)

// cache implements types.Cache interface.
type cache[M types.Model] struct {
	modelCmap  cmap.ConcurrentMap[string, []M]
	intCmap    cmap.ConcurrentMap[string, []int64]
	boolCmap   cmap.ConcurrentMap[string, []bool]
	floatCmap  cmap.ConcurrentMap[string, []float64]
	stringCmap cmap.ConcurrentMap[string, []string]
	anyCmap    cmap.ConcurrentMap[string, any]
}

var cacheMap = cmap.New[any]()

func Init() error {
	return nil
}

func (c *cache[M]) Set(key string, values ...M)            { c.modelCmap.Set(key, values) }
func (c *cache[M]) SetInt(key string, values ...int64)     { c.intCmap.Set(key, values) }
func (c *cache[M]) SetBool(key string, values ...bool)     { c.boolCmap.Set(key, values) }
func (c *cache[M]) SetFloat(key string, values ...float64) { c.floatCmap.Set(key, values) }
func (c *cache[M]) SetString(key string, values ...string) { c.stringCmap.Set(key, values) }
func (c *cache[M]) SetAny(key string, value any)           { c.anyCmap.Set(key, value) }

func (c *cache[M]) Get(key string) ([]M, bool)            { return c.modelCmap.Get(key) }
func (c *cache[M]) GetInt(key string) ([]int64, bool)     { return c.intCmap.Get(key) }
func (c *cache[M]) GetBool(key string) ([]bool, bool)     { return c.boolCmap.Get(key) }
func (c *cache[M]) GetFloat(key string) ([]float64, bool) { return c.floatCmap.Get(key) }
func (c *cache[M]) GetString(key string) ([]string, bool) { return c.stringCmap.Get(key) }
func (c *cache[M]) GetAny(key string) (any, bool)         { return c.anyCmap.Get(key) }

func (c *cache[M]) GetAll() map[string][]M            { return c.modelCmap.Items() }
func (c *cache[M]) GetAllInt() map[string][]int64     { return c.intCmap.Items() }
func (c *cache[M]) GetAllBool() map[string][]bool     { return c.boolCmap.Items() }
func (c *cache[M]) GetAllFloat() map[string][]float64 { return c.floatCmap.Items() }
func (c *cache[M]) GetAllString() map[string][]string { return c.stringCmap.Items() }
func (c *cache[M]) GetAllAny() map[string]any         { return c.anyCmap.Items() }

func (c *cache[M]) Peek(key string) ([]M, bool)            { return c.modelCmap.Get(key) }
func (c *cache[M]) PeekInt(key string) ([]int64, bool)     { return c.intCmap.Get(key) }
func (c *cache[M]) PeekBool(key string) ([]bool, bool)     { return c.boolCmap.Get(key) }
func (c *cache[M]) PeekFloat(key string) ([]float64, bool) { return c.floatCmap.Get(key) }
func (c *cache[M]) PeekString(key string) ([]string, bool) { return c.stringCmap.Get(key) }
func (c *cache[M]) PeekAny(key string) (any, bool)         { return c.anyCmap.Get(key) }

func (c *cache[M]) Remove(key string)       { c.modelCmap.Remove(key) }
func (c *cache[M]) RemoveInt(key string)    { c.intCmap.Remove(key) }
func (c *cache[M]) RemoveBool(key string)   { c.boolCmap.Remove(key) }
func (c *cache[M]) RemoveFloat(key string)  { c.floatCmap.Remove(key) }
func (c *cache[M]) RemoveString(key string) { c.stringCmap.Remove(key) }
func (c *cache[M]) RemoveAny(key string)    { c.anyCmap.Remove(key) }

func (c *cache[M]) Exists(key string) bool       { return c.modelCmap.Has(key) }
func (c *cache[M]) ExistsInt(key string) bool    { return c.intCmap.Has(key) }
func (c *cache[M]) ExistsBool(key string) bool   { return c.boolCmap.Has(key) }
func (c *cache[M]) ExistsFloat(key string) bool  { return c.floatCmap.Has(key) }
func (c *cache[M]) ExistsString(key string) bool { return c.stringCmap.Has(key) }
func (c *cache[M]) ExistsAny(key string) bool    { return c.anyCmap.Has(key) }

func (c *cache[M]) Keys() []string       { return c.modelCmap.Keys() }
func (c *cache[M]) KeysInt() []string    { return c.intCmap.Keys() }
func (c *cache[M]) KeysBool() []string   { return c.boolCmap.Keys() }
func (c *cache[M]) KeysFloat() []string  { return c.floatCmap.Keys() }
func (c *cache[M]) KeysString() []string { return c.stringCmap.Keys() }
func (c *cache[M]) KeysAny() []string    { return c.anyCmap.Keys() }

func (c *cache[M]) Flush() {
	c.modelCmap.Clear()
	c.intCmap.Clear()
	c.boolCmap.Clear()
	c.floatCmap.Clear()
	c.stringCmap.Clear()
	c.anyCmap.Clear()
}

func Cache[M types.Model]() types.Cache[M] {
	key := reflect.TypeOf(*new(M)).Elem().String()
	val, exists := cacheMap.Get(key)
	if !exists {
		val = &cache[M]{
			modelCmap:  cmap.New[[]M](),
			intCmap:    cmap.New[[]int64](),
			boolCmap:   cmap.New[[]bool](),
			floatCmap:  cmap.New[[]float64](),
			stringCmap: cmap.New[[]string](),
			anyCmap:    cmap.New[any](),
		}
		cacheMap.Set(key, val)
	}
	return val.(*cache[M])
}

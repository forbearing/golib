package cache

import (
	"reflect"

	"github.com/forbearing/golib/types"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var _ types.Cache[types.Model] = (*cache[types.Model])(nil)

// cache implements types.Cache interface.
type cache[M types.Model] struct {
	_model  cmap.ConcurrentMap[string, []M]
	_int    cmap.ConcurrentMap[string, []int64]
	_bool   cmap.ConcurrentMap[string, []bool]
	_float  cmap.ConcurrentMap[string, []float64]
	_string cmap.ConcurrentMap[string, []string]
	_any    cmap.ConcurrentMap[string, any]
}

var cacheMap = cmap.New[any]()

func Init() error {
	return nil
}

func (c *cache[M]) Set(key string, values ...M)            { c._model.Set(key, values) }
func (c *cache[M]) SetInt(key string, values ...int64)     { c._int.Set(key, values) }
func (c *cache[M]) SetBool(key string, values ...bool)     { c._bool.Set(key, values) }
func (c *cache[M]) SetFloat(key string, values ...float64) { c._float.Set(key, values) }
func (c *cache[M]) SetString(key string, values ...string) { c._string.Set(key, values) }
func (c *cache[M]) SetAny(key string, value any)           { c._any.Set(key, value) }

func (c *cache[M]) Get(key string) ([]M, bool)            { return c._model.Get(key) }
func (c *cache[M]) GetInt(key string) ([]int64, bool)     { return c._int.Get(key) }
func (c *cache[M]) GetBool(key string) ([]bool, bool)     { return c._bool.Get(key) }
func (c *cache[M]) GetFloat(key string) ([]float64, bool) { return c._float.Get(key) }
func (c *cache[M]) GetString(key string) ([]string, bool) { return c._string.Get(key) }
func (c *cache[M]) GetAny(key string) (any, bool)         { return c._any.Get(key) }

func (c *cache[M]) GetAll() map[string][]M            { return c._model.Items() }
func (c *cache[M]) GetAllInt() map[string][]int64     { return c._int.Items() }
func (c *cache[M]) GetAllBool() map[string][]bool     { return c._bool.Items() }
func (c *cache[M]) GetAllFloat() map[string][]float64 { return c._float.Items() }
func (c *cache[M]) GetAllString() map[string][]string { return c._string.Items() }
func (c *cache[M]) GetAllAny() map[string]any         { return c._any.Items() }

func (c *cache[M]) Peek(key string) ([]M, bool)            { return c._model.Get(key) }
func (c *cache[M]) PeekInt(key string) ([]int64, bool)     { return c._int.Get(key) }
func (c *cache[M]) PeekBool(key string) ([]bool, bool)     { return c._bool.Get(key) }
func (c *cache[M]) PeekFloat(key string) ([]float64, bool) { return c._float.Get(key) }
func (c *cache[M]) PeekString(key string) ([]string, bool) { return c._string.Get(key) }
func (c *cache[M]) PeekAny(key string) (any, bool)         { return c._any.Get(key) }

func (c *cache[M]) Remove(key string)       { c._model.Remove(key) }
func (c *cache[M]) RemoveInt(key string)    { c._int.Remove(key) }
func (c *cache[M]) RemoveBool(key string)   { c._bool.Remove(key) }
func (c *cache[M]) RemoveFloat(key string)  { c._float.Remove(key) }
func (c *cache[M]) RemoveString(key string) { c._string.Remove(key) }
func (c *cache[M]) RemoveAny(key string)    { c._any.Remove(key) }

func (c *cache[M]) Exists(key string) bool       { return c._model.Has(key) }
func (c *cache[M]) ExistsInt(key string) bool    { return c._int.Has(key) }
func (c *cache[M]) ExistsBool(key string) bool   { return c._bool.Has(key) }
func (c *cache[M]) ExistsFloat(key string) bool  { return c._float.Has(key) }
func (c *cache[M]) ExistsString(key string) bool { return c._string.Has(key) }
func (c *cache[M]) ExistsAny(key string) bool    { return c._any.Has(key) }

func (c *cache[M]) Keys() []string       { return c._model.Keys() }
func (c *cache[M]) KeysInt() []string    { return c._int.Keys() }
func (c *cache[M]) KeysBool() []string   { return c._bool.Keys() }
func (c *cache[M]) KeysFloat() []string  { return c._float.Keys() }
func (c *cache[M]) KeysString() []string { return c._string.Keys() }
func (c *cache[M]) KeysAny() []string    { return c._any.Keys() }

func (c *cache[M]) Flush() {
	c._model.Clear()
	c._int.Clear()
	c._bool.Clear()
	c._float.Clear()
	c._string.Clear()
	c._any.Clear()
}

func Cache[M types.Model]() types.Cache[M] {
	key := reflect.TypeOf(*new(M)).Elem().String()
	val, exists := cacheMap.Get(key)
	if !exists {
		val = &cache[M]{
			_model:  cmap.New[[]M](),
			_int:    cmap.New[[]int64](),
			_bool:   cmap.New[[]bool](),
			_float:  cmap.New[[]float64](),
			_string: cmap.New[[]string](),
			_any:    cmap.New[any](),
		}
		cacheMap.Set(key, val)
	}
	return val.(*cache[M])
}

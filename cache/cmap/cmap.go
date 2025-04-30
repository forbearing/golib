package cmap

import (
	"reflect"

	"github.com/forbearing/golib/types"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (
	_ types.Cache[types.Model] = (*_Cache[types.Model])(nil)
	_ types.Cache[int]         = _Int{}
	_ types.Cache[int8]        = _Int8{}
	_ types.Cache[int16]       = _Int16{}
	_ types.Cache[int32]       = _Int32{}
	_ types.Cache[int64]       = _Int64{}
	_ types.Cache[uint]        = _Uint{}
	_ types.Cache[uint8]       = _Uint8{}
	_ types.Cache[uint16]      = _Uint16{}
	_ types.Cache[uint32]      = _Uint32{}
	_ types.Cache[uint64]      = _Uint64{}
	_ types.Cache[float32]     = _Float32{}
	_ types.Cache[float64]     = _Float64{}
	_ types.Cache[bool]        = _Bool{}
	_ types.Cache[rune]        = _Rune{}
	_ types.Cache[string]      = _String{}
	_ types.Cache[byte]        = _Byte{}
	_ types.Cache[[]byte]      = _Bytes{}
	_ types.Cache[any]         = _Any{}
)

var cacheMap = cmap.New[any]()

var (
	Int     = _Int{}
	Int8    = _Int8{}
	Int16   = _Int16{}
	Int32   = _Int32{}
	Int64   = _Int64{}
	Uint    = _Uint{}
	Uint8   = _Uint8{}
	Uint16  = _Uint16{}
	Uint32  = _Uint32{}
	Uint64  = _Uint64{}
	Float32 = _Float32{}
	Float64 = _Float64{}
	Bool    = _Bool{}
	Rune    = _Rune{}
	String  = _String{}
	Byte    = _Byte{}
	Bytes   = _Bytes{}
	Any     = _Any{}

	_int     cmap.ConcurrentMap[string, int]
	_int8    cmap.ConcurrentMap[string, int8]
	_int16   cmap.ConcurrentMap[string, int16]
	_int32   cmap.ConcurrentMap[string, int32]
	_int64   cmap.ConcurrentMap[string, int64]
	_uint    cmap.ConcurrentMap[string, uint]
	_uint8   cmap.ConcurrentMap[string, uint8]
	_uint16  cmap.ConcurrentMap[string, uint16]
	_uint32  cmap.ConcurrentMap[string, uint32]
	_uint64  cmap.ConcurrentMap[string, uint64]
	_float32 cmap.ConcurrentMap[string, float32]
	_float64 cmap.ConcurrentMap[string, float64]
	_bool    cmap.ConcurrentMap[string, bool]
	_rune    cmap.ConcurrentMap[string, rune]
	_string  cmap.ConcurrentMap[string, string]
	_byte    cmap.ConcurrentMap[string, byte]
	_bytes   cmap.ConcurrentMap[string, []byte]
	_any     cmap.ConcurrentMap[string, any]
)

type _Cache[T any] struct {
	_cmap cmap.ConcurrentMap[string, T]
}
type (
	_Int     struct{}
	_Int8    struct{}
	_Int16   struct{}
	_Int32   struct{}
	_Int64   struct{}
	_Uint    struct{}
	_Uint8   struct{}
	_Uint16  struct{}
	_Uint32  struct{}
	_Uint64  struct{}
	_Float32 struct{}
	_Float64 struct{}
	_Bool    struct{}
	_String  struct{}
	_Byte    struct{}
	_Rune    struct{}
	_Bytes   struct{}
	_Any     struct{}
)

func Init() error {
	_int = cmap.New[int]()
	_int8 = cmap.New[int8]()
	_int16 = cmap.New[int16]()
	_int32 = cmap.New[int32]()
	_int64 = cmap.New[int64]()
	_uint = cmap.New[uint]()
	_uint8 = cmap.New[uint8]()
	_uint16 = cmap.New[uint16]()
	_uint32 = cmap.New[uint32]()
	_uint64 = cmap.New[uint64]()
	_bool = cmap.New[bool]()
	_float32 = cmap.New[float32]()
	_float64 = cmap.New[float64]()
	_string = cmap.New[string]()
	_rune = cmap.New[rune]()
	_byte = cmap.New[byte]()
	_bytes = cmap.New[[]byte]()
	_any = cmap.New[any]()
	return nil
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	if !exists {
		val = &_Cache[T]{_cmap: cmap.New[T]()}
		cacheMap.Set(key, val)
	}
	return val.(*_Cache[T])
}

func (_Int) Set(key string, value int)   { _int.Set(key, value) }
func (_Int) Get(key string) (int, bool)  { return _int.Get(key) }
func (_Int) Peek(key string) (int, bool) { return _int.Get(key) }
func (_Int) Remove(key string)           { _int.Remove(key) }
func (_Int) Exists(key string) bool      { return _int.Has(key) }
func (_Int) Keys() []string              { return _int.Keys() }
func (_Int) Count() int                  { return _int.Count() }
func (_Int) Flush()                      { _int.Clear() }

func (_Int8) Set(key string, value int8)   { _int8.Set(key, value) }
func (_Int8) Get(key string) (int8, bool)  { return _int8.Get(key) }
func (_Int8) Peek(key string) (int8, bool) { return _int8.Get(key) }
func (_Int8) Remove(key string)            { _int8.Remove(key) }
func (_Int8) Exists(key string) bool       { return _int8.Has(key) }
func (_Int8) Keys() []string               { return _int8.Keys() }
func (_Int8) Count() int                   { return _int8.Count() }
func (_Int8) Flush()                       { _int8.Clear() }

func (_Int16) Set(key string, value int16)   { _int16.Set(key, value) }
func (_Int16) Get(key string) (int16, bool)  { return _int16.Get(key) }
func (_Int16) Peek(key string) (int16, bool) { return _int16.Get(key) }
func (_Int16) Remove(key string)             { _int16.Remove(key) }
func (_Int16) Exists(key string) bool        { return _int16.Has(key) }
func (_Int16) Keys() []string                { return _int16.Keys() }
func (_Int16) Count() int                    { return _int16.Count() }
func (_Int16) Flush()                        { _int16.Clear() }

func (_Int32) Set(key string, value int32)   { _int32.Set(key, value) }
func (_Int32) Get(key string) (int32, bool)  { return _int32.Get(key) }
func (_Int32) Peek(key string) (int32, bool) { return _int32.Get(key) }
func (_Int32) Remove(key string)             { _int32.Remove(key) }
func (_Int32) Exists(key string) bool        { return _int32.Has(key) }
func (_Int32) Keys() []string                { return _int32.Keys() }
func (_Int32) Count() int                    { return _int32.Count() }
func (_Int32) Flush()                        { _int32.Clear() }

func (_Int64) Set(key string, value int64)   { _int64.Set(key, value) }
func (_Int64) Get(key string) (int64, bool)  { return _int64.Get(key) }
func (_Int64) Peek(key string) (int64, bool) { return _int64.Get(key) }
func (_Int64) Remove(key string)             { _int64.Remove(key) }
func (_Int64) Exists(key string) bool        { return _int64.Has(key) }
func (_Int64) Keys() []string                { return _int64.Keys() }
func (_Int64) Count() int                    { return _int64.Count() }
func (_Int64) Flush()                        { _int64.Clear() }

func (_Uint) Set(key string, value uint)   { _uint.Set(key, value) }
func (_Uint) Get(key string) (uint, bool)  { return _uint.Get(key) }
func (_Uint) Peek(key string) (uint, bool) { return _uint.Get(key) }
func (_Uint) Remove(key string)            { _uint.Remove(key) }
func (_Uint) Exists(key string) bool       { return _uint.Has(key) }
func (_Uint) Keys() []string               { return _uint.Keys() }
func (_Uint) Count() int                   { return _uint.Count() }
func (_Uint) Flush()                       { _uint.Clear() }

func (_Uint8) Set(key string, value uint8)   { _uint8.Set(key, value) }
func (_Uint8) Get(key string) (uint8, bool)  { return _uint8.Get(key) }
func (_Uint8) Peek(key string) (uint8, bool) { return _uint8.Get(key) }
func (_Uint8) Remove(key string)             { _uint8.Remove(key) }
func (_Uint8) Exists(key string) bool        { return _uint8.Has(key) }
func (_Uint8) Keys() []string                { return _uint8.Keys() }
func (_Uint8) Count() int                    { return _uint8.Count() }
func (_Uint8) Flush()                        { _uint8.Clear() }

func (_Uint16) Set(key string, value uint16)   { _uint16.Set(key, value) }
func (_Uint16) Get(key string) (uint16, bool)  { return _uint16.Get(key) }
func (_Uint16) Peek(key string) (uint16, bool) { return _uint16.Get(key) }
func (_Uint16) Remove(key string)              { _uint16.Remove(key) }
func (_Uint16) Exists(key string) bool         { return _uint16.Has(key) }
func (_Uint16) Keys() []string                 { return _uint16.Keys() }
func (_Uint16) Count() int                     { return _uint16.Count() }
func (_Uint16) Flush()                         { _uint16.Clear() }

func (_Uint32) Set(key string, value uint32)   { _uint32.Set(key, value) }
func (_Uint32) Get(key string) (uint32, bool)  { return _uint32.Get(key) }
func (_Uint32) Peek(key string) (uint32, bool) { return _uint32.Get(key) }
func (_Uint32) Remove(key string)              { _uint32.Remove(key) }
func (_Uint32) Exists(key string) bool         { return _uint32.Has(key) }
func (_Uint32) Keys() []string                 { return _uint32.Keys() }
func (_Uint32) Count() int                     { return _uint32.Count() }
func (_Uint32) Flush()                         { _uint32.Clear() }

func (_Uint64) Set(key string, value uint64)   { _uint64.Set(key, value) }
func (_Uint64) Get(key string) (uint64, bool)  { return _uint64.Get(key) }
func (_Uint64) Peek(key string) (uint64, bool) { return _uint64.Get(key) }
func (_Uint64) Remove(key string)              { _uint64.Remove(key) }
func (_Uint64) Exists(key string) bool         { return _uint64.Has(key) }
func (_Uint64) Keys() []string                 { return _uint64.Keys() }
func (_Uint64) Count() int                     { return _uint64.Count() }
func (_Uint64) Flush()                         { _uint64.Clear() }

func (_Float32) Set(key string, value float32)   { _float32.Set(key, value) }
func (_Float32) Get(key string) (float32, bool)  { return _float32.Get(key) }
func (_Float32) Peek(key string) (float32, bool) { return _float32.Get(key) }
func (_Float32) Remove(key string)               { _float32.Remove(key) }
func (_Float32) Exists(key string) bool          { return _float32.Has(key) }
func (_Float32) Keys() []string                  { return _float32.Keys() }
func (_Float32) Count() int                      { return _float32.Count() }
func (_Float32) Flush()                          { _float32.Clear() }

func (_Float64) Set(key string, value float64)   { _float64.Set(key, value) }
func (_Float64) Get(key string) (float64, bool)  { return _float64.Get(key) }
func (_Float64) Peek(key string) (float64, bool) { return _float64.Get(key) }
func (_Float64) Remove(key string)               { _float64.Remove(key) }
func (_Float64) Exists(key string) bool          { return _float64.Has(key) }
func (_Float64) Keys() []string                  { return _float64.Keys() }
func (_Float64) Count() int                      { return _float64.Count() }
func (_Float64) Flush()                          { _float64.Clear() }

func (_Bool) Set(key string, value bool)   { _bool.Set(key, value) }
func (_Bool) Get(key string) (bool, bool)  { return _bool.Get(key) }
func (_Bool) Peek(key string) (bool, bool) { return _bool.Get(key) }
func (_Bool) Remove(key string)            { _bool.Remove(key) }
func (_Bool) Exists(key string) bool       { return _bool.Has(key) }
func (_Bool) Keys() []string               { return _bool.Keys() }
func (_Bool) Count() int                   { return _bool.Count() }
func (_Bool) Flush()                       { _bool.Clear() }

func (_Rune) Set(key string, value rune)   { _rune.Set(key, value) }
func (_Rune) Get(key string) (rune, bool)  { return _rune.Get(key) }
func (_Rune) Peek(key string) (rune, bool) { return _rune.Get(key) }
func (_Rune) Remove(key string)            { _rune.Remove(key) }
func (_Rune) Exists(key string) bool       { return _rune.Has(key) }
func (_Rune) Keys() []string               { return _rune.Keys() }
func (_Rune) Count() int                   { return _rune.Count() }
func (_Rune) Flush()                       { _rune.Clear() }

func (_String) Set(key string, value string)   { _string.Set(key, value) }
func (_String) Get(key string) (string, bool)  { return _string.Get(key) }
func (_String) Peek(key string) (string, bool) { return _string.Get(key) }
func (_String) Remove(key string)              { _string.Remove(key) }
func (_String) Exists(key string) bool         { return _string.Has(key) }
func (_String) Keys() []string                 { return _string.Keys() }
func (_String) Count() int                     { return _string.Count() }
func (_String) Flush()                         { _string.Clear() }

func (_Byte) Set(key string, value byte)   { _byte.Set(key, value) }
func (_Byte) Get(key string) (byte, bool)  { return _byte.Get(key) }
func (_Byte) Peek(key string) (byte, bool) { return _byte.Get(key) }
func (_Byte) Remove(key string)            { _byte.Remove(key) }
func (_Byte) Exists(key string) bool       { return _byte.Has(key) }
func (_Byte) Keys() []string               { return _byte.Keys() }
func (_Byte) Count() int                   { return _byte.Count() }
func (_Byte) Flush()                       { _byte.Clear() }

func (_Bytes) Set(key string, value []byte)   { _bytes.Set(key, value) }
func (_Bytes) Get(key string) ([]byte, bool)  { return _bytes.Get(key) }
func (_Bytes) Peek(key string) ([]byte, bool) { return _bytes.Get(key) }
func (_Bytes) Remove(key string)              { _bytes.Remove(key) }
func (_Bytes) Exists(key string) bool         { return _bytes.Has(key) }
func (_Bytes) Keys() []string                 { return _bytes.Keys() }
func (_Bytes) Count() int                     { return _bytes.Count() }
func (_Bytes) Flush()                         { _bytes.Clear() }

func (_Any) Set(key string, value any)   { _any.Set(key, value) }
func (_Any) Get(key string) (any, bool)  { return _any.Get(key) }
func (_Any) Peek(key string) (any, bool) { return _any.Get(key) }
func (_Any) Remove(key string)           { _any.Remove(key) }
func (_Any) Exists(key string) bool      { return _any.Has(key) }
func (_Any) Keys() []string              { return _any.Keys() }
func (_Any) Count() int                  { return _any.Count() }
func (_Any) Flush()                      { _any.Clear() }

func (c *_Cache[M]) Set(key string, value M)   { c._cmap.Set(key, value) }
func (c *_Cache[M]) Get(key string) (M, bool)  { return c._cmap.Get(key) }
func (c *_Cache[M]) Peek(key string) (M, bool) { return c._cmap.Get(key) }
func (c *_Cache[M]) Remove(key string)         { c._cmap.Remove(key) }
func (c *_Cache[M]) Exists(key string) bool    { return c._cmap.Has(key) }
func (c *_Cache[M]) Keys() []string            { return c._cmap.Keys() }
func (c *_Cache[M]) Count() int                { return c._cmap.Count() }
func (c *_Cache[M]) Flush()                    { c._cmap.Clear() }

package lru

import (
	"reflect"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/types"
	lru "github.com/hashicorp/golang-lru/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
)

const DEFAULT_SIZE = 1000000

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

	_int     *lru.Cache[string, int]
	_int8    *lru.Cache[string, int8]
	_int16   *lru.Cache[string, int16]
	_int32   *lru.Cache[string, int32]
	_int64   *lru.Cache[string, int64]
	_uint    *lru.Cache[string, uint]
	_uint8   *lru.Cache[string, uint8]
	_uint16  *lru.Cache[string, uint16]
	_uint32  *lru.Cache[string, uint32]
	_uint64  *lru.Cache[string, uint64]
	_float32 *lru.Cache[string, float32]
	_float64 *lru.Cache[string, float64]
	_bool    *lru.Cache[string, bool]
	_rune    *lru.Cache[string, rune]
	_string  *lru.Cache[string, string]
	_byte    *lru.Cache[string, byte]
	_bytes   *lru.Cache[string, []byte]
	_any     *lru.Cache[string, any]
)

type _Cache[T any] struct {
	_lru *lru.Cache[string, T]
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

func Init() (err error) {
	if _int, err = lru.New[string, int](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _int8, err = lru.New[string, int8](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _int16, err = lru.New[string, int16](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _int32, err = lru.New[string, int32](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _int64, err = lru.New[string, int64](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _uint, err = lru.New[string, uint](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _uint8, err = lru.New[string, uint8](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _uint16, err = lru.New[string, uint16](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _uint32, err = lru.New[string, uint32](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _uint64, err = lru.New[string, uint64](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _float32, err = lru.New[string, float32](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _float64, err = lru.New[string, float64](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _bool, err = lru.New[string, bool](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _rune, err = lru.New[string, rune](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _string, err = lru.New[string, string](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _byte, err = lru.New[string, byte](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _bytes, err = lru.New[string, []byte](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	if _any, err = lru.New[string, any](DEFAULT_SIZE); err != nil {
		return errors.WithStack(err)
	}
	return
}

func Cache[T any]() types.Cache[T] {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	key := typ.PkgPath() + "|" + typ.String()
	val, exists := cacheMap.Get(key)
	// lru.New() only error on negative size.
	if !exists {
		_lru, _ := lru.New[string, T](DEFAULT_SIZE)
		val = &_Cache[T]{_lru: _lru}
		cacheMap.Set(key, val)
	}
	return val.(*_Cache[T])
}

func (_Int) Set(key string, value int)   { _int.Add(key, value) }
func (_Int) Get(key string) (int, bool)  { return _int.Get(key) }
func (_Int) Peek(key string) (int, bool) { return _int.Get(key) }
func (_Int) Remove(key string)           { _int.Remove(key) }
func (_Int) Exists(key string) bool      { return _int.Contains(key) }
func (_Int) Keys() []string              { return _int.Keys() }
func (_Int) Count() int                  { return _int.Len() }
func (_Int) Flush()                      { _int.Purge() }

func (_Int8) Set(key string, value int8)   { _int8.Add(key, value) }
func (_Int8) Get(key string) (int8, bool)  { return _int8.Get(key) }
func (_Int8) Peek(key string) (int8, bool) { return _int8.Get(key) }
func (_Int8) Remove(key string)            { _int8.Remove(key) }
func (_Int8) Exists(key string) bool       { return _int8.Contains(key) }
func (_Int8) Keys() []string               { return _int8.Keys() }
func (_Int8) Count() int                   { return _int8.Len() }
func (_Int8) Flush()                       { _int8.Purge() }

func (_Int16) Set(key string, value int16)   { _int16.Add(key, value) }
func (_Int16) Get(key string) (int16, bool)  { return _int16.Get(key) }
func (_Int16) Peek(key string) (int16, bool) { return _int16.Get(key) }
func (_Int16) Remove(key string)             { _int16.Remove(key) }
func (_Int16) Exists(key string) bool        { return _int16.Contains(key) }
func (_Int16) Keys() []string                { return _int16.Keys() }
func (_Int16) Count() int                    { return _int16.Len() }
func (_Int16) Flush()                        { _int16.Purge() }

func (_Int32) Set(key string, value int32)   { _int32.Add(key, value) }
func (_Int32) Get(key string) (int32, bool)  { return _int32.Get(key) }
func (_Int32) Peek(key string) (int32, bool) { return _int32.Get(key) }
func (_Int32) Remove(key string)             { _int32.Remove(key) }
func (_Int32) Exists(key string) bool        { return _int32.Contains(key) }
func (_Int32) Keys() []string                { return _int32.Keys() }
func (_Int32) Count() int                    { return _int32.Len() }
func (_Int32) Flush()                        { _int32.Purge() }

func (_Int64) Set(key string, value int64)   { _int64.Add(key, value) }
func (_Int64) Get(key string) (int64, bool)  { return _int64.Get(key) }
func (_Int64) Peek(key string) (int64, bool) { return _int64.Get(key) }
func (_Int64) Remove(key string)             { _int64.Remove(key) }
func (_Int64) Exists(key string) bool        { return _int64.Contains(key) }
func (_Int64) Keys() []string                { return _int64.Keys() }
func (_Int64) Count() int                    { return _int64.Len() }
func (_Int64) Flush()                        { _int64.Purge() }

func (_Uint) Set(key string, value uint)   { _uint.Add(key, value) }
func (_Uint) Get(key string) (uint, bool)  { return _uint.Get(key) }
func (_Uint) Peek(key string) (uint, bool) { return _uint.Get(key) }
func (_Uint) Remove(key string)            { _uint.Remove(key) }
func (_Uint) Exists(key string) bool       { return _uint.Contains(key) }
func (_Uint) Keys() []string               { return _uint.Keys() }
func (_Uint) Count() int                   { return _uint.Len() }
func (_Uint) Flush()                       { _uint.Purge() }

func (_Uint8) Set(key string, value uint8)   { _uint8.Add(key, value) }
func (_Uint8) Get(key string) (uint8, bool)  { return _uint8.Get(key) }
func (_Uint8) Peek(key string) (uint8, bool) { return _uint8.Get(key) }
func (_Uint8) Remove(key string)             { _uint8.Remove(key) }
func (_Uint8) Exists(key string) bool        { return _uint8.Contains(key) }
func (_Uint8) Keys() []string                { return _uint8.Keys() }
func (_Uint8) Count() int                    { return _uint8.Len() }
func (_Uint8) Flush()                        { _uint8.Purge() }

func (_Uint16) Set(key string, value uint16)   { _uint16.Add(key, value) }
func (_Uint16) Get(key string) (uint16, bool)  { return _uint16.Get(key) }
func (_Uint16) Peek(key string) (uint16, bool) { return _uint16.Get(key) }
func (_Uint16) Remove(key string)              { _uint16.Remove(key) }
func (_Uint16) Exists(key string) bool         { return _uint16.Contains(key) }
func (_Uint16) Keys() []string                 { return _uint16.Keys() }
func (_Uint16) Count() int                     { return _uint16.Len() }
func (_Uint16) Flush()                         { _uint16.Purge() }

func (_Uint32) Set(key string, value uint32)   { _uint32.Add(key, value) }
func (_Uint32) Get(key string) (uint32, bool)  { return _uint32.Get(key) }
func (_Uint32) Peek(key string) (uint32, bool) { return _uint32.Get(key) }
func (_Uint32) Remove(key string)              { _uint32.Remove(key) }
func (_Uint32) Exists(key string) bool         { return _uint32.Contains(key) }
func (_Uint32) Keys() []string                 { return _uint32.Keys() }
func (_Uint32) Count() int                     { return _uint32.Len() }
func (_Uint32) Flush()                         { _uint32.Purge() }

func (_Uint64) Set(key string, value uint64)   { _uint64.Add(key, value) }
func (_Uint64) Get(key string) (uint64, bool)  { return _uint64.Get(key) }
func (_Uint64) Peek(key string) (uint64, bool) { return _uint64.Get(key) }
func (_Uint64) Remove(key string)              { _uint64.Remove(key) }
func (_Uint64) Exists(key string) bool         { return _uint64.Contains(key) }
func (_Uint64) Keys() []string                 { return _uint64.Keys() }
func (_Uint64) Count() int                     { return _uint64.Len() }
func (_Uint64) Flush()                         { _uint64.Purge() }

func (_Float32) Set(key string, value float32)   { _float32.Add(key, value) }
func (_Float32) Get(key string) (float32, bool)  { return _float32.Get(key) }
func (_Float32) Peek(key string) (float32, bool) { return _float32.Get(key) }
func (_Float32) Remove(key string)               { _float32.Remove(key) }
func (_Float32) Exists(key string) bool          { return _float32.Contains(key) }
func (_Float32) Keys() []string                  { return _float32.Keys() }
func (_Float32) Count() int                      { return _float32.Len() }
func (_Float32) Flush()                          { _float32.Purge() }

func (_Float64) Set(key string, value float64)   { _float64.Add(key, value) }
func (_Float64) Get(key string) (float64, bool)  { return _float64.Get(key) }
func (_Float64) Peek(key string) (float64, bool) { return _float64.Get(key) }
func (_Float64) Remove(key string)               { _float64.Remove(key) }
func (_Float64) Exists(key string) bool          { return _float64.Contains(key) }
func (_Float64) Keys() []string                  { return _float64.Keys() }
func (_Float64) Count() int                      { return _float64.Len() }
func (_Float64) Flush()                          { _float64.Purge() }

func (_Bool) Set(key string, value bool)   { _bool.Add(key, value) }
func (_Bool) Get(key string) (bool, bool)  { return _bool.Get(key) }
func (_Bool) Peek(key string) (bool, bool) { return _bool.Get(key) }
func (_Bool) Remove(key string)            { _bool.Remove(key) }
func (_Bool) Exists(key string) bool       { return _bool.Contains(key) }
func (_Bool) Keys() []string               { return _bool.Keys() }
func (_Bool) Count() int                   { return _bool.Len() }
func (_Bool) Flush()                       { _bool.Purge() }

func (_Rune) Set(key string, value rune)   { _rune.Add(key, value) }
func (_Rune) Get(key string) (rune, bool)  { return _rune.Get(key) }
func (_Rune) Peek(key string) (rune, bool) { return _rune.Get(key) }
func (_Rune) Remove(key string)            { _rune.Remove(key) }
func (_Rune) Exists(key string) bool       { return _rune.Contains(key) }
func (_Rune) Keys() []string               { return _rune.Keys() }
func (_Rune) Count() int                   { return _rune.Len() }
func (_Rune) Flush()                       { _rune.Purge() }

func (_String) Set(key string, value string)   { _string.Add(key, value) }
func (_String) Get(key string) (string, bool)  { return _string.Get(key) }
func (_String) Peek(key string) (string, bool) { return _string.Get(key) }
func (_String) Remove(key string)              { _string.Remove(key) }
func (_String) Exists(key string) bool         { return _string.Contains(key) }
func (_String) Keys() []string                 { return _string.Keys() }
func (_String) Count() int                     { return _string.Len() }
func (_String) Flush()                         { _string.Purge() }

func (_Byte) Set(key string, value byte)   { _byte.Add(key, value) }
func (_Byte) Get(key string) (byte, bool)  { return _byte.Get(key) }
func (_Byte) Peek(key string) (byte, bool) { return _byte.Get(key) }
func (_Byte) Remove(key string)            { _byte.Remove(key) }
func (_Byte) Exists(key string) bool       { return _byte.Contains(key) }
func (_Byte) Keys() []string               { return _byte.Keys() }
func (_Byte) Count() int                   { return _byte.Len() }
func (_Byte) Flush()                       { _byte.Purge() }

func (_Bytes) Set(key string, value []byte)   { _bytes.Add(key, value) }
func (_Bytes) Get(key string) ([]byte, bool)  { return _bytes.Get(key) }
func (_Bytes) Peek(key string) ([]byte, bool) { return _bytes.Get(key) }
func (_Bytes) Remove(key string)              { _bytes.Remove(key) }
func (_Bytes) Exists(key string) bool         { return _bytes.Contains(key) }
func (_Bytes) Keys() []string                 { return _bytes.Keys() }
func (_Bytes) Count() int                     { return _bytes.Len() }
func (_Bytes) Flush()                         { _bytes.Purge() }

func (_Any) Set(key string, value any)   { _any.Add(key, value) }
func (_Any) Get(key string) (any, bool)  { return _any.Get(key) }
func (_Any) Peek(key string) (any, bool) { return _any.Get(key) }
func (_Any) Remove(key string)           { _any.Remove(key) }
func (_Any) Exists(key string) bool      { return _any.Contains(key) }
func (_Any) Keys() []string              { return _any.Keys() }
func (_Any) Count() int                  { return _any.Len() }
func (_Any) Flush()                      { _any.Purge() }

func (c *_Cache[M]) Set(key string, value M)   { c._lru.Add(key, value) }
func (c *_Cache[M]) Get(key string) (M, bool)  { return c._lru.Get(key) }
func (c *_Cache[M]) Peek(key string) (M, bool) { return c._lru.Get(key) }
func (c *_Cache[M]) Remove(key string)         { c._lru.Remove(key) }
func (c *_Cache[M]) Exists(key string) bool    { return c._lru.Contains(key) }
func (c *_Cache[M]) Keys() []string            { return c._lru.Keys() }
func (c *_Cache[M]) Count() int                { return c._lru.Len() }
func (c *_Cache[M]) Flush()                    { c._lru.Purge() }

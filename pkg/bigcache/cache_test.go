package bigcache

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheStruct(t *testing.T) {
	assert.NoError(t, Init())
	cache := New(&service{})
	defer func() {
		if err := cache.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	testCases := []struct {
		name string
		key  string
		val  any
	}{
		{"svc1", "svc1", &service{
			Name:        "svc1",
			Protocol:    "tcp",
			ForwardIP:   "127.0.0.1",
			RemotePort:  10080,
			ForwardPort: 80,
		}},
		{"svc2", "svc2", &service{}},
		{"svc3", "svc3", nil},
		{"svc4", "svc4", service{
			Name:        "svc4",
			Protocol:    "tcp",
			ForwardIP:   "127.0.0.1",
			RemotePort:  10080,
			ForwardPort: 80,
		}},
		{"svc5", "svc5", service{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(t, cache.Set(tc.key, tc.val))
			entry, err := cache.Get(tc.key)
			assert.NoError(t, err)
			assert.IsType(t, &service{}, entry)

			val := &service{}
			if tc.val != nil {
				if reflect.ValueOf(tc.val).Kind() == reflect.Pointer {
					val = tc.val.(*service)
				} else {
					v := tc.val.(service)
					val = &v
				}
			}
			assert.Equal(t, val, entry)
		})
	}
	cache.Close()
	cache.Close()
	cache.Close()

	cache2 := New(service{})
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(t, cache2.Set(tc.key, tc.val))
			entry, err := cache2.Get(tc.key)
			assert.NoError(t, err)
			assert.IsType(t, &service{}, entry)

			val := &service{}
			if tc.val != nil {
				if reflect.ValueOf(tc.val).Kind() == reflect.Pointer {
					val = tc.val.(*service)
				} else {
					v := tc.val.(service)
					val = &v
				}
			}
			assert.Equal(t, val, entry)
		})
	}
	cache2.Close()
	cache2.Close()
	cache2.Close()
}

func TestCacheString(t *testing.T) {
	assert.NoError(t, Init())

	s := "val3"
	testCasesString := []struct {
		name string
		key  string
		val  any
	}{
		{"str1", "str1", "val1"},
		{"str1", "str1", nil},
		{"str1", "str1", &s},
	}

	// cache := New(new(string))
	cache := New(&s)
	for _, tc := range testCasesString {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(t, cache.Set(tc.key, tc.val))
			entry, err := cache.Get(tc.key)
			assert.NoError(t, err)
			assert.IsType(t, entry, new(string))
			val := new(string)
			if tc.val != nil {
				if reflect.ValueOf(tc.val).Kind() == reflect.Pointer {
					val = (tc.val.(*string))
				} else {
					v := tc.val.(string)
					val = &v
				}
			}
			assert.Equal(t, val, entry)
		})
	}

	cache2 := New(s)
	for _, tc := range testCasesString {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(t, cache2.Set(tc.key, tc.val))
			entry, err := cache2.Get(tc.key)
			assert.NoError(t, err)
			assert.IsType(t, entry, new(string))
			val := new(string)
			if tc.val != nil {
				if reflect.ValueOf(tc.val).Kind() == reflect.Pointer {
					val = (tc.val.(*string))
				} else {
					v := tc.val.(string)
					val = &v
				}
			}
			assert.Equal(t, val, entry)
		})
	}
}

func TestCacheInt(t *testing.T) {
	assert.NoError(t, Init())

	i := 10
	testCases := []struct {
		name string
		key  string
		val  any
	}{
		{"int1", "int1", 1},
		{"int2", "int2", nil},
		{"int3", "int3", &i},
	}

	cache := New(&i)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(t, cache.Set(tc.key, tc.val))
			entry, err := cache.Get(tc.key)
			assert.NoError(t, err)
			assert.IsType(t, entry, &i)
			val := new(int)
			if tc.val != nil {
				if reflect.ValueOf(tc.val).Kind() == reflect.Pointer {
					val = tc.val.(*int)
				} else {
					v := tc.val.(int)
					val = &v
				}
			}
			assert.Equal(t, val, entry)
		})
	}
}

type service struct {
	Name        string
	Protocol    string
	ForwardIP   string
	RemotePort  int
	ForwardPort int
}

package bigcache

import (
	"encoding/json"
	"reflect"
)

type SessionStore Cache

// NewSessionStore creates a session store that expires automatically.
func NewSessionStore(model any) *SessionStore {
	if reflect.ValueOf(model).Kind() == reflect.Pointer {
		return &SessionStore{
			cache: sessionCache,
			model: reflect.TypeOf(model).Elem(),
		}
	}
	return &SessionStore{
		cache: sessionCache,
		model: reflect.TypeOf(model),
	}
}
func (c *SessionStore) Len() int                { return c.cache.Len() }
func (c *SessionStore) Cap() int                { return c.cache.Capacity() }
func (c *SessionStore) Reset() error            { return c.cache.Reset() }
func (c *SessionStore) Delete(key string) error { return c.cache.Delete(key) }

func (c *SessionStore) Get(key string) (any, error) {
	entry, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	i := (reflect.New(c.model)).Interface()
	if err := json.Unmarshal(entry, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (c *SessionStore) Set(key string, entry any) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return c.cache.Set(key, data)
}

func (c *SessionStore) Close() (err error) {
	once.Do(func() {
		err = sessionCache.Close()
	})
	return
}

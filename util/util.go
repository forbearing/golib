package util

import (
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
)

func UUID(name ...string) string {
	var _name string
	if len(name) > 0 {
		_name = name[0]
	}
	if len(_name) == 0 {
		_name = uuid.NewV4().String()
	}
	return uuid.NewV5(uuid.NewV4(), _name).String()
}

func IndexedUUID() string {
	return ksuid.New().String()
}

// Pointer will return a pointer to T with given value.
func Pointer[T comparable](t T) *T {
	if reflect.DeepEqual(t, nil) {
		return new(T)
	}
	return &t
}

// Depointer will return a T with given value.
func Depointer[T comparable](t *T) T {
	if t == nil {
		return *new(T)
	}
	return *t
}

package util

import (
	"bytes"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
)

// UUID is a generic uuid generator.
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

// LightUUID used in logger request id.
// It will use uuid.NewV4() instead of uuid.NewV5() which have less cpu cost.
func LightUUID() string { return uuid.NewV4().String() }

// IndexedUUID generate indexable uuid.
func IndexedUUID() string { return ksuid.New().String() }

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

// CharSpliter is the custom split function for bufio.Scanner.
func CharSpliter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}
	if atEOF {
		return len(data), data, nil
	}
	if data[0] == '|' {
		return 1, data[:1], nil
	}
	return 0, nil, nil
}

// SplitByDoublePipe is the custom split function for bufio.Scanner.
func SplitByDoublePipe(data []byte, atEOF bool) (advance int, token []byte, err error) {
	delimiter := []byte("||")

	// Search for the delimiter in the input data
	if i := bytes.Index(data, delimiter); i >= 0 {
		return i + len(delimiter), data[:i], nil
	}

	// If the delimiter is not found, and it's at the end of the input data, return it
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}

	// If no delimiter is found, return no data and wait for more input
	return 0, nil, nil
}

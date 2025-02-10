package rbtree

import (
	"strings"
	"sync"
)

// Option is a functional option type for configuring a red-black tree.
type Option[K comparable, V any] func(t *Tree[K, V]) error

// WithSafe creates a option that make the red-black tree safe for concurrent use.
func WithSafe[K comparable, V any]() Option[K, V] {
	return func(t *Tree[K, V]) error {
		t.safe = true
		t.mu = &sync.RWMutex{}
		return nil
	}
}

// WithNodeFormat creates a option that sets the node format when call tree.String().
// Default node format is fmt.Sprintf("%v", n.Key).
// The format must contains two placeholders, the first is the key and the second is the value.
// For example: "%d:%s"
func WithNodeFormat[K comparable, V any](format string) Option[K, V] {
	return func(t *Tree[K, V]) error {
		format = strings.TrimFunc(format, func(r rune) bool {
			return r == '\n'
		})
		t.nodeFormat = format
		return nil
	}
}

// WithColorfulString creates a option that make the red-black tree colorful output.
func WithColorfulString[K comparable, V any]() Option[K, V] {
	return func(t *Tree[K, V]) error {
		t.color = true
		return nil
	}
}

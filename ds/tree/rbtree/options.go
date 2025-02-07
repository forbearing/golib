package rbtree

import "sync"

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

// WithColorfulString creates a option that make the red-black tree colorful output.
func WithColorfulString[K comparable, V any]() Option[K, V] {
	return func(t *Tree[K, V]) error {
		t.color = true
		return nil
	}
}

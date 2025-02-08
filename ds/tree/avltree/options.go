package avltree

import "sync"

// Option is a functional option type for configuring a AVL tree.
type Option[K comparable, V any] func(*Tree[K, V]) error

// WithSafe creates a option that make the AVL tree safe for concurrent use.
func WithSafe[K comparable, V any]() Option[K, V] {
	return func(t *Tree[K, V]) error {
		t.safe = true
		t.mu = &sync.RWMutex{}
		return nil
	}
}

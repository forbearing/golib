package trie

import (
	"sync"

	"github.com/forbearing/golib/ds/types"
)

// Option is a function option type for configuring a trie.
type Option[K comparable, V any] func(*Trie[K, V]) error

// WithSafe returns a option that makes the trie safe for concurrent use.
func WithSafe[K comparable, V any]() Option[K, V] {
	return func(t *Trie[K, V]) error {
		t.safe = true
		t.mu = &sync.RWMutex{}
		return nil
	}
}

func WithNodeFormatter[K comparable, V any](fn func(n *Node[K, V]) string) Option[K, V] {
	return func(t *Trie[K, V]) error {
		if fn == nil {
			return types.ErrFuncNil
		}
		t.nodeFormatter = fn
		return nil
	}
}

func WithKeyFormatter[K comparable, V any](fn func(K, *Node[K, V]) string) Option[K, V] {
	return func(t *Trie[K, V]) error {
		if fn == nil {
			return types.ErrFuncNil
		}
		t.keyFormatter = fn
		return nil
	}
}

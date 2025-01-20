package linkedlist

import "sync"

type Option[V any] func(*List[V]) error

// WithSafe creates a Option that make the doublely-linked list safe for concurrent use.
func WithSafe[V any]() Option[V] {
	return func(m *List[V]) error {
		m.mu = new(sync.RWMutex)
		return nil
	}
}

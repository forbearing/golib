package arraylist

import "sync"

type Option[V any] func(*List[V]) error

// WithSafe creates a Option that make the array-backed list safe for concurrent use.
func WithSafe[V any]() Option[V] {
	return func(l *List[V]) error {
		l.mu = new(sync.RWMutex)
		return nil
	}
}

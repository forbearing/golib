package set

import "sync"

// Option is a functional option type for configuring a Set.
type Option[T comparable] func(*Set[T]) error

// WithSafe creates a Option that make the Set safe for concurrent use.
func WithSafe[T comparable]() Option[T] {
	return func(s *Set[T]) error {
		s.mu = new(sync.RWMutex)
		s.safe = true
		return nil
	}
}

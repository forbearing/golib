package arraystack

// Option is funtional option type for configuring a Stack.
type Option[E any] func(*Stack[E]) error

// WithSafe creates a option that make the array-backed stack safe for concurrent use.
func WithSafe[E any]() Option[E] {
	return func(s *Stack[E]) error {
		s.safe = true
		return nil
	}
}

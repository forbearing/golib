package arraystack

import (
	"fmt"
	"slices"
	"strings"

	"github.com/forbearing/golib/ds/list/arraylist"
)

// Stack represents a stack based on array-backed list.
// The stack provides typical LIFO(last-in, first-out) behavior.
type Stack[E any] struct {
	list *arraylist.List[E]
	cmp  func(E, E) int
	safe bool
}

// New creates and initializes a empty stack.
// The "cmp" function is used to compare elements for equality.
// Options can be provided to customize the stack's properties (e.g., thread safety).
func New[E any](cmp func(E, E) int, ops ...Option[E]) (s *Stack[E], err error) {
	s = new(Stack[E])
	s.cmp = cmp
	for _, op := range ops {
		if op == nil {
			continue
		}
		if err = op(s); err != nil {
			return nil, err
		}
	}
	if s.safe {
		s.list, err = arraylist.New(cmp, arraylist.WithSafe[E]())
	} else {
		s.list, err = arraylist.New(cmp)
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// NewFromSlice creates and initializes a stack from the provided slice.
// The "cmp" function is used to compare elements for equality.
// Options can be provided to customize the stack's properties (e.g., thread safety).
func NewFromSlice[E any](cmp func(E, E) int, slice []E, ops ...Option[E]) (*Stack[E], error) {
	s, err := New(cmp, ops...)
	if err != nil {
		return nil, err
	}
	if len(slice) > 0 {
		s.list.Append(slice...)
	}
	return s, nil
}

// NewFromMapKeys creates and initializes a stack from the provided map keys.
// The "cmp" function is used to compare elements for equality.
// Options can be provided to customize the stack's properties (e.g., thread safety).
// Returns an empty stack if the provided map is nil or empty.
func NewFromMapKeys[K comparable, V any](cmp func(K, K) int, m map[K]V, ops ...Option[K]) (*Stack[K], error) {
	s, err := New(cmp, ops...)
	if err != nil {
		return nil, err
	}
	if len(m) > 0 {
		for k := range m {
			s.list.Append(k)
		}
	}
	return s, nil
}

// NewFromMapValues creates a stack from the provided map values.
// Options can be provided to customize the stack's properties (e.g., thread safety).
// Returns an empty stack if the provided map is nil or empty.
func NewFromMapValues[K comparable, V any](cmp func(V, V) int, m map[K]V, ops ...Option[V]) (*Stack[V], error) {
	s, err := New(cmp, ops...)
	if err != nil {
		return nil, err
	}
	if len(m) > 0 {
		for _, v := range m {
			s.list.Append(v)
		}
	}
	return s, nil
}

// Push adds an element to the top of the stack.
// The stack's size increases by one.
func (s *Stack[E]) Push(e E) {
	s.list.Append(e)
}

// Pop removes and returns the top element of the stack.
// Returns the zero value of E and false if the stack is empty.
func (s *Stack[E]) Pop() (E, bool) {
	var e E
	if s.Len() == 0 {
		return e, false
	}
	e = s.list.RemoveAt(s.Len() - 1)
	return e, true
}

// Peek returns the top element of the stack without removing it.
// Returns the zero value of E and false if the stack is empty.
func (s *Stack[E]) Peek() (E, bool) {
	var e E
	if s.Len() == 0 {
		return e, false
	}
	return s.list.Get(s.Len() - 1)
}

// IsEmpty reports whether the stack has no elements.
func (s *Stack[E]) IsEmpty() bool {
	return s.Len() == 0
}

// Len returns the number of elements currently in the stack.
func (s *Stack[E]) Len() int {
	return s.list.Len()
}

// Values returns all elements in the stack in LIFO(last-in, first-out) order.
// If the stack is empty, it returns an empty slice (not nil).
func (s *Stack[E]) Values() []E {
	el := s.list.Values()
	slices.Reverse(el)
	return el
}

// Clear removes all elements from the stack.
func (s *Stack[E]) Clear() {
	s.list.Clear()
}

// Clone returns a deep copy of the stack.
func (s *Stack[E]) Clone() *Stack[E] {
	clone, _ := NewFromSlice(s.cmp, s.list.Values(), s.options()...)
	return clone
}

func (s *Stack[E]) options() []Option[E] {
	ops := make([]Option[E], 0)
	if s.safe {
		ops = append(ops, WithSafe[E]())
	}
	return ops
}

// String returns a string representation of the stack.
func (s *Stack[E]) String() string {
	el := s.list.Values()
	slices.Reverse(el)
	items := make([]string, 0, s.list.Len())
	for _, e := range el {
		items = append(items, fmt.Sprintf("%v", e))
	}
	return fmt.Sprintf("stack:{%s}", strings.Join(items, ", "))
}

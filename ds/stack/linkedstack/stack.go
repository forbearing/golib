package linkedstack

import (
	"fmt"
	"slices"
	"strings"

	"github.com/forbearing/golib/ds/list/linkedlist"
)

// Stack represents a stack based on linkedlist..
// The stack provides typical LIFO(last-in, first-out) behavior.
type Stack[E any] struct {
	list *linkedlist.List[E]
	safe bool
}

// New creates and initializes a empty stack.
// Options can be provided to customize the stack's properties (e.g., thread safety).
func New[E any](ops ...Option[E]) (s *Stack[E], err error) {
	s = new(Stack[E])
	for _, op := range ops {
		if op == nil {
			continue
		}
		if err = op(s); err != nil {
			return nil, err
		}
	}
	if s.safe {
		s.list, err = linkedlist.New(linkedlist.WithSafe[E]())
	} else {
		s.list, err = linkedlist.New[E]()
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// NewFromSlice creates and initializes a stack from the provided slice.
// Options can be provided to customize the stack's properties (e.g., thread safety).
func NewFromSlice[E any](slice []E, ops ...Option[E]) (*Stack[E], error) {
	s, err := New(ops...)
	if err != nil {
		return nil, err
	}
	if len(slice) > 0 {
		for _, e := range slice {
			s.list.PushBack(e)
		}
	}
	return s, nil
}

// NewFromMapKeys creates and initializes a stack from the provided map keys.
// Options can be provided to customize the stack's properties (e.g., thread safety).
// Returns an empty stack if the provided map is nil or empty.
func NewFromMapKeys[K comparable, V any](m map[K]V, ops ...Option[K]) (*Stack[K], error) {
	s, err := New(ops...)
	if err != nil {
		return nil, err
	}
	if len(m) > 0 {
		for k := range m {
			s.list.PushBack(k)
		}
	}
	return s, nil
}

// NewFromMapValues creates a stack from the provided map values.
// Options can be provided to customize the stack's properties (e.g., thread safety).
// Returns an empty stack if the provided map is nil or empty.
func NewFromMapValues[K comparable, V any](cmp func(V, V) int, m map[K]V, ops ...Option[V]) (*Stack[V], error) {
	s, err := New(ops...)
	if err != nil {
		return nil, err
	}
	if len(m) > 0 {
		for _, v := range m {
			s.list.PushBack(v)
		}
	}
	return s, nil
}

// Push adds an element to the top of the stack.
// The stack's size increases by one.
func (s *Stack[E]) Push(e E) {
	s.list.PushBack(e)
}

// Pop removes and returns the top element of the stack.
// Returns the zero value of E and false if the stack is empty.
func (s *Stack[E]) Pop() (E, bool) {
	var e E
	if s.Len() == 0 {
		return e, false
	}
	return s.list.PopBack(), true
}

// Peek returns the top element of the stack without removing it.
// Returns the zero value of E and false if the stack is empty.
func (s *Stack[E]) Peek() (E, bool) {
	var e E
	if s.Len() == 0 {
		return e, false
	}
	return s.list.Tail.Value, true
}

// IsEmpty reports whether the stack has no elements.
func (s *Stack[E]) IsEmpty() bool {
	return s.list.Len() == 0
}

// Len returns the number of elements currently in the stack.
func (s *Stack[E]) Len() int {
	return s.list.Len()
}

// Values returns all elements in the stack in LIFO(last-in, first-out) order.
// If the stack is empty, it returns an empty slice (not nil).
func (s *Stack[E]) Values() []E {
	el := s.list.Slice()
	slices.Reverse(el)
	return el
}

// Clear removes all elements from the stack.
func (s *Stack[E]) Clear() {
	s.list.Clear()
}

// Clone returns a deep copy of the stack.
func (s *Stack[E]) Clone() *Stack[E] {
	clone, _ := NewFromSlice(s.list.Slice())
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
	el := make([]E, 0, s.list.Len())
	s.list.Range(func(e E) bool {
		el = append(el, e)
		return true
	})
	slices.Reverse(el)
	items := make([]string, 0, s.list.Len())
	for _, e := range el {
		items = append(items, fmt.Sprintf("%v", e))
	}
	return fmt.Sprintf("stack:{%s}", strings.Join(items, ", "))
}

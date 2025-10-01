package mapset

import (
	"fmt"
	"slices"
	"strings"

	"github.com/forbearing/gst/ds/types"
)

type Set[E comparable] struct {
	set    map[E]struct{}
	mu     types.Locker
	safe   bool
	cmp    func(E, E) int
	sorted bool
}

// New creates a new set without pre-allocates space.
// Options can be provided to customize the set's properties (e.g., thread safety).
func New[E comparable](ops ...Option[E]) (*Set[E], error) {
	return NewWithSize(0, ops...)
}

// NewWithSize creates a new set and pre-allocates space for the given size.
// Options can be provided to customize the set's properties (e.g., thread safety).
func NewWithSize[T comparable](size int, ops ...Option[T]) (*Set[T], error) {
	s := &Set[T]{
		set: make(map[T]struct{}, size),
		mu:  types.FakeLocker{},
	}
	for _, op := range ops {
		if op == nil {
			continue
		}
		if err := op(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// NewFromSlice creates a new set from the provided slice.
// If the provided slice is nil or empty, creates an empty set.
// Options can be provided to customize the set's properties (e.g., thread safety).
func NewFromSlice[E comparable](slice []E, ops ...Option[E]) (*Set[E], error) {
	if len(slice) == 0 {
		return New(ops...)
	}
	s, err := NewWithSize(len(slice), ops...)
	if err != nil {
		return nil, err
	}
	for _, e := range slice {
		s.set[e] = struct{}{}
	}
	return s, nil
}

// NewFromMapKeys creates a new set from the provided map of keys.
// If map "m" is nil or empty, creates an empty set.
// Options can be provided to customize the set's properties (e.g., thread safety).
func NewFromMapKeys[K comparable, V any](m map[K]V, ops ...Option[K]) (*Set[K], error) {
	if len(m) == 0 {
		return New(ops...)
	}
	s, err := NewWithSize(len(m), ops...)
	if err != nil {
		return nil, err
	}
	for k := range m {
		s.set[k] = struct{}{}
	}
	return s, nil
}

// NewFromMapValues creates a new set from the provided map of values.
// If map "m" is nil or empty, creates an empty set.
// Options can be provided to customize the set's properties (e.g., thread safety).
func NewFromMapValues[K comparable, V comparable](m map[K]V, ops ...Option[V]) (*Set[V], error) {
	if len(m) == 0 {
		return New(ops...)
	}
	s, err := NewWithSize(len(m), ops...)
	if err != nil {
		return nil, err
	}
	for _, v := range m {
		s.set[v] = struct{}{}
	}
	return s, nil
}

// Add one or more elements into the set.
// Returns the number of elements added.
func (s *Set[E]) Add(el ...E) int {
	if s.safe {
		s.mu.Lock()
		defer s.mu.Unlock()
	}

	prevLen := len(s.set)
	for _, e := range el {
		s.set[e] = struct{}{}
	}
	return len(s.set) - prevLen
}

// Pop removes and returns a single, arbitrary element from the set.
// The order of removal is non-deterministic.
// If the set is empty, it returns zero value of element type and false.
func (s *Set[E]) Pop() (e E, ok bool) {
	if s.safe {
		s.mu.Lock()
		defer s.mu.Unlock()
	}

	for v := range s.set {
		delete(s.set, v)
		return v, true
	}
	return e, false
}

// Remove one or more elements from the set.
func (s *Set[E]) Remove(el ...E) {
	if s.safe {
		s.mu.Lock()
		defer s.mu.Unlock()
	}

	for _, e := range el {
		delete(s.set, e)
	}
}

// Clear removes all elements from the set.
func (s *Set[E]) Clear() {
	if s.safe {
		s.mu.Lock()
		defer s.mu.Unlock()
	}

	for e := range s.set {
		delete(s.set, e)
	}
}

// Len returns the number of elements in the set.
func (s *Set[E]) Len() int {
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	return len(s.set)
}

// Clone creates and returns a deep copy of the set.
//
// The property of the cloned set is the same as the original set.
// - If the original set is concurrent safe, the cloned set is concurrent safe.
func (s *Set[E]) Clone() *Set[E] {
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	return s.clone()
}

func (s *Set[E]) clone() *Set[E] {
	var cloned *Set[E]
	cloned, _ = NewFromMapKeys(s.set, s.options()...)
	return cloned
}

func (s *Set[E]) options() []Option[E] {
	ops := []Option[E]{}
	if s.safe {
		ops = append(ops, WithSafe[E]())
	}
	if s.sorted {
		ops = append(ops, WithSorted(s.cmp))
	}
	return ops
}

// Contains reports whether the set contains all the given elements.
// It always returns true if the provided slice is nil or empty.
func (s *Set[E]) Contains(el ...E) bool {
	if len(el) == 0 {
		return true
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	var ok bool
	for _, e := range el {
		if _, ok = s.set[e]; !ok {
			return false
		}
	}
	return true
}

// ContainsOne reports whether the set contains the given element.
func (s *Set[E]) ContainsOne(v E) bool {
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	_, ok := s.set[v]
	return ok
}

// ContainsAny reports whether the set contains any of the given element.
// It returns true if the provided slice is nil or empty.
func (s *Set[E]) ContainsAny(el ...E) bool {
	if len(el) == 0 {
		return true
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	var ok bool
	for _, e := range el {
		if _, ok = s.set[e]; ok {
			return true
		}
	}
	return false
}

// ContainsAnyElement reports whether the set contains any element of the given set.
// If the given set is nil or empty, it returns false.
func (s *Set[E]) ContainsAnyElement(other *Set[E]) bool {
	if other == nil {
		return false
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.safe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	if len(other.set) == 0 {
		return false
	}

	var ok bool
	if len(s.set) < len(other.set) {
		for e := range s.set {
			if _, ok = other.set[e]; ok {
				return true
			}
		}
	} else {
		for e := range other.set {
			if _, ok = s.set[e]; ok {
				return true
			}
		}
	}
	return false
}

// range calls fn for each element in the set.
// If fn returns false, "Range" stops the iteration.
// If fn is nil, "Range" does nothing.
func (s *Set[E]) Range(fn func(e E) bool) {
	if fn == nil {
		return
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	if s.sorted {
		el := s.sortedSlice(s.cmp)
		for _, e := range el {
			if !fn(e) {
				return
			}
		}
	} else {
		for e := range s.set {
			if !fn(e) {
				return
			}
		}
	}
}

// Equal reports whether two sets have the same elements.
func (s *Set[E]) Equal(other *Set[E]) bool {
	if other == nil {
		return false
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.safe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	if len(s.set) != len(other.set) {
		return false
	}

	var ok bool
	for e := range s.set {
		if _, ok = other.set[e]; !ok {
			return false
		}
	}
	return true
}

// IsEmpty reports whether the set is empty.
func (s *Set[E]) IsEmpty() bool {
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	return len(s.set) == 0
}

// Iter returns a channel of elements that caller can range over.
func (s *Set[E]) Iter() <-chan E {
	ch := make(chan E)
	go func() {
		if s.safe {
			s.mu.RLock()
			defer s.mu.RUnlock()
		}
		if s.sorted {
			el := s.sortedSlice(s.cmp)
			for _, e := range el {
				ch <- e
			}
		} else {
			for e := range s.set {
				ch <- e
			}
		}
		close(ch)
	}()
	return ch
}

// IsSubset checks if the current set is a subset of the given set.
// A subset means every element of the current set is also in the given set.
// If the given set is nil, the function always returns false.
func (s *Set[E]) IsSubset(other *Set[E]) bool {
	if other == nil {
		return false
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.safe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}

	return s.isSubset(other)
}

func (s *Set[E]) isSubset(other *Set[E]) bool {
	if len(s.set) > len(other.set) {
		return false
	}
	var ok bool
	for e := range s.set {
		if _, ok = other.set[e]; !ok {
			return false
		}
	}
	return true
}

// IsProperSubset checks if the current set is a proper subset of the given set.
// A proper subset means every element of the current set is in the given set,
// and the given set contains more elements than the current set.
func (s *Set[E]) IsProperSubset(other *Set[E]) bool {
	if other == nil {
		return false
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.safe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}

	return len(s.set) < len(other.set) && s.isSubset(other)
}

// IsSuperset checks if the current set is a superset of the given set.
// A superset means the current set contains every element of the given set.
// If the given set is nil or empty, the function always returns true.
func (s *Set[E]) IsSuperset(other *Set[E]) bool {
	if other == nil {
		return true
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.safe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}

	return s.isSuperset(other)
}

func (s *Set[E]) isSuperset(other *Set[E]) bool {
	if len(other.set) == 0 {
		return true
	}
	var ok bool
	for e := range other.set {
		if _, ok = s.set[e]; !ok {
			return false
		}
	}
	return true
}

// IsProperSuperset checks if the current set is a proper superset of given set.
// A proper superset means all elements of given set are present int the current set.
// and the current set has additional element not present in the given set.
func (s *Set[E]) IsProperSuperset(other *Set[E]) bool {
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other == nil && len(s.set) > 0 {
		return true
	}
	if other.safe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	return len(s.set) > len(other.set) && s.isSuperset(other)
}

// Difference computes the difference between the current set and the given set.
// The resulting set contains element that are present in the current set
// but not in the given set.
//
// The returned set inherits the properties of the current set.
// For example: if the current set is concurrent-safe, the returned set is also
// be concurrent-safe.
func (s *Set[E]) Difference(other *Set[E]) *Set[E] {
	if other == nil {
		return s.clone()
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.safe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	if len(other.set) == 0 || len(s.set) == 0 {
		return s.clone()
	}

	diff, _ := New(s.options()...)
	for e := range s.set {
		if _, ok := other.set[e]; !ok {
			diff.set[e] = struct{}{}
		}
	}
	return diff
}

// SymmetricDifference computes the symmetric difference between the current set
// and the given set.
// The symmetric difference includes elements present in either set but not in both.
//
// The returned set inherits the properties of the current set.
// For example, if the current set is concurrent-safe, the returned set is also
// be concurrent-safe
func (s *Set[E]) SymmetricDifference(other *Set[E]) *Set[E] {
	if other == nil {
		return s.Clone()
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if len(other.set) == 0 {
		return s.Clone()
	}

	diff, _ := New(s.options()...)
	for e := range s.set {
		if _, ok := other.set[e]; !ok {
			diff.set[e] = struct{}{}
		}
	}
	for e := range other.set {
		if _, ok := s.set[e]; !ok {
			diff.set[e] = struct{}{}
		}
	}
	return diff
}

// Union returns computes union of the current set and the given set.
// The resulting is contains all the elements that are present in
// either the current set or the given set.
//
// If the given set is nil or empty, returns the deep clone of the current set.
//
// The returned set inherits the properties of the current set.
// For example, if the current set is concurrent-safe, the returned set is also
// be concurrent-safe
func (s *Set[E]) Union(other *Set[E]) *Set[E] {
	if other == nil {
		return s.Clone()
	}
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.safe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	if len(other.set) == 0 {
		return s.Clone()
	}

	union, _ := New(s.options()...)
	for e := range s.set {
		union.set[e] = struct{}{}
	}
	for e := range other.set {
		union.set[e] = struct{}{}
	}
	return union
}

// Intersect computes the intersection of the current set and the given set.
// The resulting set contains elements that are present in both the current set and the given set.
//
// If the given set is nil or empty, returns an empty set.
// The returned set inherits the properties of the current set.
// For example, if the current set is concurrent-safe, the returned set is also
func (s *Set[E]) Intersect(other *Set[E]) *Set[E] {
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if other.safe {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	inter, _ := New(s.options()...)
	if other == nil {
		return inter
	}
	if len(other.set) == 0 || len(s.set) == 0 {
		return inter
	}

	if len(s.set) < len(other.set) {
		for e := range s.set {
			if _, ok := other.set[e]; ok {
				inter.set[e] = struct{}{}
			}
		}
	} else {
		for e := range other.set {
			if _, ok := s.set[e]; ok {
				inter.set[e] = struct{}{}
			}
		}
	}

	return inter
}

// String returns a string representation of the set.
func (s *Set[E]) String() string {
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	el := make([]string, 0, len(s.set))
	if s.sorted {
		elements := s.sortedSlice(s.cmp)
		for _, e := range elements {
			el = append(el, fmt.Sprintf("%v", e))
		}
	} else {
		for e := range s.set {
			el = append(el, fmt.Sprintf("%v", e))
		}
	}
	return fmt.Sprintf("Set{%s}", strings.Join(el, ", "))
}

// Slice returns a slice of the elements in the set.
// The order of the elements is non-deterministic.
func (s *Set[E]) Slice() []E {
	if s.safe {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	if s.sorted {
		return s.sortedSlice(s.cmp)
	}
	return s.unsortedSlice()
}

func (s *Set[E]) sortedSlice(cmp func(E, E) int) []E {
	el := make([]E, 0, len(s.set))
	for e := range s.set {
		el = append(el, e)
	}
	slices.SortFunc(el, cmp)
	return el
}

func (s *Set[E]) unsortedSlice() []E {
	el := make([]E, 0, len(s.set))
	for e := range s.set {
		el = append(el, e)
	}
	return el
}

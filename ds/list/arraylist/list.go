// Package arraylist provides a generic implementation of a resizable array-backed list.

package arraylist

import (
	"errors"
	"slices"

	"github.com/forbearing/golib/ds/types"
	"github.com/forbearing/golib/ds/util"
)

const (
	growthFactor = float32(2.0)
	shrinkFactor = float32(0.25)

	minCap = 16
)

var ErrEqualNil = errors.New("equal function is nil")

// List represents a resizable array-backend list.
// Call New or NewFromSlice default creates a list witout concurrent safety.
// Call New or NewFromSlice with `WithSafe` option to make the List safe for concurrent use.
type List[T any] struct {
	elements []T
	equal    func(T, T) bool
	mu       types.Locker

	safe bool
}

// New creates and returns a new array-backed list.
// The provided equal function is used to compare values for equality.
// Optional options can be passed to modify the list's behavior, such as enabling concurrent safety.
func New[T any](equal func(T, T) bool, ops ...Option[T]) (*List[T], error) {
	if equal == nil {
		return nil, ErrEqualNil
	}
	l := &List[T]{
		elements: make([]T, 0, minCap), // NOTE: zero capacity will cause growBy blocked.
		mu:       types.FakeLocker{},
		equal:    equal,
	}
	for _, op := range ops {
		if op == nil {
			continue
		}
		if err := op(l); err != nil {
			return nil, err
		}
	}
	return l, nil
}

// NewFromSlice creates a new array-backed list from the given slice.
// The provided equal function is used to compare values for equality.
// Optional options can be passed to modify the list's behavior, such as enabling concurrent safety.
func NewFromSlice[T any](equal func(T, T) bool, values []T, ops ...Option[T]) (*List[T], error) {
	l, err := New(equal, ops...)
	if err != nil {
		return nil, err
	}
	l.growBy(len(values))
	copy(l.elements, values)
	return l, nil
}

// Get returns the value at the given index.
func (l *List[T]) Get(index int) (T, bool) {
	// Checking "l.safe" before acquiring the lock is more efficient based on benchmark result.
	if l.safe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	if !l.withinRange(index, false) {
		var v T
		return v, false
	}
	return l.elements[index], true
}

// Append appends specified values to the end of the list.
func (l *List[T]) Append(values ...T) {
	if len(values) == 0 {
		return
	}
	// Checking "l.safe" before acquiring the lock is more efficient based on benchmark result.
	if l.safe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	l.append(values...)
}

func (l *List[T]) append(values ...T) {
	oldLen := len(l.elements)
	l.growBy(len(values))
	for i := range values {
		l.elements[oldLen+i] = values[i]
	}
}

// Insert inserts values at the given index.
// If the index is the length of the list, the values will be appended.
// If the index out of range, this function is no-op.
func (l *List[T]) Insert(index int, values ...T) {
	if len(values) == 0 {
		return
	}
	// Checking "l.safe" before acquiring the lock is more efficient based on benchmark result.
	if l.safe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	if !l.withinRange(index, true) {
		return
	}
	if index == len(l.elements) {
		l.append(values...)
		return
	}

	addLen := len(values)
	oldLen := len(l.elements)
	l.growBy(addLen)
	// move elements after index + length of values.
	copy(l.elements[index+addLen:], l.elements[index:oldLen])
	// copy values after index
	copy(l.elements[index:index+addLen], values)
}

// Set sets the value at the given index.
// If the index is the length of the list, the value will be appended.
// If the index out of range, this function is no-op.
func (l *List[T]) Set(index int, value T) {
	// Checking "l.safe" before acquiring the lock is more efficient based on benchmark result.
	if l.safe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	if !l.withinRange(index, true) {
		return
	}
	if index == len(l.elements) {
		l.append(value)
		return
	}
	l.elements[index] = value
}

// Remove removes all the value from the list.
func (l *List[T]) Remove(v T) {
	// Checking "l.safe" before acquiring the lock is more efficient based on benchmark result.
	if l.safe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	i := 0
	for i < len(l.elements) {
		if l.equal(v, l.elements[i]) {
			l.removeAt(i)
		} else {
			i++
		}
	}
}

// RemoveAt removes the value at the given index.
// If the index out of range, this function is no-op and returns zero value of T.
func (l *List[T]) RemoveAt(index int) T {
	// Checking "l.safe" before acquiring the lock is more efficient based on benchmark result.
	if l.safe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	return l.removeAt(index)
}

func (l *List[T]) removeAt(index int) T {
	var v T
	if !l.withinRange(index, false) {
		return v
	}
	v = l.elements[index]
	// equivalent to
	// l.elements = append(l.elements[:index], l.elements[index+1:]...)
	l.elements = slices.Delete(l.elements, index, index+1)
	l.shrink()
	return v
}

// Clear removes all elements from the list.
func (l *List[T]) Clear() {
	if l.safe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	// l.elements = l.elements[:0]
	// l.elements = nil
	l.elements = make([]T, 0)
}

// Contains reports whether the list contains all the given values.
// Returns true if all values are present in the list, false otherwise.
func (l *List[T]) Contains(values ...T) bool {
	if len(values) == 0 {
		return false
	}

	for _, v := range values {
		if !l.contains(v) {
			return false
		}
	}
	return true
}

func (l *List[T]) contains(value T) bool {
	// Skipping the "l.safe" check is more efficient based on benchmark result.
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, v := range l.elements {
		if l.equal(v, value) {
			return true
		}
	}
	return false
}

// Values returns a slice containing all values in the list.
// The returned slice is a copy of the internal slice,
// and modifications to it will not affect the list.
func (l *List[T]) Values() []T {
	// Skipping the "l.safe" check is more efficient based on benchmark result.
	if l.safe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}
	return slices.Clone(l.elements)
}

// IndexOf returns the index of the first occurrence of value in the list.
func (l *List[T]) IndexOf(value T) int {
	// Skipping the "l.safe" check is more efficient based on benchmark result.
	l.mu.RLock()
	defer l.mu.RUnlock()

	for i, v := range l.elements {
		if l.equal(v, value) {
			return i
		}
	}
	return -1
}

// IsEmpty reports whether the list has no elements.
func (l *List[T]) IsEmpty() bool {
	// Checking "l.safe" before acquiring the lock is more efficient based on benchmark result.
	if l.safe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}
	return len(l.elements) == 0
}

// Len returns the number of elements in the list.
func (l *List[T]) Len() int {
	// Checking "l.safe" before acquiring the lock is more efficient based on benchmark result.
	if l.safe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	return len(l.elements)
}

// Sort sorts the list using the given comparator
// if cmp is nil, the function is no-op.
// cmp should return:
// - A negative value if first argument is less than second.
// - Zero if the arguments are equal.
// - A positive value if first argument is greater than second.
func (l *List[T]) Sort(cmp util.Comparator[T]) {
	if cmp == nil {
		return
	}
	// Whether to check "l.safe" has no significant performance impact according to benchmark.
	if l.safe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	if len(l.elements) < 2 {
		return
	}
	slices.SortFunc(l.elements, cmp)
}

// Swap swaps the values at the given indexes.
func (l *List[T]) Swap(i, j int) {
	// Check "l.safe" before acquiring the lock is more efficient based on benchmark result.
	if l.safe {
		l.mu.Lock()
		defer l.mu.Unlock()
	}

	if l.withinRange(i, false) && l.withinRange(j, false) {
		l.elements[i], l.elements[j] = l.elements[j], l.elements[i]
	}
}

// Range call function fn on each value in the list.
// if `fn` returns false, the iteration stops.
// if `fn` is nil, the method does nothing.
func (l *List[T]) Range(fn func(v T) bool) {
	if fn == nil {
		return
	}
	// Whether to check "l.safe" has no significant performance impact according to benchmark.
	if l.safe {
		l.mu.RLock()
		defer l.mu.RUnlock()
	}

	for _, v := range l.elements {
		if !fn(v) {
			return
		}
	}
}

func (l *List[T]) resize(len, cap int) {
	newElements := make([]T, len, cap)
	copy(newElements, l.elements)
	l.elements = newElements
}

func (l *List[T]) growBy(n int) {
	currCap := cap(l.elements)
	newLen := len(l.elements) + n
	if newLen > currCap {
		// // method 1:
		// newCap := int(growthFactor * float32(currCap+n))
		// l.resize(newLen, newCap)

		// method 2:
		newCap := int(growthFactor * float32(currCap))
		for newCap < newLen {
			newCap = int(growthFactor * float32(newCap))
		}
		l.resize(newLen, newCap)
	} else {
		l.elements = l.elements[:newLen]
	}
}

func (l *List[T]) shrink() {
	currCap := cap(l.elements)
	if len(l.elements) <= int(shrinkFactor*float32(currCap)) {
		newCap := int(shrinkFactor * float32(currCap))
		if newCap < minCap {
			newCap = minCap
		}
		l.resize(len(l.elements), newCap)
	}
}

func (l *List[T]) withinRange(index int, allowEnd bool) bool {
	if allowEnd {
		return index >= 0 && index <= len(l.elements)
	}
	return index >= 0 && index < len(l.elements)
}

// todo: replace slices package method by myself logic.

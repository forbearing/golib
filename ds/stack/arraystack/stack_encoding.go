package arraystack

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// MarshalJSON will marshal the stack into a JSON-based representation.
func (s *Stack[E]) MarshalJSON() ([]byte, error) {
	el := make([]E, 0, s.list.Len())
	s.list.Range(func(e E) bool {
		el = append(el, e)
		return true
	})
	slices.Reverse(el)
	items := make([]string, 0, s.list.Len())
	for _, e := range el {
		b, err := json.Marshal(e)
		if err != nil {
			return nil, err
		}
		items = append(items, string(b))
	}
	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

// UnmarshalJSON will unmarshal a JSON-based representation byte slice into the stack.
func (s *Stack[E]) UnmarshalJSON(data []byte) (err error) {
	el := make([]E, 0)
	if err = json.Unmarshal(data, &el); err != nil {
		return err
	}
	s1, err := New(s.cmp, s.options()...)
	if err != nil {
		return err
	}
	(*s) = *s1
	for _, e := range el {
		s.Push(e)
	}
	return nil
}

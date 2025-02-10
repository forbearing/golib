package splaytree

import (
	"encoding/json"
)

// MarshalJSON implements the json.Marshaler interface.
func (t *Tree[K, V]) MarshalJSON() ([]byte, error) {
	if t.safe {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}

	m := make(map[K]V)
	var iter func(n *Node[K, V])
	iter = func(n *Node[K, V]) {
		if n == nil {
			return
		}
		iter(n.Children[0])
		m[n.Key] = n.Value
		iter(n.Children[1])
	}
	iter(t.root)

	return json.Marshal(m)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Tree[K, V]) UnmarshalJSON(data []byte) error {
	if t.safe {
		t.mu.Lock()
		defer t.mu.Unlock()
	}

	m := make(map[K]V)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	t.root = nil
	t.size = 0
	for k, v := range m {
		t.put(k, v)
	}
	return nil
}

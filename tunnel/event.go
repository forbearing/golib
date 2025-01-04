// package tunnel is communication protocol between with server and user, server and client.
package tunnel

import (
	"sync"

	"github.com/vmihailenco/msgpack/v5"
)

type Cmd uint32

var (
	Ping Cmd = NewCmd("ping", 1)
	Pong Cmd = NewCmd("pong", 2)
)

var cmdMap sync.Map

// NewCmd build a Cmd.
// custom always greater than 1000
func NewCmd(name string, value uint32) Cmd {
	cmdMap.Store(value, name)
	return Cmd(value)
}

func (c Cmd) String() string {
	if v, loaded := cmdMap.Load(uint32(c)); loaded {
		return v.(string)
	}
	return "Unknown"
}

type Event struct {
	ID    string `json:"id,omitempty" msgpack:"id,omitempty"`   // event id
	Cmd   Cmd    `json:"cmd,omitempty" msgpack:"cmd,omitempty"` // event cmd
	Error string `json:"error,omitempty" msgpack:"error,omitempty"`

	Payload any `json:"payload,omitempty" msgpack:"payload,omitempty"`
}

// DecodePayload converts Event.Payload to the specified type T.
// It handles three cases:
// 1. If the payload is nil, returns the zero value of type T
// 2. If the payload is already type T, returns it directly
// 3. Otherwise, uses msgpack marshal/unmarshal to perform the conversion
func DecodePayload[T any](data any) (T, error) {
	var result T
	if data == nil {
		return result, nil
	}
	if v, ok := data.(T); ok {
		return v, nil
	}
	b, err := msgpack.Marshal(data)
	if err != nil {
		return result, err
	}
	if err := msgpack.Unmarshal(b, &result); err != nil {
		return result, err
	}
	return result, nil
}

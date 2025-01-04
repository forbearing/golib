// package tunnel is communication protocol between with server and user, server and client.
package tunnel

import "sync"

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

// package tunnel is communication protocol between with server and user, server and client.
package tunnel

// Cmd is the command type
// NOTE: custom cmd must be greater than 1000.
type Cmd uint32

// NOTE:赋予具体的值的目的是防止 Cmd 值变化了, 如果 Cmd 值变化了, 相当于发送了其他
// 的命令, 这会直接导致 server 和 client 之间的通信出现问题.
const (
	Ping Cmd = 1
	Pong Cmd = 2
)

var cmdMap = map[Cmd]string{
	Ping: "Ping",
	Pong: "Pong",
}

func (c Cmd) String() string {
	v, ok := cmdMap[c]
	if ok {
		return v
	}
	return "Unknow"
}

func (c Cmd) Number() int { return int(c) }

type Event struct {
	ID    string `json:"id,omitempty" msgpack:"id,omitempty"`   // event id
	Cmd   Cmd    `json:"cmd,omitempty" msgpack:"cmd,omitempty"` // event cmd
	Error string `json:"error,omitempty" msgpack:"error,omitempty"`

	Payload any `json:"payload,omitempty" msgpack:"payload,omitempty"`
}

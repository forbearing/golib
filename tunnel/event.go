// package tunnel is communication protocol between with server and user, server and client.
package tunnel

// Cmd is the command type
type Cmd uint32

// NOTE:赋予具体的值的目的是防止 Cmd 值变化了, 如果 Cmd 值变化了, 相当于发送了其他
// 的命令, 这会直接导致 server 和 client 之间的通信出现问题.
const (
	RegisterAsControl Cmd = 1
	RegisterSuccess   Cmd = 2
	RegisterFailure   Cmd = 3
	ClientShutdown    Cmd = 5

	Heartbeat     Cmd = 1000
	ClientSysInfo Cmd = 1001

	Ping Cmd = 2001
	Pong Cmd = 2002
)

var CmdMap = map[Cmd]string{
	RegisterAsControl: "RegisterAsControl",
	RegisterSuccess:   "RegisterSuccess",
	RegisterFailure:   "RegisterFailure",
	ClientShutdown:    "ClientShutdown",

	Heartbeat:     "Heartbeat",
	ClientSysInfo: "ClientSysInfo",

	Ping: "Ping",
	Pong: "Pong",
}

func (c Cmd) String() string {
	v, ok := CmdMap[c]
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

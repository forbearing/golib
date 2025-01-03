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

func (c Cmd) Number() int {
	return int(c)
}

type Event struct {
	ID  string `json:"id,omitempty"`  // event id
	Cmd Cmd    `json:"cmd,omitempty"` // event cmd

	Error string `json:"error,omitempty"`

	// payload of event "RegisterAsControl"
	Hostname string    `json:"hostname,omitempty"`
	CtrlID   string    `json:"ctrl_id,omitempty"`
	HostInfo *HostInfo `json:"host_info,omitempty"`

	ClientSysInfoPayload *ClientSysInfoPayload `json:"client_sysinfo_payload,omitempty"` // CientSysInfo Payload
}

type HostInfo struct {
	Hostname             string `json:"hostname" mapstructure:"hostname"`
	Uptime               uint64 `json:"uptime" mapstructure:"uptime"`
	BootTime             uint64 `json:"bootTime" mapstructure:"bootTime"`
	Procs                uint64 `json:"procs" mapstructure:"procs"`                     // number of processes
	OS                   string `json:"os" mapstructure:"os"`                           // ex: freebsd, linux
	Platform             string `json:"platform" mapstructure:"platform"`               // ex: ubuntu, linuxmint
	PlatformFamily       string `json:"platformFamily" mapstructure:"platformFamily"`   // ex: debian, rhel
	PlatformVersion      string `json:"platformVersion" mapstructure:"platformVersion"` // version of the complete OS
	KernelVersion        string `json:"kernelVersion" mapstructure:"kernelVersion"`     // version of the OS kernel (if available)
	KernelArch           string `json:"kernelArch" mapstructure:"kernelArch"`           // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
	VirtualizationSystem string `json:"virtualizationSystem" mapstructure:"virtualizationSystem"`
	VirtualizationRole   string `json:"virtualizationRole" mapstructure:"virtualizationRole"` // guest or host
	HostID               string `json:"hostId" mapstructure:"hostId"`                         // ex: uuid

	AllIPv4 []string `json:"allIpv4" mapstructure:"allIpv4"`
	AllIPv6 []string `json:"allIpv6" mapstructure:"allIpv6"`
	Version string   `json:"clientVersion" mapstructure:"clientVersion"`
}

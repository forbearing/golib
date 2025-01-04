package tunnel_test

import (
	"net"
	"testing"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/tunnel"
	"github.com/forbearing/golib/types/consts"
	"github.com/stretchr/testify/assert"
)

var (
	addr   = "0.0.0.0:12345"
	doneCh = make(chan struct{}, 1)

	Bye   = tunnel.NewCmd("bye", 1000)
	Hello = tunnel.NewCmd("hello", 1001)
)

type ByePaylod struct {
	Field1 string
	Field2 uint64
}
type HelloPaylod struct {
	Field3 string
	Field4 float64
}

var (
	byePayload1   = ByePaylod{Field1: "bye1", Field2: 123}
	byePayload2   = ByePaylod{Field1: "bye2", Field2: 456}
	helloPayload1 = HelloPaylod{Field3: "hello1", Field4: 3.14}
	helloPayload2 = HelloPaylod{Field3: "hello2", Field4: 3.14}
)

func TestSession(t *testing.T) {
	assert.NoError(t, bootstrap.Bootstrap())
	go server(t)
	client(t)
	<-doneCh
}

func server(t *testing.T) {
	l, err := net.Listen("tcp", addr)
	assert.NoError(t, err)
	defer l.Close()

	doneCh <- struct{}{}

	conn, err := l.Accept()
	assert.NoError(t, err)
	defer conn.Close()

	session, _ := tunnel.NewSession(conn, consts.Server)
	for {
		event, err := session.Read()
		assert.NoError(t, err)
		switch event.Cmd {
		case tunnel.Ping:
			t.Log("client ping")
			session.Write(&tunnel.Event{Cmd: tunnel.Pong})
		case Hello:
			payload, err := tunnel.DecodePayload[HelloPaylod](event.Payload)
			assert.NoError(t, err)
			t.Logf("client hello: %+v\n", payload)
			session.Write(&tunnel.Event{Cmd: Hello, Payload: helloPayload2})
		case Bye:
			payload, err := tunnel.DecodePayload[ByePaylod](event.Payload)
			assert.NoError(t, err)
			t.Logf("client bye: %+v\n", payload)
			session.Write(&tunnel.Event{Cmd: Bye, Payload: byePayload2})
		}
	}
}

func client(t *testing.T) {
	<-doneCh

	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	defer conn.Close()

	session, _ := tunnel.NewSession(conn, consts.Client)
	session.Write(&tunnel.Event{Cmd: tunnel.Ping})

	for {
		event, err := session.Read()
		assert.NoError(t, err)

		switch event.Cmd {
		case tunnel.Pong:
			t.Log("server pong")
			session.Write(&tunnel.Event{Cmd: Hello, Payload: helloPayload1})
		case Hello:
			payload, err := tunnel.DecodePayload[HelloPaylod](event.Payload)
			assert.NoError(t, err)
			t.Logf("server hello: %+v\n", payload)
			session.Write(&tunnel.Event{Cmd: Bye, Payload: byePayload1})
		case Bye:
			payload, err := tunnel.DecodePayload[ByePaylod](event.Payload)
			assert.NoError(t, err)
			t.Logf("server bye: %+v\n", payload)
			doneCh <- struct{}{}
			return
		}
	}
}

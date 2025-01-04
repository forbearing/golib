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
			t.Log("client hello")
			session.Write(&tunnel.Event{Cmd: Hello})
		case Bye:
			t.Log("client bye")
			session.Write(&tunnel.Event{Cmd: Bye})
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
			session.Write(&tunnel.Event{Cmd: Hello})
		case Hello:
			t.Log("server hello")
			session.Write(&tunnel.Event{Cmd: Bye})
		case Bye:
			t.Log("server bye")
			doneCh <- struct{}{}
			return
		}
	}
}

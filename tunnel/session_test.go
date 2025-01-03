package tunnel_test

import (
	"net"
	"testing"
	"time"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/tunnel"
	"github.com/forbearing/golib/types/consts"
	"github.com/stretchr/testify/assert"
)

var (
	addr   = "0.0.0.0:12345"
	doneCh = make(chan struct{}, 1)
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

	session, err := tunnel.NewSession(conn, consts.Server)
	assert.NoError(t, err)

	event, err := session.Read()
	assert.NoError(t, err)
	switch event.Cmd {
	case tunnel.Ping:
		t.Log("client ping")
		session.Write(&tunnel.Event{Cmd: tunnel.Pong})
	}
}

func client(t *testing.T) {
	<-doneCh

	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	defer conn.Close()

	session, err := tunnel.NewSession(conn, consts.Client)
	assert.NoError(t, err)
	assert.NoError(t, session.Write(&tunnel.Event{Cmd: tunnel.Ping}))

	for {
		event, err := session.Read()
		assert.NoError(t, err)
		if event.Cmd == tunnel.Pong {
			t.Log("server pong")
			break
		}
		time.Sleep(300 * time.Millisecond)
	}

	doneCh <- struct{}{}
}

package tunnel

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/types/consts"
	"github.com/forbearing/golib/util"
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
)

// connMap is concurrency secure map.
// connMap key is net.Conn's file descriptor(int)
// connMap value is *connMapVal
var connMap sync.Map

// connMapVal to ensure a net.Conn has only one corresponding locker.
type connMapVal struct {
	key any

	// rmu used by Read function to ensure read binary data and data size
	// is an atomic operation, otherwise dirty data will be generated.
	rmu sync.Mutex
	// wmu used by Write function to write binary data and data size
	// is an atomic operation, otherwise dirty data will be generated.
	wmu sync.Mutex

	// connection contains the connection of remote/local ip and remote/local port.
	connection util.Connection
}

// List will list all tunnel session connection.
func List() map[any]any {
	kv := make(map[any]any)
	connMap.Range(func(key, val any) bool {
		kv[key] = val
		return true
	})
	return kv
}

type Session struct {
	tcpconn net.Conn
	wsconn  *websocket.Conn

	locker *connMapVal
}

// Read reads binary data and the size of binary data.
// Read should be locked to ensure that reading data and size of data
// is an atomic operation.
// FIXME: handle io.EOF error
func (s *Session) Read() (*Event, error) {
	// logger.Protocol.Infow("Session.Read() S", "key", s.locker.key) // S: start
	// s.locker.rmu.Lock()
	// logger.Protocol.Infow("Session.Read() L", "key", s.locker.key) // L: lock
	// defer s.locker.rmu.Unlock()

	// cmdBuf, err := internal.ReadBinary(s.tcpconn)
	// if err != nil {
	// 	return nil, err
	// }

	// 1.read data size
	var size uint32
	if err := binary.Read(s.tcpconn, binary.BigEndian, &size); err != nil {
		return nil, err
	}
	if size == 0 {
		return nil, errors.New("data size is nil")
	}
	// 2.read data
	cmdBuf := make([]byte, size)
	if _, err := io.ReadFull(s.tcpconn, cmdBuf); err != nil {
		return nil, err
	}
	logger.Binary.Infoz("Read", zap.Uint32("size", size), zap.ByteString("data", cmdBuf), zap.Binary("binary", cmdBuf))
	// 3.unmarshal data.
	event := &Event{}
	if err := msgpack.Unmarshal(cmdBuf, event); err != nil {
		return nil, err
	}

	logger.Protocol.Infow("Session.Read()", "key", s.locker.key, "event", event.Cmd)
	return event, nil
}

// Write writes binary data and the size of binary data.
// Write should be locked to ensure that writing data and size of data
// is an atomic operation.
// FIXME: handle io.EOF error
func (s *Session) Write(event *Event) error {
	if event == nil {
		return nil
	}
	logger.Protocol.Infow("Session.Write() S", "key", s.locker.key, "event", event.Cmd) // S: start
	s.locker.wmu.Lock()
	logger.Protocol.Infow("Session.Write() L", "key", s.locker.key, "event", event.Cmd) // L: lock
	defer s.locker.wmu.Unlock()
	// 这个 ID 极其重要, 如果人家设置了 ID, 你千万别动人家的 I
	// 因为 Event 的 ID 相同就表明这是正在进行的同一个通信.
	// 如果 Event 的 ID 就表明这不是同一个通信.
	if event.ID == "" {
		event.ID = util.UUID()
	}
	// 1.marshal data
	buf, err := msgpack.Marshal(event)
	if err != nil {
		return err
	}

	// 2.write data size
	if err := binary.Write(s.tcpconn, binary.BigEndian, uint32(len(buf))); err != nil {
		return err
	}
	logger.Binary.Infoz("Write", zap.Int("size", len(buf)), zap.ByteString("data", buf), zap.Binary("binary", buf))
	// 3.write data
	if err := binary.Write(s.tcpconn, binary.BigEndian, buf); err != nil {
		return err
	}
	return nil
}

// NewSession
// cid is the client control connection uuid.
// server-side should always provides the cid when send or receive event.
func NewSession(_conn any, appSide consts.AppSide, cid ...string) (*Session, error) {
	switch conn := _conn.(type) {
	case *net.TCPConn:
		var key any
		var nd util.Connection
		switch appSide {
		case consts.Server:
			// // NOTE:file descriptor is the key in server-side app.
			// // remote port could be same when too many clients.
			// file, err := conn.File()
			// if err != nil {
			// 	panic(err)
			// }
			// file.Close()
			// key = int(file.Fd())

			nd = util.GetConnection(conn)
			key = nd.RemoteIP
			if len(cid) > 0 {
				if len(cid[0]) > 0 {
					key = cid[0]
				}
			}
			// logger.Protocol.Infow(key.(string), "rip", nd.Rip, "rport", nd.Rport, "lip", nd.Lip, "lport", nd.Lport)
		case consts.Client:
			// NOTE:local port is the key in client-side app.
			// local port always is server listen port in server-side.
			nd = util.GetConnection(conn)
			key = nd.LocalPort
			// logger.Protocol.Infow(strconv.Itoa(key.(int)), "rip", nd.Rip, "rport", nd.Rport, "lip", nd.Lip, "lport", nd.Lport)
		default:
			panic("unknow app side: " + appSide)
		}
		locker := new(connMapVal)
		locker.connection = nd
		locker.key = key
		if val, ok := connMap.LoadOrStore(key, locker); ok {
			locker = val.(*connMapVal)
		}
		return &Session{tcpconn: conn, locker: locker}, nil
	case *websocket.Conn:
		// return &Session{wsconn: conn}, nil
		return nil, errors.New("only support tcp.Conn")
	case websocket.Conn:
		// return &Session{wsconn: &conn}, nil
		return nil, errors.New("only support tcp.Conn")
	default:
		return nil, errors.New("only support tcp.Conn")
	}
}

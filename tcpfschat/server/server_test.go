package server

import (
	"github.com/merisho/tcp-fs-chat/test"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func NewTestListener() *TestListener {
	return &TestListener{
		acceptConnections: make(chan net.Conn, 4096),
	}
}

type TestListener struct {
	acceptConnections chan net.Conn
}

func (l *TestListener) connectionsToAccept(conns ...net.Conn) {
	for _, c := range conns {
		l.acceptConnections <- c
	}
}

func (l *TestListener) Accept() (net.Conn, error) {
	return <- l.acceptConnections, nil
}

func (l *TestListener) Close() error {
	return nil
}

func (l *TestListener) Addr() net.Addr {
	return nil
}

func TestAcceptClient(t *testing.T) {
	ln := NewTestListener()
	ln.connectionsToAccept(test.NewTestConnection())

	s := NewServer(ln)
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 1, s.ConnectionCount())
}

func TestHandleConnectionClosing(t *testing.T) {
	ln := NewTestListener()
	conn := test.NewTestConnection().EOFOnRead()
	ln.connectionsToAccept(conn)

	s := NewServer(ln)
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 0, s.ConnectionCount())
	assert.True(t, conn.Closed())
}

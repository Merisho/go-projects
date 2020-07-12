package connections

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

type TestConn struct {
	LastWrite []byte
}

func (c *TestConn) Read(b []byte) (n int, err error) {
	return len(b), nil
}

func (c *TestConn) Write(b []byte) (n int, err error) {
	c.LastWrite = b
	return len(b), nil
}

func (c *TestConn) Close() error {
	return nil
}

func (c *TestConn) LocalAddr() net.Addr {
	return nil
}

func (c *TestConn) RemoteAddr() net.Addr {
	return nil
}

func (c *TestConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *TestConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *TestConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestConnectionsBroadcast(t *testing.T) {
	conns := &Connections{}

	conn1 := &TestConn{}
	conn2 := &TestConn{}
	conn3 := &TestConn{}

	conns.Add(conn1)
	conns.Add(conn2)
	conns.Add(conn3)

	b := []byte("test")
	err := conns.Broadcast(b)

	assert.NoError(t, err)
	assert.Equal(t, b, conn1.LastWrite)
	assert.Equal(t, b, conn2.LastWrite)
	assert.Equal(t, b, conn3.LastWrite)
}

func TestConnectionsBroadcastFrom(t *testing.T) {
	conns := &Connections{}

	conn1 := &TestConn{}
	conn2 := &TestConn{}
	conn3 := &TestConn{}

	conns.Add(conn1)
	conns.Add(conn2)
	conns.Add(conn3)

	b := []byte("test")
	err := conns.BroadcastFrom(conn1, b)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(conn1.LastWrite))
	assert.Equal(t, b, conn2.LastWrite)
	assert.Equal(t, b, conn3.LastWrite)
}

func TestRemoveConnection(t *testing.T) {
	conns := &Connections{}
	conn := &TestConn{}
	conns.Add(conn)

	conns.Remove(conn)

	assert.Equal(t, 0, len(conns.conns))
}

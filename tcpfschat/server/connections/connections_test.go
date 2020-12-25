package connections

import (
	"github.com/merisho/tcp-fs-chat/test"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"testing"
)

type TestErr struct {
	Temp bool
	Time bool
	Err string
}

func (e TestErr) Temporary() bool {
	return e.Temp
}

func (e TestErr) Error() string {
	return e.Err
}

func (e TestErr) Timeout() bool {
	return e.Time
}

func TestHandleConnectionClosing(t *testing.T) {
	conns := &Connections{}
	conn, _ := conns.Add(test.NewTestConnection())

	connOk := conns.HandleConnectionErr(conn, io.EOF)

	assert.False(t, connOk)
	assert.Equal(t, 0, conns.Count())
}

func TestHandleNilConnectionError(t *testing.T) {
	conns := &Connections{}

	connOk := conns.HandleConnectionErr(nil, nil)

	assert.True(t, connOk)
}

func TestHandleNonTemporaryError(t *testing.T) {
	conns := &Connections{}
	conn, _ := conns.Add(test.NewTestConnection())

	connOk := conns.HandleConnectionErr(conn, TestErr{})

	assert.False(t, connOk)
	assert.Equal(t, 0, conns.Count())
}

func TestHandleTemporaryError(t *testing.T) {
	conns := &Connections{}
	conn, _ := conns.Add(test.NewTestConnection())

	connOk := conns.HandleConnectionErr(conn, TestErr{Temp: true})

	assert.True(t, connOk)
	assert.Equal(t, 1, conns.Count())
}

func TestHandleTimeoutError(t *testing.T) {
	conns := &Connections{}
	conn, _ := conns.Add(test.NewTestConnection())

	connOk := conns.HandleConnectionErr(conn, TestErr{Time: true})

	assert.False(t, connOk)
	assert.Equal(t, 0, conns.Count())
}

func TestGenerateUUIDOnConnectionAdd(t *testing.T) {
	conns := &Connections{}
	c := test.NewTestConnection()

	conn, err := conns.Add(c)

	assert.NoError(t, err)
	assert.Len(t, conn.ID(), 16)
}

func TestRemoveByID(t *testing.T) {
	conns := &Connections{}

	conn, _ := conns.Add(test.NewTestConnection())

	removed := conns.RemoveByID(conn.ID())
	assert.Equal(t, conn, removed)
	assert.Zero(t, conns.Count())

	nonExistentID := []byte{1,2,3,4,5,6,7}
	removed = conns.RemoveByID(nonExistentID)
	assert.Nil(t, removed)
}

func TestForEach(t *testing.T) {
	conns := &Connections{}

	conn1, conn2 := net.Pipe()
	conn3, conn4 := net.Pipe()

	conns.Add(conn1)
	conns.Add(conn2)
	conns.Add(conn3)
	conns.Add(conn4)

	count := make(chan struct{}, 4)
	conns.ForEach(func(c Conn) {
		count <- struct{}{}
	})

	assert.Len(t, count, 4)
}

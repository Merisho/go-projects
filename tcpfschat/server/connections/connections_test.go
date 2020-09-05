package connections

import (
	"github.com/merisho/tcp-fs-chat/test"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestConnectionsBroadcast(t *testing.T) {
	conns := &Connections{}

	conn1 := test.NewTestConnection()
	conn2 := test.NewTestConnection()
	conn3 := test.NewTestConnection()

	conns.Add(conn1)
	conns.Add(conn2)
	conns.Add(conn3)

	msg := "test"
	err := conns.Broadcast([]byte(msg))

	assert.NoError(t, err)
	assert.Equal(t, msg, conn1.FrontWrittenChunk())
	assert.Equal(t, msg, conn2.FrontWrittenChunk())
	assert.Equal(t, msg, conn3.FrontWrittenChunk())
}

func TestConnectionsBroadcastFrom(t *testing.T) {
	conns := &Connections{}

	conn1 := test.NewTestConnection()
	conn2 := test.NewTestConnection()
	conn3 := test.NewTestConnection()

	conns.Add(conn1)
	conns.Add(conn2)
	conns.Add(conn3)

	msg := "test"
	err := conns.BroadcastFrom(conn1, []byte(msg))

	assert.NoError(t, err)
	assert.Equal(t, "", conn1.FrontWrittenChunk())
	assert.Equal(t, msg, conn2.FrontWrittenChunk())
	assert.Equal(t, msg, conn3.FrontWrittenChunk())
}

func TestRemoveConnection(t *testing.T) {
	conns := &Connections{}
	conn := test.NewTestConnection()
	conns.Add(conn)

	conns.Remove(conn)

	assert.Equal(t, 0, len(conns.conns))
}

func TestHandleConnectionClosing(t *testing.T) {
	conns := &Connections{}
	conn := test.NewTestConnection()
	conns.Add(conn)

	connOk := conns.HandleConnectionErr(conn, io.EOF)

	assert.False(t, connOk)
	assert.Equal(t, 0, conns.Count())
}

func TestHandleNilConnectionError(t *testing.T) {
	conns := &Connections{}

	connOk := conns.HandleConnectionErr(nil, nil)

	assert.True(t, connOk)
}

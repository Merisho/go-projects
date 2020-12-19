package connections

import (
	"github.com/merisho/tcp-fs-chat/test"
	"github.com/stretchr/testify/assert"
	"io"
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

func TestConnectionsBroadcast(t *testing.T) {
	conns := &Connections{}

	conn1 := test.NewTestConnection()
	conn2 := test.NewTestConnection()
	conn3 := test.NewTestConnection()

	conns.Add(conn1)
	conns.Add(conn2)
	conns.Add(conn3)

	msg := "test"
	errs := conns.Broadcast([]byte(msg))

	assert.Empty(t, errs)
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
	errs := conns.BroadcastFrom(conn1, []byte(msg))

	assert.Empty(t, errs)
	assert.Equal(t, "", conn1.FrontWrittenChunk())
	assert.Equal(t, msg, conn2.FrontWrittenChunk())
	assert.Equal(t, msg, conn3.FrontWrittenChunk())
}

func TestBroadcastErrors(t *testing.T) {
	conns := &Connections{}

	conn1 := test.NewTestConnection().ErrorOnWrite()
	conn2 := test.NewTestConnection().ErrorOnWrite()
	conns.Add(conn1)
	conns.Add(conn2)

	connErrs := conns.Broadcast([]byte("Hello"))

	assert.Equal(t, 2, len(connErrs))
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

func TestHandleNonTemporaryError(t *testing.T) {
	conns := &Connections{}
	conn := test.NewTestConnection()
	conns.Add(conn)

	connOk := conns.HandleConnectionErr(conn, TestErr{})

	assert.False(t, connOk)
	assert.Equal(t, 0, conns.Count())
}

func TestHandleTemporaryError(t *testing.T) {
	conns := &Connections{}
	conn := test.NewTestConnection()
	conns.Add(conn)

	connOk := conns.HandleConnectionErr(conn, TestErr{Temp: true})

	assert.True(t, connOk)
	assert.Equal(t, 1, conns.Count())
}

func TestHandleTimeoutError(t *testing.T) {
	conns := &Connections{}
	conn := test.NewTestConnection()
	conns.Add(conn)

	connOk := conns.HandleConnectionErr(conn, TestErr{Time: true})

	assert.False(t, connOk)
	assert.Equal(t, 0, conns.Count())
}

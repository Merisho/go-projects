package client_test

import (
    "github.com/merisho/tcp-fs-chat/client"
    "github.com/merisho/tcp-fs-chat/test"
    "github.com/stretchr/testify/assert"
    "net"
    "testing"
    "time"
)

func TestReceiveMessage(t *testing.T) {
    conn := test.NewTestConnection()
    c := client.New(conn, nil)

    r := c.Receive()
    conn.ChunksToRead("Hello")

    chanMessageEqual(t, r, "Hello")
}

func TestSkipNotReadyReceiver(t *testing.T) {
    conn := test.NewTestConnection()
    c := client.New(conn, nil)

    _ = c.Receive()
    ready := c.Receive()
    conn.ChunksToRead("Hello")

    chanMessageEqual(t, ready, "Hello")
}

func TestSendMessage(t *testing.T) {
    conn := test.NewTestConnection()
    c := client.New(conn, nil)

    err := c.Send("Hello")
    assert.NoError(t, err)

    // Client appends 0 byte to the end of each message
    assert.Equal(t, "Hello\x00", conn.FrontWrittenChunk())
}

func TestDisconnectFromServer(t *testing.T) {
    conn := test.NewTestConnection().EOFOnRead()
    c := client.New(conn, nil)

    _, ok := <-c.Receive()
    assert.False(t, ok)
}

func TestSplitTCPDataChunkIntoMessages(t *testing.T) {
    testConn, conn := net.Pipe()
    c := client.New(conn, nil)

    msg := []byte("Hello\x00World\x00")
    _, _ = testConn.Write(msg)

    r := c.Receive()
    assert.Equal(t, "Hello", <-r)
    assert.Equal(t, "World", <-r)
}

func chanMessageEqual(t *testing.T, c chan string, expected string) {
    select {
    case m := <- c:
        assert.Equal(t, expected, m)
    case <- time.After(5 * time.Millisecond):
        assert.Fail(t, "timeout")
    }
}

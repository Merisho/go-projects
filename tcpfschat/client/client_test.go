package client_test

import (
    "github.com/merisho/tcp-fs-chat/client"
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

func TestReceiveMessage(t *testing.T) {
    conn := NewTestConnection()
    c := client.New(conn)

    r := c.Receive()
    conn.chunksToRead("Hello")

    chanMessageEqual(t, r, "Hello")
}

func TestSkipNotReadyReceiver(t *testing.T) {
    conn := NewTestConnection()
    c := client.New(conn)

    _ = c.Receive()
    ready := c.Receive()
    conn.chunksToRead("Hello")

    chanMessageEqual(t, ready, "Hello")
}

func TestRetryDeliveryToNotReadyReceiver(t *testing.T) {
    conn := NewTestConnection()
    c := client.New(conn)

    notReady := c.Receive()
    ready := c.Receive()
    conn.chunksToRead("Hello")

    chanMessageEqual(t, ready, "Hello")
    chanMessageEqual(t, notReady, "Hello")
}

func TestRemoveNotReadyReceiversAfterDeliveryRetryTimeout(t *testing.T) {
    conn := NewTestConnection()
    c := client.New(conn)

    notReady := c.Receive()
    conn.chunksToRead("Hello")

    time.Sleep(15 * time.Millisecond)

    conn.chunksToRead("Hello")

    select {
    case <- time.After(10 * time.Millisecond):
    case _, ok := <- notReady:
        assert.False(t, ok, "must NOT receive the 'Hello' message, but instead a channel closing signal")
    }
}

func TestSendMessage(t *testing.T) {
    conn := NewTestConnection()
    c := client.New(conn)

    err := c.Send("Hello")
    assert.NoError(t, err)

    assert.Equal(t, "Hello", conn.frontWrittenChunk())
}

func chanMessageEqual(t *testing.T, c chan string, expected string) {
    select {
    case m := <- c:
        assert.Equal(t, expected, m)
    case <- time.After(5 * time.Millisecond):
        assert.Fail(t, "timeout")
    }
}

package test

import (
	"github.com/merisho/tcp-fs-chat/client"
	"github.com/merisho/tcp-fs-chat/server"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChat(t *testing.T) {
	_, err := server.Serve(1337)
	assert.NoError(t, err)

	c1, err := client.Connect("localhost", 1337)
	assert.NoError(t, err)

	c2, err := client.Connect("localhost", 1337)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	c1.Send("Hello from client 1")

	time.Sleep(10 * time.Millisecond)
	select {
	case msg := <-c2.Receive():
		assert.Equal(t, "Hello from client 1", msg)
	case <-time.After(1 * time.Second):
		assert.Fail(t, "timeout")
	}

	c2.Send("Hello from client 2")

	select {
	case msg := <-c1.Receive():
		assert.Equal(t, "Hello from client 1", msg)
	case <-time.After(1 * time.Second):
		assert.Fail(t, "timeout")
	}

	select {
	case msg := <-c1.Receive():
		assert.Equal(t, "Hello from client 2", msg)
	case <-time.After(1 * time.Second):
		assert.Fail(t, "timeout")
	}
}

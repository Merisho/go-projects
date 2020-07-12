package client_test

import (
	"github.com/merisho/tcp-fs-chat/client"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClientConnect(t *testing.T) {
	port, _ := StartFakeServer()

	c, err := client.Connect("localhost", port)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestReceiveMessage(t *testing.T) {
	port, msgs := StartFakeServer()

	c, err := client.Connect("localhost", port)
	assert.NoError(t, err)

	r := c.Receive()

	msgs <- "Hello test"

	select {
	case m := <- r:
		assert.Equal(t, "Hello test", m)
	case <- time.After(300 * time.Millisecond):
		assert.Fail(t, "timeout")
	}
}

func TestAuth(t *testing.T) {
	port, msgs := StartFakeServer()

	c, err := client.Connect("localhost", port)
	assert.NoError(t, err)

	time.AfterFunc(10 * time.Millisecond, func() {
		msgs <- "auth success"
	})

	err = c.Auth("test", "password")
	assert.NoError(t, err)

	time.AfterFunc(10 * time.Millisecond, func() {
		msgs <- "auth fail"
	})

	err = c.Auth("test", "invalid_password")
	assert.Error(t, err)
}

//func TestSendMessage(t *testing.T) {
//	port, msgs, srvMsgs := StartFakeServer()
//
//	c, err := client.Connect("localhost", port)
//	assert.NoError(t, err)
//
//	c.Send("test")
//}

package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	_, addr, stop := startTestServer()
	defer stop()

	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	n, err := conn.Write([]byte("Hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
}

func TestBroadcast(t *testing.T) {
	_, addr, stop := startTestServer()
	defer stop()

	sender := createClient(addr, "sender")
	receiver1 := createClient(addr, "receiver1")
	receiver2 := createClient(addr, "receiver2")

	_, err := sender.Write([]byte("hello"))
	assert.NoError(t, err)

	res := make([]byte, 5)

	_, err = receiver1.Read(res)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(res))

	_, err = receiver2.Read(res)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(res))
}

func TestRemoveClosedConnection(t *testing.T) {
	server, addr, stop := startTestServer()
	defer stop()

	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	_, err = conn.Write([]byte("test::test"))
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 1, server.ConnectionCount())

	err = conn.Close()
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 0, server.ConnectionCount())
}

func TestClientAuthenticationFail(t *testing.T) {
	_, addr, stop := startTestServer()
	defer stop()

	client, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	_, err = client.Write([]byte("Hello"))
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	res := make([]byte, 9)
	_, err = client.Read(res)
	assert.Equal(t, "auth fail", string(res))

	_, err = client.Read(res)
	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
}

func TestClientAuthenticationSuccess(t *testing.T) {
	_, addr, stop := startTestServer()
	defer stop()

	client, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	_, err = client.Write([]byte("test::123456"))
	assert.NoError(t, err)

	res := make([]byte, 12)
	_, err = client.Read(res)
	assert.Equal(t, "auth success", string(res))
}

func TestAuthCreds(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Error(e)
		}
	}()

	s := &Server{}

	cases := []struct{
		credsMessage []byte
		username string
		password string
	}{
		{
			credsMessage: append([]byte("test::123123"), 0, 0),
			username: "test",
			password: "123123",
		},
		{
			credsMessage: append([]byte{0}, []byte("test::123123")...),
			username: "test",
			password: "123123",
		},
		{
			credsMessage: []byte("::"),
			username: "",
			password: "",
		},
		{
			credsMessage: []byte(""),
			username: "",
			password: "",
		},
		{
			credsMessage: []byte("te::st:123123"),
			username: "te",
			password: "st:123123",
		},
		{
			credsMessage: []byte(":"),
			username: "",
			password: "",
		},
	}

	for i, c := range cases {
		username, password := s.authCreds(c.credsMessage)

		require.Equal(t, c.username, username, "Test: %d", i + 1)
		require.Equal(t, c.password, password, "Test: %d", i + 1)
	}
}

func startTestServer() (s *Server, addr string, stopServer func()) {
	p, addr := getAddr()
	server, err := Serve(p)
	if err != nil {
		panic(err)
	}

	return server, addr, func() {
		stop(server)
	}
}

func stop(s *Server) {
	if err := recover(); err != nil {
		fmt.Println(err)
	}

	err := s.Close()
	if err != nil {
		panic(err)
	}
}

var port = portGen()

func getAddr() (uint16, string) {
	p := port()
	return p, fmt.Sprintf("localhost:%d", p)
}

func portGen() func() uint16 {
	base := uint32(1336)
	return func() uint16 {
		if base == 65535 {
			panic("port out of range")
		}

		return uint16(atomic.AddUint32(&base, 1))
	}
}

func createClient(addr, name string) net.Conn {
	client, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	_, err = client.Write([]byte(name + "::" + name))
	if err != nil {
		panic(err)
	}

	res := make([]byte, 12)
	_, err = client.Read(res)
	if err != nil {
		panic(err)
	}

	if string(res) == "auth success" {
		return client
	}

	panic(string(res))
}

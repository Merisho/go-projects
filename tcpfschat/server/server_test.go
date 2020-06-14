package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	p, addr := getAddr()
	server, err := Serve(p)
	assert.NoError(t, err)

	defer stop(server)

	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	n, err := conn.Write([]byte("Hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
}

func TestBroadcast(t *testing.T) {
	p, addr := getAddr()
	server, err := Serve(p)
	assert.NoError(t, err)

	defer stop(server)

	sender, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	receiver1, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	receiver2, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	_, err = sender.Write([]byte("Hello test"))
	assert.NoError(t, err)

	err = receiver1.SetDeadline(time.Now().Add(1 * time.Second))
	assert.NoError(t, err)

	b := make([]byte, 1024)
	_, err = receiver1.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, "Hello test", strings.TrimRight(string(b), "\x00"))

	err = receiver2.SetDeadline(time.Now().Add(1 * time.Second))
	assert.NoError(t, err)

	_, err = receiver2.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, "Hello test", strings.TrimRight(string(b), "\x00"))
}

func TestRemoveClosedConnection(t *testing.T) {
	p, addr := getAddr()
	server, err := Serve(p)
	assert.NoError(t, err)

	defer stop(server)

	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)

	err = conn.Close()
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 0, server.ConnectionCount())
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

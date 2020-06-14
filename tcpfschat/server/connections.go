package server

import (
	"net"
	"sync"
)

type Connections struct {
	conns []net.Conn
	mu sync.Mutex
}

func (conns *Connections) Add(c net.Conn) {
	conns.mu.Lock()
	conns.conns = append(conns.conns, c)
	conns.mu.Unlock()
}

func (conns *Connections) Remove(c net.Conn) {
	conns.mu.Lock()
	defer conns.mu.Unlock()
	for i, conn := range conns.conns {
		if conn == c {
			conns.conns = append(conns.conns[:i], conns.conns[i + 1:]...)
			break
		}
	}
}

func (conns *Connections) Broadcast(b []byte) error {
	return conns.broadcastFrom(nil, b)
}

func (conns *Connections) BroadcastFrom(from net.Conn, b []byte) error {
	return conns.broadcastFrom(from, b)
}

func (conns *Connections) broadcastFrom(from net.Conn, b []byte) error {
	conns.mu.Lock()
	defer conns.mu.Unlock()

	var err error
	for _, c := range conns.conns {
		if c == from {
			continue
		}

		_, e := c.Write(b)
		if e != nil {
			err = e
		}
	}

	return err
}

func (conns *Connections) Count() int {
	conns.mu.Lock()
	defer conns.mu.Unlock()
	return len(conns.conns)
}

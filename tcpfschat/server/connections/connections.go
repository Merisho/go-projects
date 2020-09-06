package connections

import (
	"io"
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

func (conns *Connections) Broadcast(b []byte) []ConnErr {
	return conns.broadcastFrom(nil, b)
}

func (conns *Connections) BroadcastFrom(from net.Conn, b []byte) []ConnErr {
	return conns.broadcastFrom(from, b)
}

func (conns *Connections) broadcastFrom(from net.Conn, b []byte) []ConnErr {
	conns.mu.Lock()
	defer conns.mu.Unlock()

	var errs []ConnErr
	for _, c := range conns.conns {
		if c == from {
			continue
		}

		_, err := c.Write(b)
		if err != nil {
			errs = append(errs, ConnErr{
				Conn: c,
				Err: err,
			})
		}
	}

	return errs
}

func (conns *Connections) Count() int {
	conns.mu.Lock()
	defer conns.mu.Unlock()
	return len(conns.conns)
}

func (conns *Connections) HandleConnectionErr(c net.Conn, err error) (connectionOk bool) {
	if err == nil {
		return true
	}

	if err == io.EOF {
		conns.Remove(c)
		return false
	}

	if e, ok := err.(net.Error); ok {
		if e.Temporary() {
			return true
		}

		conns.Remove(c)
		return false
	}

	return true
}

type ConnErr struct {
	Conn net.Conn
	Err error
}

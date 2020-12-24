package connections

import (
	"bytes"
	"io"
	"net"
	"sync"
	"github.com/google/uuid"
)

type Connections struct {
	conns []Conn
	mu sync.Mutex
}

func (conns *Connections) Add(c net.Conn) (conn Conn, err error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	conn = newConnection(c, uid[:])

	conns.mu.Lock()
	conns.conns = append(conns.conns, conn)
	conns.mu.Unlock()

	return conn, nil
}

func (conns *Connections) RemoveByID(id []byte) Conn {
	conns.mu.Lock()
	defer conns.mu.Unlock()

	var removed Conn
	for i, conn := range conns.conns {
		if bytes.Equal(conn.ID(), id) {
			removed = conn
			conns.conns = append(conns.conns[:i], conns.conns[i + 1:]...)
			break
		}
	}

	return removed
}

func (conns *Connections) Broadcast(b []byte) []ConnErr {
	return conns.broadcastFrom(nil, b)
}

func (conns *Connections) BroadcastFrom(from Conn, b []byte) []ConnErr {
	return conns.broadcastFrom(from, b)
}

func (conns *Connections) broadcastFrom(from Conn, b []byte) []ConnErr {
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

func (conns *Connections) HandleConnectionErr(c Conn, err error) (connectionOk bool) {
	if err == nil {
		return true
	}

	if err == io.EOF {
		conns.RemoveByID(c.ID())
		return false
	}

	if e, ok := err.(net.Error); ok {
		if e.Temporary() {
			return true
		}

		conns.RemoveByID(c.ID())
		return false
	}

	return true
}

type ConnErr struct {
	Conn Conn
	Err  error
}

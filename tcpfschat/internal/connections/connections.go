package connections

import (
    "bytes"
	"github.com/google/uuid"
	"github.com/merisho/tcp-fs-chat/internal/chaterrors"
	"net"
	"sync"
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

func (conns *Connections) ForEach(f func(c Conn)) {
	conns.mu.Lock()
	defer conns.mu.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(conns.conns))
	for _, c := range conns.conns {
		go func(c Conn) {
			defer func() {
				wg.Done()
			}()

			f(c)
		}(c)
	}

	wg.Wait()
}

func (conns *Connections) Count() int {
	conns.mu.Lock()
	defer conns.mu.Unlock()
	return len(conns.conns)
}

func (conns *Connections) HandleConnectionErr(c Conn, err error) (connectionOk bool) {
	if chaterrors.IsTemporary(err) {
		return true
	}

	conns.RemoveByID(c.ID())
	return false
}

type ConnErr struct {
	Conn Conn
	Err  error
}

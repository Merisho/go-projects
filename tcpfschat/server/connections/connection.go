package connections

import (
    "net"
)

type Conn interface {
    net.Conn
    ID() []byte
}

func newConnection(conn net.Conn, id []byte) Conn {
    return &clientConn{
        Conn: conn,
        id: id,
    }
}

type clientConn struct {
    net.Conn
    id []byte
}

func (c *clientConn) ID() []byte {
    return c.id
}

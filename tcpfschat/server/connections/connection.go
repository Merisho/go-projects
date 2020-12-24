package connections

import (
    "net"
    "time"
)

type Conn interface {
    net.Conn
    ID() []byte
}

func newConnection(conn net.Conn, id []byte) Conn {
    return &clientConn{
        c: conn,
        id: id,
    }
}

type clientConn struct {
    c net.Conn
    id []byte
}

func (c *clientConn) ID() []byte {
    return c.id
}

func (c *clientConn) Read(b []byte) (n int, err error) {
    return c.c.Read(b)
}

func (c *clientConn) Write(b []byte) (n int, err error) {
    return c.c.Write(b)
}

func (c *clientConn) Close() error {
    return c.c.Close()
}

func (c *clientConn) LocalAddr() net.Addr {
    return c.c.LocalAddr()
}

func (c *clientConn) RemoteAddr() net.Addr {
    return c.c.RemoteAddr()
}

func (c *clientConn) SetDeadline(t time.Time) error {
    return c.c.SetDeadline(t)
}

func (c *clientConn) SetReadDeadline(t time.Time) error {
    return c.c.SetReadDeadline(t)
}

func (c *clientConn) SetWriteDeadline(t time.Time) error {
    return c.c.SetWriteDeadline(t)
}

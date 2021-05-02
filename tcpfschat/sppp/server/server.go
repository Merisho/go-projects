package server

import (
    "github.com/merisho/tcp-fs-chat/sppp/conn"
    "net"
    "strconv"
)

func Listen(host string, port uint16) (net.Listener, error) {
    address := host + ":" + strconv.FormatUint(uint64(port), 10)
    listener, err := net.Listen("tcp", address)
    if err != nil {
        return nil, err
    }

    return &Listener{
        listener: listener,
        addr: Addr{
            addr: address,
        },
    }, nil
}

type Listener struct {
    listener net.Listener
    addr Addr
}

func (l *Listener) Accept() (net.Conn, error) {
    c, err := l.listener.Accept()
    if err != nil {
        return nil, err
    }

    return conn.NewConn(c), nil
}

func (l *Listener) Close() error {
    return l.listener.Close()
}

func (l *Listener) Addr() net.Addr {
    return l.addr
}

type Addr struct {
    addr string
}

func (a Addr) Network() string {
    return "sppp"
}

func (a Addr) String() string {
    return a.addr
}

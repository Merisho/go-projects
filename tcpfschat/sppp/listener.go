package sppp

import (
    "net"
)

func NewSPPPListener(l net.Listener) *Listener {
    return &Listener{
        listener: l,
    }
}

type Listener struct {
    listener net.Listener
}

func (l *Listener) Accept() (*Conn, error) {
    c, err := l.listener.Accept()
    if err != nil {
        return nil, err
    }

    return NewConn(c), nil
}

func (l *Listener) Close() error {
    return l.listener.Close()
}

func (l *Listener) Addr() net.Addr {
    return l.listener.Addr()
}

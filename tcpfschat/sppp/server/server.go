package server

import (
    "github.com/merisho/tcp-fs-chat/sppp/conn"
    "net"
)

func NewSPPPServer(l net.Listener) *Server {
    return &Server{
        listener: l,
    }
}

type Server struct {
    listener net.Listener
}

func (l *Server) Accept() (*conn.Conn, error) {
    c, err := l.listener.Accept()
    if err != nil {
        return nil, err
    }

    return conn.NewConn(c), nil
}

func (l *Server) Close() error {
    return l.listener.Close()
}

func (l *Server) Addr() net.Addr {
    return l.listener.Addr()
}

package server

import (
    "github.com/merisho/tcp-fs-chat/server/connections"
    "log"
    "net"
)

func NewServer(ln net.Listener) Server {
    s := Server{
        ln: ln,
        conns: connections.Connections{},
    }

    s.start()

    return s
}

type Server struct {
    ln net.Listener
    conns connections.Connections
}

func (s *Server) start() *Server {
    ready := make(chan struct{})
    go func() {
        close(ready)
        for {
            c, err := s.ln.Accept()
            if err != nil {
                log.Println(err)
                continue
            }

            s.handleConnection(c)
        }
    }()

    <- ready

    return s
}

func (s *Server) handleConnection(c net.Conn) {
    s.conns.Add(c)

    go func() {
        b := make([]byte, 1024)
        for {
            n, err := c.Read(b)
            if err != nil {
                log.Println(err)
                continue
            }

            s.broadcast(c, b[:n])
        }
    }()
}

func (s *Server) broadcast(c net.Conn, msg []byte) {
    err := s.conns.BroadcastFrom(c, msg)
    if err != nil {
        log.Println(err)
    }
}

func (s *Server) ConnectionCount() int {
    return s.conns.Count()
}

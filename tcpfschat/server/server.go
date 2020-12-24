package server

import (
    "github.com/merisho/tcp-fs-chat/server/connections"
    "github.com/pkg/errors"
    "log"
    "net"
)

func NewServer(ln net.Listener) *Server {
    s := Server{
        ln: ln,
        conns: connections.Connections{},
    }

    s.start()

    return &s
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
    conn, err := s.conns.Add(c)
    if err != nil {
        log.Println(errors.Wrap(err, "could not add connection"))
        err = c.Close()
        if err != nil {
            log.Println(errors.Wrap(err, "could not close connection"))
        }

        return
    }

    go func() {
        _, err := conn.Write(conn.ID())
        if !s.conns.HandleConnectionErr(conn, err) {
            return
        }

        defer conn.Close()
        b := make([]byte, 1024 * 1024)
        for {
            n, err := conn.Read(b)
            if !s.conns.HandleConnectionErr(conn, err) {
                return
            }

            s.broadcast(conn, b[:n])
        }
    }()
}

func (s *Server) broadcast(c connections.Conn, msg []byte) {
    errs := s.conns.BroadcastFrom(c, msg)
    if len(errs) != 0 {
        for _, e := range errs {
            log.Println(e.Err.Error())
        }
    }
}

func (s *Server) ConnectionCount() int {
    return s.conns.Count()
}

func (s *Server) Disconnect(clientID []byte) error {
    conn := s.conns.RemoveByID(clientID)
    if conn == nil {
        return nil
    }

    return conn.Close()
}

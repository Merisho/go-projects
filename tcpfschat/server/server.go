package server

import (
    "github.com/merisho/tcp-fs-chat/server/connections"
    "github.com/pkg/errors"
    "log"
    "net"
    "time"
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

func (s *Server) broadcast(conn connections.Conn, msg []byte) {
    errChan := make(chan error, s.conns.Count())
    s.conns.ForEach(func(c connections.Conn) {
        if c == conn {
            return
        }

        _ = c.SetWriteDeadline(time.Now().Add(1 * time.Second))

        _, err := c.Write(msg)
        if err != nil {
            errChan <- err
        }
    })

    for len(errChan) > 0 {
        log.Println(<-errChan)
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

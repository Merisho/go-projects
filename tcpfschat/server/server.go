package server

import (
    "github.com/merisho/tcp-fs-chat/sppp"
    "io"
    "log"
    "sync"
)

func newServer(ln *sppp.Listener) *Server {
    s := &Server{
        ln: ln,
    }

    s.start()

    return s
}

type Server struct {
    ln *sppp.Listener
    conns []*sppp.Conn
    connsMx sync.Mutex
}

func (srv *Server) start() {
    go func() {
        for {
            cn, err := srv.ln.Accept()
            if err != nil {
                log.Println(err)
                return
            }

            go srv.handleConnection(cn)
        }
    }()
}

func (srv *Server) handleConnection(cn *sppp.Conn) {
    srv.connsMx.Lock()
    srv.conns = append(srv.conns, cn)
    srv.connsMx.Unlock()

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()

        metaMsg, err := cn.ReadMsg()
        if err != nil {
            log.Println(err)
            return
        }

        username := metaMsg.Content

        for {
           msg, err := cn.ReadMsg()
           if err != nil {
               if err == io.EOF {
                   return
               }

               log.Println(err)
               continue
           }

           srv.connsMx.Lock()
           for _, c := range srv.conns {
               if c == cn {
                   continue
               }

               m := append(username, ": "...)
               err := c.WriteMsg(append(m, msg.Content...))
               if err != nil {
                   log.Println(err)
               }
           }
           srv.connsMx.Unlock()
        }
    }()

    wg.Wait()

    srv.connsMx.Lock()
    defer srv.connsMx.Unlock()

    for i, c := range srv.conns {
        if c == cn {
            srv.conns = append(srv.conns[:i], srv.conns[i + 1:]...)
            break
        }
    }
}

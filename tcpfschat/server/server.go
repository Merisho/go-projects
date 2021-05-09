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

    wg.Add(1)
    go func() {
        defer wg.Done()

        for {
            rs := cn.ReadStream()
            go srv.handleStream(rs, cn)
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

func (srv *Server) handleStream(rs sppp.ReadStream, cn *sppp.Conn) {
    streamMeta, err := rs.ReadData()
    if err != nil {
        log.Printf("Could not read stream: %s", err)
        return
    }

    var wstreams []sppp.WriteStream

    srv.connsMx.Lock()
    for _, c := range srv.conns {
        if c == cn {
            continue
        }

        ws, err := c.WriteStream(streamMeta)
        if err != nil {
            log.Printf("Could not open write stream to connection: %s", err)
            continue
        }
        wstreams = append(wstreams, ws)
    }
    srv.connsMx.Unlock()

    for {
        chunk, err := rs.ReadData()
        if err != nil {
            if err != io.EOF {
                log.Printf("Could not read stream chunk: %s", err)
            }

            for _, ws := range wstreams {
                err := ws.Close()
                if err != nil {
                    log.Printf("Could not close write stream: %s", err)
                }
            }

            return
        }

        for _, ws := range wstreams {
            err := ws.WriteData(chunk)
            if err != nil {
                log.Printf("Could not write stream chunk: %s", err)
            }
        }
    }
}

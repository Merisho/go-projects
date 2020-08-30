package server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/merisho/tcp-fs-chat/server/connections"
	"log"
	"net"
	"strconv"
)

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", ":" + strconv.FormatUint(uint64(port), 10))
	if err != nil {
		return nil, err
	}

	s := &Server{
		ln: listener,
		Err: make(chan error),
		conns: &connections.Connections{},
	}
	go s.serve()

	return s, nil
}

type Server struct {
	ln net.Listener
	Err chan error
	conns *connections.Connections
}

func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) ConnectionCount() int {
	return s.conns.Count()
}

func (s *Server) serve() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			s.Err <- err
		} else {
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	s.conns.Add(conn)

	r := bufio.NewReader(conn)
	for {
		b := make([]byte, 1024)
		_, err := r.Read(b)
		if err != nil {
			s.conns.Remove(conn)
			err := conn.Close()
			if err != nil {
				log.Println(err)
			}

			s.Err <- err

			break
		}

		msg := s.formatMessage("", b)
		err = s.conns.BroadcastFrom(conn, msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) formatMessage(username string, b []byte) []byte {
	m := bytes.Replace(b, []byte{0}, []byte{}, -1)
	msg := fmt.Sprintf("[%s,%s]", username, string(m))

	return []byte(msg)
}

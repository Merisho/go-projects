package server

import (
	"bufio"
	"fmt"
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
		conns: &Connections{},
	}
	go s.serve()

	return s, nil
}

type Server struct {
	ln net.Listener
	Err chan error
	conns *Connections
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
			s.conns.Add(conn)
			fmt.Println("new conn")
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
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

		err = s.conns.Broadcast(b)
		if err != nil {
			log.Println(err)
		}
	}
}

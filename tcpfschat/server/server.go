package server

import (
	"bufio"
	"bytes"
	"errors"
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
	r := bufio.NewReader(conn)

	if s.authenticate(conn) != nil {
		return
	}

	s.conns.Add(conn)

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

func (s *Server) authenticate(conn net.Conn) error {
	b := make([]byte, 1024)
	_, err := conn.Read(b)
	if err != nil {
		return s.failAuth(conn)
	}

	username, password := s.authCreds(b)
	if username == "" || password == "" {
		return s.failAuth(conn)
	}

	_, err = conn.Write([]byte("auth success"))
	if err != nil {
		return s.failAuth(conn)
	}

	return nil
}

func (s *Server) authCreds(msg []byte) (username, password string) {
	prepared := bytes.Replace(msg, []byte{0}, []byte{}, -1)
	creds := bytes.Split(prepared, []byte("::"))

	if l := len(creds); l != 2 {
		return "", ""
	}

	return string(creds[0]), string(creds[1])
}

func (s *Server) failAuth(conn net.Conn) error {
	_, err := conn.Write([]byte("auth fail"))
	if err != nil {
		log.Println(err)
	}

	err = conn.Close()
	if err != nil {
		log.Println(err)
	}

	return errors.New("auth fail")
}

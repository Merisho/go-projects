package test

import (
    "github.com/merisho/tcp-fs-chat/sppp"
    "github.com/stretchr/testify/suite"
    "io"
    "net"
    "testing"
)

func TestSPPP(t *testing.T) {
    suite.Run(t, new(E2ESPPPTestSuite))
}

type E2ESPPPTestSuite struct {
    suite.Suite
}

func (s *E2ESPPPTestSuite) TestMessages() {
    tcp, err := net.Listen("tcp", ":7357")
    s.Require().NoError(err)

    srv := sppp.NewSPPPListener(tcp)

    c, err := net.Dial("tcp", ":7357")
    s.Require().NoError(err)

    client := sppp.NewConn(c)

    srvConn, err := srv.Accept()
    s.Require().NoError(err)

    err = client.WriteMsg([]byte("Hello World!"))
    s.Require().NoError(err)

    msg, err := srvConn.ReadMsg()
    s.Require().NoError(err)
    s.Require().Equal("Hello World!", string(msg.Content))

    err = srvConn.WriteMsg([]byte("Listener Message"))
    s.Require().NoError(err)

    msg, err = client.ReadMsg()
    s.Require().NoError(err)
    s.Require().Equal("Listener Message", string(msg.Content))
}

func (s *E2ESPPPTestSuite) TestStreams() {
    tcp, err := net.Listen("tcp", ":7358")
    s.Require().NoError(err)

    srv := sppp.NewSPPPListener(tcp)

    c, err := net.Dial("tcp", ":7358")
    s.Require().NoError(err)

    client := sppp.NewConn(c)

    ws, err := client.WriteStream([]byte("stream meta"))
    s.Require().NoError(err)

    err = ws.WriteData([]byte("chunk 1"))
    s.Require().NoError(err)

    err = ws.WriteData([]byte("chunk 2"))
    s.Require().NoError(err)

    err = ws.Close()
    s.Require().NoError(err)

    srvConn, err := srv.Accept()
    s.Require().NoError(err)

    rs, err := srvConn.ReadStream()
    s.Require().NoError(err)

    meta, err := rs.ReadData()
    s.Require().NoError(err)
    s.Require().Equal("stream meta", string(meta))

    chunk, err := rs.ReadData()
    s.Require().NoError(err)
    s.Require().Equal("chunk 1", string(chunk))

    chunk, err = rs.ReadData()
    s.Require().NoError(err)
    s.Require().Equal("chunk 2", string(chunk))

    chunk, err = rs.ReadData()
    s.Require().Equal(io.EOF, err)
    s.Require().Nil(chunk)
}

//func (s *E2ESPPPTestSuite) TestConnectionClose() {
//    tcp, err := net.Listen("tcp", ":7359")
//    s.Require().NoError(err)
//
//    srv := listener.NewSPPPListener(tcp)
//
//    c, err := net.Dial("tcp", ":7359")
//    s.Require().NoError(err)
//
//    client := conn.NewConn(c)
//
//    srvConn, err := srv.Accept()
//    s.Require().NoError(err)
//
//    client.Close()
//}

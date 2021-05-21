package server

import (
    "github.com/merisho/tcp-fs-chat/sppp"
    "github.com/stretchr/testify/suite"
    "io"
    "net"
    "testing"
    "time"
)

func TestChatServer(t *testing.T) {
    suite.Run(t, new(ChatServerTestSuite))
}

type ChatServerTestSuite struct {
    suite.Suite
}

func (s *ChatServerTestSuite) SetupSuite() {
   ln, err := net.Listen("tcp", ":1337")
   s.Require().NoError(err)

   sln := sppp.NewSPPPListener(ln)
   newServer(sln)
}

func (s *ChatServerTestSuite) TestMessagesE2E() {
    require := s.Require()

    conn1, err := net.Dial("tcp", ":1337")
    require.NoError(err)
    client1 := sppp.NewConn(conn1)
    require.NoError(client1.WriteMsg([]byte("user1")))

    // Wait a bit to ensure client1 is connected
    time.Sleep(5 * time.Millisecond)

    conn2, err := net.Dial("tcp", ":1337")
    require.NoError(err)
    client2 := sppp.NewConn(conn2)
    require.NoError(client2.WriteMsg([]byte("user2")))

    err = client2.WriteMsg([]byte("client 2"))
    require.NoError(err)

    msg, err := client1.ReadMsg()
    require.NoError(err)
    require.Equal("user2: client 2", string(msg.Content))

    require.Equal(0, client1.MsgCount())

    err = client1.Close()
    require.NoError(err)

    err = client2.Close()
    require.NoError(err)
}

func (s *ChatServerTestSuite) TestStreamsE2E() {
   require := s.Require()

   ln, err := net.Listen("tcp", ":1338")
   require.NoError(err)

   sln := sppp.NewSPPPListener(ln)
   newServer(sln)

   conn1, err := net.Dial("tcp", ":1338")
   require.NoError(err)
   client1 := sppp.NewConn(conn1)
   require.NoError(client1.WriteMsg([]byte("user1")))

   conn2, err := net.Dial("tcp", ":1338")
   require.NoError(err)
   client2 := sppp.NewConn(conn2)
   require.NoError(client2.WriteMsg([]byte("user2")))

   ws, err := client1.WriteStream([]byte("client 1 stream"))
   require.NoError(err)

   require.NoError(ws.WriteData([]byte("client 1 chunk 1")))
   require.NoError(ws.WriteData([]byte("client 1 chunk 2")))
   require.NoError(ws.WriteData([]byte("client 1 chunk 3")))
   require.NoError(ws.Close())

   rs, err := client2.ReadStream()
   require.NoError(err)

   meta, err := rs.ReadData()
   require.NoError(err)
   require.Equal("client 1 stream", string(meta))

   chunk, err := rs.ReadData()
   require.NoError(err)
   require.Equal("client 1 chunk 1", string(chunk))

   chunk, err = rs.ReadData()
   require.NoError(err)
   require.Equal("client 1 chunk 2", string(chunk))

   chunk, err = rs.ReadData()
   require.NoError(err)
   require.Equal("client 1 chunk 3", string(chunk))

   chunk, err = rs.ReadData()
   require.Equal(io.EOF, err)
   require.Nil(chunk)
}

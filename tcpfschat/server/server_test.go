package server

import (
    "github.com/merisho/tcp-fs-chat/sppp"
    "github.com/stretchr/testify/suite"
    "net"
    "testing"
)

func TestChatServer(t *testing.T) {
    suite.Run(t, new(ChatServerTestSuite))
}

type ChatServerTestSuite struct {
    suite.Suite
}

func (s *ChatServerTestSuite) TestMessagesE2E() {
    require := s.Require()
    
    ln, err := net.Listen("tcp", ":1337")
    require.NoError(err)

    sln := sppp.NewSPPPListener(ln)
    newServer(sln)

    conn1, err := net.Dial("tcp", ":1337")
    require.NoError(err)
    client1 := sppp.NewConn(conn1)
    require.NoError(client1.WriteMsg([]byte("user1")))

    conn2, err := net.Dial("tcp", ":1337")
    require.NoError(err)
    client2 := sppp.NewConn(conn2)
    require.NoError(client2.WriteMsg([]byte("user2")))

    err = client1.WriteMsg([]byte("client 1"))
    require.NoError(err)

    msg, err := client2.ReadMsg()
    require.NoError(err)
    require.Equal("user1: client 1", string(msg.Content))

    require.Equal(0, client1.MsgCount())
}

package server

import (
    "github.com/merisho/tcp-fs-chat/sppp"
    "github.com/stretchr/testify/suite"
    "math/rand"
    "net"
    "testing"
    "time"
)

func TestConn(t *testing.T) {
    suite.Run(t, new(ConnTestSuite))
}

type ConnTestSuite struct {
    suite.Suite
}

func (s *ConnTestSuite) TestMsgRead() {
    c1, c2 := net.Pipe()
    reader := NewConn(c2)

    id := rand.Int63()

    rawMsg := sppp.NewMessage(id, sppp.TextType, []byte("test")).Marshal()
    _, _ = c1.Write(rawMsg[:])

    rawMsg = sppp.NewMessage(id, sppp.TextType, []byte(" message")).Marshal()
    _, _ = c1.Write(rawMsg[:])

    rawMsg = sppp.NewMessage(id, sppp.MsgEndType, nil).Marshal()
    _, _ = c1.Write(rawMsg[:])

    msg, err := reader.ReadMsg()
    s.Require().NoError(err)
    s.Require().Equal("test message", string(msg.Content))
}

func (s *ConnTestSuite) TestMsgReadTimeout() {
   c1, c2 := net.Pipe()
   reader := NewConn(c2)
   reader.SetMessageReadTimeout(500 * time.Millisecond)

   id := rand.Int63()
   rawMsg := sppp.NewMessage(id, sppp.TextType, []byte("test")).Marshal()
   _, _ = c1.Write(rawMsg[:])

   msgChan := make(chan sppp.Message)
   go func() {
       m, _ := reader.ReadMsg()
       msgChan <- m
   }()

   select {
   case <- msgChan:
       s.Fail("Must not receive a message")
   case <- time.After(500 * time.Millisecond):
   }

   var rawTimeoutResponse [1024]byte
   _, _ = c1.Read(rawTimeoutResponse[:])

   timeoutRes, err := sppp.UnmarshalMessage(rawTimeoutResponse)
   s.Require().NoError(err)
   s.Require().EqualValues(sppp.TimeoutType, timeoutRes.Type)
   s.Require().Equal(id, timeoutRes.ID)
}

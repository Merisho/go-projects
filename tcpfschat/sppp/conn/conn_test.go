package conn

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
    rand *rand.Rand
}

func (s *ConnTestSuite) SetupSuite() {
    s.rand = rand.New(rand.NewSource(time.Now().Unix()))
}

func (s *ConnTestSuite) TestMsgRead() {
    c1, c2 := net.Pipe()
    reader := NewConn(c2)

    id := s.rand.Int63()

    rawMsg := sppp.NewMessage(id, sppp.TextType, []byte("test")).Marshal()
    _, _ = c1.Write(rawMsg[:])

    rawMsg = sppp.NewMessage(id, sppp.TextType, []byte(" message")).Marshal()
    _, _ = c1.Write(rawMsg[:])

    rawMsg = sppp.NewMessage(id, sppp.EndType, nil).Marshal()
    _, _ = c1.Write(rawMsg[:])

    msg, err := reader.ReadMsg()
    s.Require().NoError(err)
    s.Require().Equal("test message", string(msg.Content))
}

func (s *ConnTestSuite) TestMsgReadTimeout() {
    c1, c2 := net.Pipe()
    reader := NewConn(c2)
    reader.SetMessageReadTimeout(50 * time.Millisecond)

    id := s.rand.Int63()
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
    case <- time.After(60 * time.Millisecond):
    }

    var rawTimeoutResponse [1024]byte
    _, _ = c1.Read(rawTimeoutResponse[:])

    timeoutRes, err := sppp.UnmarshalMessage(rawTimeoutResponse)
    s.Require().NoError(err)
    s.Require().EqualValues(sppp.TimeoutType, timeoutRes.Type)
    s.Require().Equal(id, timeoutRes.ID)
}

func (s *ConnTestSuite) TestReadStream() {
    c1, c2 := net.Pipe()
    reader := NewConn(c2)

    id := s.rand.Int63()
    streamMeta := []byte("stream meta info")
    rawMsg := sppp.NewMessage(id, sppp.StreamType, streamMeta).Marshal()
    _, _  = c1.Write(rawMsg[:])

    streamData := []byte("chunk 1")
    rawMsg = sppp.NewMessage(id, sppp.StreamType, streamData).Marshal()
    _, _  = c1.Write(rawMsg[:])

    streamData = []byte("chunk 2")
    rawMsg = sppp.NewMessage(id, sppp.StreamType, streamData).Marshal()
    _, _  = c1.Write(rawMsg[:])

    stream, _ := reader.ReadStream()

    meta := <- stream
    s.Require().Equal("stream meta info", string(meta))

    chunk := <- stream
    s.Require().Equal("chunk 1", string(chunk))

    chunk = <- stream
    s.Require().Equal("chunk 2", string(chunk))

    rawMsg = sppp.NewMessage(id, sppp.EndType, nil).Marshal()
    _, _  = c1.Write(rawMsg[:])

    chunk, ok := <- stream
    s.Require().False(ok)
    s.Require().Nil(chunk)
}

func (s *ConnTestSuite) TestReadStreamTimeout() {
    c1, c2 := net.Pipe()
    reader := NewConn(c2)
    reader.SetStreamReadTimeout(50 * time.Millisecond)

    test := func() {
        id := s.rand.Int63()
        streamMeta := []byte("stream meta info")
        rawMsg := sppp.NewMessage(id, sppp.StreamType, streamMeta).Marshal()
        _, _  = c1.Write(rawMsg[:])

        stream, errs := reader.ReadStream()
        <- stream

        select {
        case err := <- errs:
            s.Require().EqualError(err, TimeoutError.Error())
        case <- time.After(70 * time.Millisecond):
            s.Fail("timeout must occur")
        }

        var timeoutRes [1024]byte
        _, _ = c1.Read(timeoutRes[:])
        timeoutMsg, err := sppp.UnmarshalMessage(timeoutRes)
        s.Require().NoError(err)
        s.Require().EqualValues(sppp.TimeoutType, timeoutMsg.Type)
    }

    test()
    test()
}

func (s *ConnTestSuite) TestHandleInvalidMessage() {
    c1, c2 := net.Pipe()
    reader := NewConn(c2)
    reader.SetStreamReadTimeout(50 * time.Millisecond)

    rawMsg := []byte("a garbage message")
    _, _ = c1.Write(rawMsg[:])

    msgChan := make(chan sppp.Message)
    go func() {
        m, _ := reader.ReadMsg()
        msgChan <- m
    }()

    select {
    case <- msgChan:
        s.Fail("Must not receive a garbage message")
    case <- time.After(60 * time.Millisecond):
    }

    var rawInvalidMsgResponse [1024]byte
    _, _ = c1.Read(rawInvalidMsgResponse[:])

    timeoutRes, err := sppp.UnmarshalMessage(rawInvalidMsgResponse)
    s.Require().NoError(err)
    s.Require().EqualValues(sppp.ErrorType, timeoutRes.Type)
}
package sppp

import (
    "bytes"
    "github.com/stretchr/testify/suite"
    "io"
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

    rawMsg := NewMessage(id, TextType, []byte("test")).Marshal()
    _, _ = c1.Write(rawMsg[:])

    rawMsg = NewMessage(id, TextType, []byte(" message")).Marshal()
    _, _ = c1.Write(rawMsg[:])

    rawMsg = NewMessage(id, EndType, nil).Marshal()
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
    rawMsg := NewMessage(id, TextType, []byte("test")).Marshal()
    _, _ = c1.Write(rawMsg[:])

    msgChan := make(chan Message)
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

    timeoutRes, err := UnmarshalMessage(rawTimeoutResponse)
    s.Require().NoError(err)
    s.Require().EqualValues(TimeoutType, timeoutRes.Type)
    s.Require().Equal(id, timeoutRes.ID)
}

func (s *ConnTestSuite) TestReadStream() {
    c1, c2 := net.Pipe()
    reader := NewConn(c2)

    id := s.rand.Int63()
    streamMeta := []byte("stream meta info")
    rawMsg := NewMessage(id, StreamType, streamMeta).Marshal()
    _, _  = c1.Write(rawMsg[:])

    streamData := []byte("chunk 1")
    rawMsg = NewMessage(id, StreamType, streamData).Marshal()
    _, _  = c1.Write(rawMsg[:])

    streamData = []byte("chunk 2")
    rawMsg = NewMessage(id, StreamType, streamData).Marshal()
    _, _  = c1.Write(rawMsg[:])

    stream, err := reader.ReadStream()
    s.Require().NoError(err)

    meta, err := stream.ReadData()
    s.Require().NoError(err)
    s.Require().Equal("stream meta info", string(meta))

    chunk, err := stream.ReadData()
    s.Require().NoError(err)
    s.Require().Equal("chunk 1", string(chunk))

    chunk, err = stream.ReadData()
    s.Require().NoError(err)
    s.Require().Equal("chunk 2", string(chunk))

    rawMsg = NewMessage(id, EndType, nil).Marshal()
    _, _  = c1.Write(rawMsg[:])

    chunk, err = stream.ReadData()
    s.Require().Equal(io.EOF, err)
    s.Require().Nil(chunk)
}

func (s *ConnTestSuite) TestReadStreamTimeout() {
   c1, c2 := net.Pipe()
   reader := NewConn(c2)
   reader.SetStreamReadTimeout(50 * time.Millisecond)

   test := func() {
       id := s.rand.Int63()
       streamMeta := []byte("stream meta info")
       rawMsg := NewMessage(id, StreamType, streamMeta).Marshal()
       _, _  = c1.Write(rawMsg[:])

       stream, err := reader.ReadStream()
       s.Require().NoError(err)

       _, err = stream.ReadData()
        s.Require().NoError(err)

       _, err = stream.ReadData()
       s.Require().Equal(TimeoutError, err)

       var timeoutRes [1024]byte
       _, _ = c1.Read(timeoutRes[:])
       timeoutMsg, err := UnmarshalMessage(timeoutRes)
       s.Require().NoError(err)
       s.Require().EqualValues(TimeoutType, timeoutMsg.Type)
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

    msgChan := make(chan Message)
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

    timeoutRes, err := UnmarshalMessage(rawInvalidMsgResponse)
    s.Require().NoError(err)
    s.Require().EqualValues(ErrorType, timeoutRes.Type)
}

func (s *ConnTestSuite) TestWriteMessage() {
    c1, c2 := net.Pipe()
    writer := NewConn(c1)
    reader := NewConn(c2)

    rawMsg := bytes.Repeat([]byte("test"), 1024)
    err := writer.WriteMsg(rawMsg)
    s.Require().NoError(err)

    msg, err := reader.ReadMsg()
    s.Require().NoError(err)
    s.Require().Equal(rawMsg, msg.Content)
}

func (s *ConnTestSuite) TestWriteStream() {
   c1, c2 := net.Pipe()
   writer := NewConn(c1)
   reader := NewConn(c2)

   metaInfo := []byte("stream meta info")
   ws, err := writer.WriteStream(metaInfo)
   s.Require().NoError(err)

   err = ws.WriteData([]byte("chunk 1"))
   s.Require().NoError(err)

   err = ws.WriteData([]byte("chunk 2"))
   s.Require().NoError(err)

   err = ws.Close()
   s.Require().NoError(err)

   rs, err := reader.ReadStream()
   s.Require().NoError(err)

   chunk, err := rs.ReadData()
   s.Require().NoError(err)
   s.Require().Equal("stream meta info", string(chunk))

   chunk, err = rs.ReadData()
   s.Require().NoError(err)
   s.Require().Equal("chunk 1", string(chunk))

   chunk, err = rs.ReadData()
   s.Require().NoError(err)
   s.Require().Equal("chunk 2", string(chunk))

   chunk, err = rs.ReadData()
   s.Require().Equal(io.EOF, err)
   s.Require().Nil(chunk)
}

func (s *ConnTestSuite) TestConnectionClosing() {
    c1, c2 := net.Pipe()
    writer := NewConn(c1)
    reader := NewConn(c2)

    reader.Close()

    _, err := writer.ReadMsg()
    s.Equal(io.EOF, err)

    _, err = writer.ReadStream()
    s.Equal(io.EOF, err)
}

func (s *ConnTestSuite) TestStreamsClosing_ReaderCloses() {
    c1, c2 := net.Pipe()
    writer := NewConn(c1)
    reader := NewConn(c2)

    ws, err := writer.WriteStream([]byte("test"))
    s.Require().NoError(err)

    reader.Close()

    err = ws.WriteData([]byte("test"))
    s.Require().Equal(io.ErrClosedPipe, err)
}

func (s *ConnTestSuite) TestStreamsClosing_WriterCloses() {
    c1, c2 := net.Pipe()
    writer := NewConn(c1)
    reader := NewConn(c2)

    _, err := writer.WriteStream([]byte("test"))
    s.Require().NoError(err)

    writer.Close()

    rs, err := reader.ReadStream()
    s.Require().NoError(err)
    s.Require().NotNil(rs)

    // First chunk (which is stream meta) should be OK
    _, err = rs.ReadData()
    s.Require().NoError(err)

    _, err = rs.ReadData()
    s.Require().Equal(io.EOF, err)
}

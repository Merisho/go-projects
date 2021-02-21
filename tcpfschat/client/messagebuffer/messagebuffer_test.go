package messagebuffer

import (
    "github.com/stretchr/testify/suite"
    "testing"
)

func TestMessageBuffer(t *testing.T) {
    suite.Run(t, new(MessageBufferTestSuite))
}

type MessageBufferTestSuite struct {
    suite.Suite
    buffer MessageBuffer
}

func (s *MessageBufferTestSuite) BeforeTest() {
    s.buffer = MessageBuffer{}
}

func (s *MessageBufferTestSuite) TestEmptyRawMessage() {
    msgs := s.buffer.Messages([]byte{})
    s.Empty(msgs)
}

func (s *MessageBufferTestSuite) TestRawMessageWithoutLeftovers() {
    testMsg := []byte("test")
    rawMsg := append(testMsg, 0)
    rawMsg = append(rawMsg, testMsg...)
    rawMsg = append(rawMsg, 0)

    msgs := s.buffer.Messages(rawMsg)
    s.Equal([][]byte{testMsg, testMsg}, msgs)
}

func (s *MessageBufferTestSuite) TestRawMessageWithLeftovers() {
    testMsg := []byte("test")
    rawMsg := append(testMsg, 0)
    rawMsg = append(rawMsg, []byte("te")...)

    msgs := s.buffer.Messages(rawMsg)
    s.Equal([][]byte{testMsg}, msgs)

    rawMsg = append([]byte("st"), 0)
    rawMsg = append(rawMsg, testMsg...)
    rawMsg = append(rawMsg, 0)
    msgs = s.buffer.Messages(rawMsg)

    s.Equal([][]byte{testMsg, testMsg}, msgs)
}

func (s *MessageBufferTestSuite) TestMultipleMessagesWithLeftovers() {
    msgs := s.buffer.Messages([]byte("t"))
    s.Empty(msgs)

    msgs = s.buffer.Messages([]byte("e"))
    s.Empty(msgs)

    msgs = s.buffer.Messages([]byte("s"))
    s.Empty(msgs)

    msgs = s.buffer.Messages([]byte("t\x00"))
    s.Equal([][]byte{[]byte("test")}, msgs)
}

func (s *MessageBufferTestSuite) TestEndPreviousMessageInSucceedingRawMessage() {
    msgs := s.buffer.Messages([]byte("test"))
    s.Empty(msgs)

    msgs = s.buffer.Messages([]byte{0})
    s.Equal([][]byte{[]byte("test")}, msgs)
}

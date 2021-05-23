package sppp

import (
    "github.com/stretchr/testify/suite"
    "strings"
    "testing"
)

func TestSPPPMessages(t *testing.T) {
    suite.Run(t, new(SPPPTestSuite))
}

type SPPPTestSuite struct {
    suite.Suite
}

func (s *SPPPTestSuite) TestUint64ToBytes() {
    cases := []struct{
        in uint64
        expected [8]byte
    }{
        {
            in: 0,
            expected: [8]byte{},
        },
        {
            in: 10,
            expected: [8]byte{10},
        },
        {
            in: 255,
            expected: [8]byte{255},
        },
        {
            in: 256,
            expected: [8]byte{0, 1},
        },
        {
            in: 511,
            expected: [8]byte{255, 1},
        },
        {
            in: 512,
            expected: [8]byte{0, 2},
        },
        {
            in: 1<<64 - 1,
            expected: [8]byte{255, 255, 255, 255, 255, 255, 255, 255},
        },
    }

    for _, c := range cases {
        b := Uint64ToBytes(c.in)
        s.Equal(c.expected, b)

        n := BytesToUint64(b)
        s.Equal(c.in, n)
    }
}

func (s *SPPPTestSuite) TestUnmarshal() {
   headerSize := 17
   var rawMsg [1024]byte
   size := uint64(100)
   expectedMsg := strings.Repeat("a", int(size))

   rawMsg[0] = TextType

   rawSize := Uint64ToBytes(size)
   copy(rawMsg[1:], rawSize[:])

   msgID := Uint64ToBytes(1)
   copy(rawMsg[9:], msgID[:])
   copy(rawMsg[headerSize:], expectedMsg)

    msg, err := UnmarshalMessage(rawMsg)
    s.NoError(err)
    s.EqualValues(TextType, msg.Type)
    s.EqualValues(size, msg.Size)
    s.EqualValues(1, msg.ID)
    s.EqualValues(expectedMsg, string(msg.Content))
}

func (s *SPPPTestSuite) TestUnmarshalBadMessage() {
    var badMsg [1024]byte
    copy(badMsg[:], "a garbage message")

    msg, err := UnmarshalMessage(badMsg)

    s.Require().EqualError(err, "bad message")
    s.Require().True(msg.Empty())
}

func (s *SPPPTestSuite) TestMarshal() {
    headerSize := 17
    var rawMsg [1024]byte
    actualMsgSize := len(rawMsg) - headerSize
    actualMsg := strings.Repeat("a", actualMsgSize)

    rawMsg[0] = TextType

    size := Uint64ToBytes(uint64(actualMsgSize))
    copy(rawMsg[1:], size[:])

    msgID := Uint64ToBytes(1)
    copy(rawMsg[9:], msgID[:])
    copy(rawMsg[headerSize:], actualMsg)

    msg := Message{
        Type:    TextType,
        Size:    uint64(actualMsgSize),
        ID:      1,
        Content: []byte(actualMsg),
    }

    b := msg.Marshal()
    s.EqualValues(rawMsg, b)
}

func (s *SPPPTestSuite) TestSplitIntoMessages() {
    msgContentChunks := 5
    msg := []byte(strings.Repeat("a", msgContentSize * msgContentChunks))

    msgs := SplitIntoMessages(0, TextType, msg)

    s.Require().Len(msgs, msgContentChunks)
}

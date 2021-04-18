package sppp

import (
    "github.com/stretchr/testify/suite"
    "strings"
    "testing"
)

func TestSPPP(t *testing.T) {
    suite.Run(t, new(SPPPTestSuite))
}

type SPPPTestSuite struct {
    suite.Suite
}

func (s *SPPPTestSuite) TestInt64ToBytes() {
    cases := []struct{
        in int64
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
    }

    for _, c := range cases {
        b := Int64ToBytes(c.in)
        s.Equal(c.expected, b)

        n := BytesToInt64(b)
        s.Equal(c.in, n)
    }
}

func (s *SPPPTestSuite) TestUnmarshal() {
   headerSize := 17
   var rawMsg [1024]byte
   actualMsg := strings.Repeat("a", len(rawMsg) - headerSize)

   rawMsg[0] = 0

   size := Int64ToBytes(1000)
   copy(rawMsg[1:], size[:])

   msgID := Int64ToBytes(1)
   copy(rawMsg[9:], msgID[:])
   copy(rawMsg[headerSize:], actualMsg)

    msg, err := UnmarshalMessage(rawMsg)
    s.NoError(err)
    s.EqualValues(TextType, msg.Type)
    s.EqualValues(1000, msg.Size)
    s.EqualValues(1, msg.ID)
    s.EqualValues(actualMsg, string(msg.Content))
}

func (s *SPPPTestSuite) TestMarshal() {
    headerSize := 17
    var rawMsg [1024]byte
    actualMsgSize := len(rawMsg) - headerSize
    actualMsg := strings.Repeat("a", actualMsgSize)

    rawMsg[0] = TextType

    size := Int64ToBytes(int64(actualMsgSize))
    copy(rawMsg[1:], size[:])

    msgID := Int64ToBytes(1)
    copy(rawMsg[9:], msgID[:])
    copy(rawMsg[headerSize:], actualMsg)

    msg := Message{
        Type:    TextType,
        Size:    int64(actualMsgSize),
        ID:      1,
        Content: []byte(actualMsg),
    }

    b, err := msg.Marshal()
    s.NoError(err)
    s.EqualValues(rawMsg, b)
}

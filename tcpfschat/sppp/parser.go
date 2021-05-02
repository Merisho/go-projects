package sppp

import "errors"

// The message is 1024 bytes
// 1 byte - type
// 8 bytes - size
// 8 bytes - ID
// 1007 bytes - content

const (
    headerSize = 17
    totalMsgSize = 1024
)

var (
    badMsgError = errors.New("bad message")
)

func NewMessage(id int64, t MessageType, content []byte) Message {
    return Message{
        Type:    t,
        Size:    int64(len(content)),
        ID:      id,
        Content: content,
    }
}

type Message struct {
    Type MessageType
    Size int64
    ID int64
    Content []byte
}

func UnmarshalMessage(msg [1024]byte) (m Message, err error) {
    header := msg[:headerSize]
    msgType := header[0]

    var rawSize [8]byte
    copy(rawSize[:], header[1:9])

    var rawMsgID [8]byte
    copy(rawMsgID[:], header[9:headerSize])

    size := BytesToInt64(rawSize)

    if invalidMessage(msgType, size) {
        return Message{}, badMsgError
    }

    m = Message{
        Type:    MessageType(header[0]),
        Size:    size,
        ID:      BytesToInt64(rawMsgID),
        Content: msg[headerSize:headerSize + size],
    }

    return m, err
}

func invalidMessage(msgType byte, size int64) bool {
    return size > totalMsgSize - headerSize || msgType >= maxMessageTypeIota
}

func (m Message) Marshal() [1024]byte {
    var b [1024]byte
    b[0] = byte(m.Type)

    size := Int64ToBytes(m.Size)
    copy(b[1:], size[:])

    id := Int64ToBytes(m.ID)
    copy(b[9:], id[:])
    copy(b[headerSize:], m.Content)

    return b
}

func (m Message) Empty() bool {
    return m.Type == 0 && m.Size == 0 && m.ID == 0
}

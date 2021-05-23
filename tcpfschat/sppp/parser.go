package sppp

import (
    "errors"
)

// The message is 8192 bytes
// 1 byte - type
// 8 bytes - size
// 8 bytes - ID
// 8175 bytes - content

const (
    headerSize = 17
    totalMsgSize = 8192
    msgContentSize = totalMsgSize - headerSize
)

var (
    badMsgError = errors.New("bad message")
)

func NewMessage(id uint64, t MessageType, content []byte) Message {
    return Message{
        Type:    t,
        Size:    uint64(len(content)),
        ID:      id,
        Content: content,
    }
}

type Message struct {
    Type MessageType
    Size uint64
    ID uint64
    Content []byte
}

func UnmarshalMessage(msg [totalMsgSize]byte) (m Message, err error) {
    header := msg[:headerSize]
    msgType := header[0]

    var rawSize [8]byte
    copy(rawSize[:], header[1:9])

    var rawMsgID [8]byte
    copy(rawMsgID[:], header[9:headerSize])

    size := BytesToUint64(rawSize)

    if invalidMessage(msgType, size) {
        return Message{}, badMsgError
    }

    m = Message{
        Type:    MessageType(header[0]),
        Size:    size,
        ID:      BytesToUint64(rawMsgID),
        Content: msg[headerSize:headerSize + size],
    }

    return m, err
}

func invalidMessage(msgType byte, size uint64) bool {
    return size > totalMsgSize - headerSize || msgType >= maxMessageTypeIota
}

func (m Message) Marshal() [totalMsgSize]byte {
    var b [totalMsgSize]byte
    b[0] = byte(m.Type)

    size := Uint64ToBytes(m.Size)
    copy(b[1:], size[:])

    id := Uint64ToBytes(m.ID)
    copy(b[9:], id[:])
    copy(b[headerSize:], m.Content)

    return b
}

func (m Message) Empty() bool {
    return m.Type == 0 && m.Size == 0 && m.ID == 0
}

func SplitIntoMessages(id uint64, t MessageType, msg []byte) []Message {
    var msgs []Message

    ln := len(msg)
    for i := 0; i < ln; i += msgContentSize {
        end := i + msgContentSize
        if end > ln {
            end = ln
        }

        m := Message{
            Type:    t,
            Size:    uint64(len(msg[i:end])),
            ID:      id,
            Content: msg[i:end],
        }
        msgs = append(msgs, m)
    }

    return msgs
}

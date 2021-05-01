package sppp

// The message is 1024 bytes
// 1 byte - type
// 8 bytes - size
// 8 bytes - ID
// 1007 bytes - content

const (
    headerSize = 17
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

func UnmarshalMessage(msg [1024]byte) (Message, error) {
    header := msg[:headerSize]

    var rawSize [8]byte
    copy(rawSize[:], header[1:9])

    var rawMsgID [8]byte
    copy(rawMsgID[:], header[9:headerSize])

    size := BytesToInt64(rawSize)

    return Message{
        Type:    MessageType(header[0]),
        Size:    size,
        ID:      BytesToInt64(rawMsgID),
        Content: msg[headerSize:headerSize + size],
    }, nil
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

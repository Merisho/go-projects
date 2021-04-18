package sppp

// The message is 1024 bytes
// 1 byte - type
// 8 bytes - size
// 8 bytes - ID
// 1007 bytes - content

const (
    headerSize = 17
)

type Message struct {
    Type MessageType
    Size int64
    ID int64
    Content []byte
}

func UnmarshalMessage(msg [1024]byte) (Message, error) {
    header := msg[:headerSize]
    msgType := header[0]
    rawSize := header[1:9]
    rawMsgID := header[9:headerSize]

    var sz [8]byte
    copy(sz[:], rawSize)

    var msgID [8]byte
    copy(msgID[:], rawMsgID)

    return Message{
        Type:    MessageType(msgType),
        Size:    BytesToInt64(sz),
        ID:      BytesToInt64(msgID),
        Content: msg[headerSize:],
    }, nil
}

func (m Message) Marshal() ([1024]byte, error) {
    var b [1024]byte
    b[0] = byte(m.Type)

    size := Int64ToBytes(m.Size)
    copy(b[1:], size[:])

    id := Int64ToBytes(m.ID)
    copy(b[9:], id[:])
    copy(b[headerSize:], m.Content)

    return b, nil
}

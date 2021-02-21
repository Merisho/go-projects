package messagebuffer

import "bytes"

type MessageBuffer struct {
    buff []byte
}

func (b *MessageBuffer) Messages(rawMsg []byte) [][]byte {
    if len(rawMsg) == 0 {
        return nil
    }

    msgs := bytes.Split(rawMsg, []byte{0})

    if len(msgs) > 0 {
        msgs[0] = append(b.buff, msgs[0]...)
    }

    b.buff = msgs[len(msgs) - 1]

    return msgs[:len(msgs) - 1]
}

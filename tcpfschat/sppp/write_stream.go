package sppp

import (
    "io"
)

type WriteStream interface {
    io.WriteCloser
    WriteData([]byte) error
}

func newWriteStream(msgID int64, w io.Writer) *writeStream {
    s := &writeStream{
        msgID: msgID,
        out:   w,
    }

    return s
}

type writeStream struct {
    msgID      int64
    out io.Writer
}

func (s *writeStream) Close() error {
    msg := NewMessage(s.msgID, EndType, nil)
    return s.write(msg)
}

func (s *writeStream) WriteData(b []byte) error {
    _, err := s.Write(b)
    return err
}

func (s *writeStream) Write(b []byte) (int, error) {
    msgs := SplitIntoMessages(s.msgID, StreamType, b)

    for _, m := range msgs {
        err := s.write(m)
        if err != nil {
            return 0, err
        }
    }

    return len(b), nil
}

func (s *writeStream) write(m Message) error {
   raw := m.Marshal()
   _, err := s.out.Write(raw[:])
   return err
}

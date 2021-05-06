package conn

import (
    "fmt"
    "github.com/merisho/tcp-fs-chat/sppp"
    "io"
    "time"
)

type ReadStream interface {
    io.ReadCloser
    ReadData() ([]byte, error)
}

type WriteStream interface {
    io.WriteCloser
    WriteData([]byte) error
}

func NewStream(msgID int64, readTimeout time.Duration) *Stream {
    s := &Stream{
        msgID:           msgID,
        readChunks:      make(chan sppp.Message, 1024),
        readErrs:        make(chan error),
        readSig:         make(chan struct{}),
        timeoutOccurred: false,
        readTimeout:     readTimeout,
        readTimeoutSig:  make(chan struct{}),

        write:           func(sppp.Message) error { return nil },
    }

    s.acceptReadSignals()

    return s
}

type Stream struct {
    msgID      int64

    readChunks chan sppp.Message
    readErrs   chan error
    readSig    chan struct{}
    readTimeoutSig    chan struct{}
    timeoutOccurred bool
    readTimeout time.Duration

    write func(sppp.Message) error
}

func (s *Stream) ReadData() ([]byte, error) {
    b := make([]byte, 2048)
    n, err := s.Read(b)
    if err != nil {
        return nil, err
    }

    return b[:n], nil
}

func (s *Stream) Read(b []byte) (int, error) {
    select {
    case m, ok := <- s.readChunks:
        if !ok {
            close(s.readErrs)
            return 0, io.EOF
        }

        copy(b, m.Content)
        return int(m.Size), nil
    case err, ok := <- s.readErrs:
        if !ok {
            return 0, io.EOF
        }

        return 0, err
    }
}

func (s *Stream) ReadTimeoutWait() chan struct{} {
    return s.readTimeoutSig
}

func (s *Stream) Close() error {
    close(s.readChunks)
    close(s.readSig)
    close(s.readTimeoutSig)

    return s.writeClose()
}

func (s *Stream) acceptReadSignals() {
    go func() {
        if s.readTimeout == 0 {
            for range s.readSig {}
            return
        }

        for {
            select {
            case <- time.After(s.readTimeout):
                s.readErrs <- TimeoutError
                s.readTimeoutSig <- struct{}{}
                return
            case <- s.readSig:
                fmt.Println("READ SIG")
            }
        }
    }()
}

func (s *Stream) feed(msg sppp.Message) {
    s.readChunks <- msg
    s.readSig <- struct{}{}
}

func (s *Stream) writeClose() error {
    msg := sppp.NewMessage(s.msgID, sppp.EndType, nil)
    return s.write(msg)
}

func (s *Stream) WriteData(b []byte) error {
    _, err := s.Write(b)
    return err
}

func (s *Stream) Write(b []byte) (int, error) {
    msgs := sppp.SplitIntoMessages(s.msgID, sppp.StreamType, b)

    for _, m := range msgs {
        err := s.write(m)
        if err != nil {
            return 0, err
        }
    }

    return len(b), nil
}

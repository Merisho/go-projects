package sppp

import (
    "io"
    "time"
)

type ReadStream interface {
    io.ReadCloser
    ReadData() ([]byte, error)
}

func newReadStream(msgID int64, readTimeout time.Duration) *readStream {
    s := &readStream{
        msgID:           msgID,
        readChunks:      make(chan Message, 1024),
        readErrs:        make(chan error),
        readSig:         make(chan struct{}),
        timeoutOccurred: false,
        readTimeout:     readTimeout,
        readTimeoutSig:  make(chan struct{}),
    }

    s.acceptReadSignals()

    return s
}

type readStream struct {
    msgID      int64
    readChunks chan Message
    readErrs   chan error
    readSig    chan struct{}
    readTimeoutSig    chan struct{}
    timeoutOccurred bool
    readTimeout time.Duration
}

func (s *readStream) ReadData() ([]byte, error) {
    b := make([]byte, 2048)
    n, err := s.Read(b)
    if err != nil {
        return nil, err
    }

    return b[:n], nil
}

func (s *readStream) Read(b []byte) (int, error) {
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

func (s *readStream) ReadTimeoutWait() chan struct{} {
    return s.readTimeoutSig
}

func (s *readStream) Close() error {
    close(s.readChunks)
    close(s.readSig)
    close(s.readTimeoutSig)

    return nil
}

func (s *readStream) acceptReadSignals() {
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
            }
        }
    }()
}

func (s *readStream) feed(msg Message) {
    s.readChunks <- msg
    s.readSig <- struct{}{}
}

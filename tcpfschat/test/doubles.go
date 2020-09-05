package test

import (
    "io"
    "net"
    "time"
)

func NewTestConnection() *TestConnection {
    return &TestConnection{
        readChunks:    make(chan string, 4096),
        writtenChunks: make(chan string, 4096),
        closed: make(chan struct{}),
    }
}

type TestConnection struct {
    readChunks    chan string
    writtenChunks chan string
    closed        chan struct{}
}

func (c *TestConnection) ChunksToRead(chunks ...string) {
    for _, b := range chunks {
        c.readChunks <- b
    }
}

func (c *TestConnection) EOFOnRead() *TestConnection {
    close(c.readChunks)
    return c
}

func (c *TestConnection) FrontWrittenChunk() string {
    select {
    case m := <-c.writtenChunks:
        return m
    default:
        return ""
    }
}

func (c *TestConnection) Read(b []byte) (int, error) {
    chunk, ok := <-c.readChunks
    if !ok {
        return 0, io.EOF
    }

    copy(b, chunk)
    return len(chunk), nil
}

func (c *TestConnection) Write(b []byte) (int, error) {
    c.writtenChunks <- string(b)
    return len(b), nil
}

func (c *TestConnection) Close() error {
    close(c.closed)
    return nil
}

func (c *TestConnection) Closed() bool {
    select {
    case <- c.closed:
        return true
    default:
        return false
    }
}

func (c *TestConnection) LocalAddr() net.Addr {
    return nil
}

func (c *TestConnection) RemoteAddr() net.Addr {
    return nil
}

func (c *TestConnection) SetDeadline(t time.Time) error {
    return nil
}

func (c *TestConnection) SetReadDeadline(t time.Time) error {
    return nil
}

func (c *TestConnection) SetWriteDeadline(t time.Time) error {
    return nil
}

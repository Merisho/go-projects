package test

import (
    "net"
    "time"
)

func NewTestConnection() *TestConnection {
    return &TestConnection{
        readChunks:    make(chan string, 4096),
        writtenChunks: make(chan string, 4096),
    }
}

type TestConnection struct {
    readChunks    chan string
    writtenChunks chan string
}

func (c *TestConnection) ChunksToRead(chunks ...string) {
    for _, b := range chunks {
        c.readChunks <- b
    }
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
    chunk := <-c.readChunks
    copy(b, chunk)
    return len(chunk), nil
}

func (c *TestConnection) Write(b []byte) (int, error) {
    c.writtenChunks <- string(b)
    return len(b), nil
}

func (c *TestConnection) Close() error {
    return nil
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

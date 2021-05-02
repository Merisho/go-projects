package server

import (
    "errors"
    "github.com/merisho/tcp-fs-chat/sppp"
    "net"
    "sync"
    "time"
)

var (
    TimeoutError = errors.New("timeout")
)

func NewConn(c net.Conn) *Conn {
    conn := &Conn{
        Conn:          c,
        textMsgChan:   make(chan sppp.Message, 1024),
        streamsChan:   make(chan sppp.Message, 1024),
        streams:       make(map[int64]Stream),
        msgReadTimeout: 5 * time.Second,
        txtBuffer: NewTextMessageBuffer(),
    }

    conn.startReading()
    conn.startTimeoutsHandling()

    return conn
}

type Conn struct {
    net.Conn
    textMsgChan   chan sppp.Message
    msgReadTimeout time.Duration
    streamReadTimeout time.Duration

    txtBuffer *TextMessageBuffer

    streamsMutex sync.Mutex
    streamsChan chan sppp.Message
    streams map[int64]Stream
}

func (c *Conn) ReadMsg() (sppp.Message, error) {
    return <- c.textMsgChan, nil
}

func (c *Conn) ReadStream() (chan []byte, chan error) {
    chunks := make(chan []byte, 1)
    errs := make(chan error, 1)

    msg := <- c.streamsChan

    go func() {
        loop:
        for {
            select {
            case m, ok := <- c.streams[msg.ID].stream:
                if !ok {
                    break loop
                }

                chunks <- m.Content
            case err := <- c.streams[msg.ID].errors:
                errs <- err
                break loop
            }
        }

        close(chunks)
        close(errs)
    }()

    return chunks, errs
}

func (c *Conn) SetMessageReadTimeout(d time.Duration) {
    c.msgReadTimeout = d
}

func (c *Conn) SetStreamReadTimeout(d time.Duration) {
    c.streamReadTimeout = d
}

func (c *Conn) startReading() {
    go func() {
        var b [1024]byte
        for {
            _, err := c.Conn.Read(b[:])
            if err != nil {
                panic(err)
            }

            msg, err := sppp.UnmarshalMessage(b)
            if err != nil {
                panic(err)
            }

            switch msg.Type {
            case sppp.TextType:
                c.handleTextMessage(msg)
            case sppp.EndType:
                c.handleMessageEnd(msg)
            case sppp.StreamType:
                c.handleStreamMessage(msg)
            }
        }
    }()
}

func (c *Conn) startTimeoutsHandling() {
    go func() {
        txtMsgTimeouts := c.txtBuffer.Timeouts()
        for {
            select {
            case id := <- txtMsgTimeouts:
                c.writeTimeout(id)
            }
        }
    }()
}

func (c *Conn) handleTextMessage(msg sppp.Message) {
    c.txtBuffer.Message(msg, c.msgReadTimeout)
}

func (c *Conn) handleStreamMessage(msg sppp.Message) {
    s, ok := c.streams[msg.ID]
    if ok {
        s.stream <- msg
    } else {
        c.deleteStreamAfterTimeout(msg)
        c.streams[msg.ID] = Stream{
            stream: make(chan sppp.Message, 1024),
            errors: make(chan error, 1024),
        }
        c.streams[msg.ID].stream <- msg
        c.streamsChan <- msg
    }
}

func (c *Conn) handleMessageEnd(msg sppp.Message) {
    m := c.txtBuffer.EndMessage(msg)
    if !m.Empty() {
        c.textMsgChan <- m
        return
    }

    c.tryRemoveStream(msg)
}

func (c *Conn) tryRemoveStream(msg sppp.Message) {
    c.streamsMutex.Lock()
    defer c.streamsMutex.Unlock()

    s, ok := c.streams[msg.ID]
    if ok {
        s.Close()
        delete(c.streams, msg.ID)
    }
}

func (c *Conn) deleteStreamAfterTimeout(msg sppp.Message) {
    if c.streamReadTimeout == 0 {
        return
    }

    go func() {
        panic("implement stream timeouts")
    }()
}

func (c *Conn) writeTimeout(id int64) {
    timeoutMsg := sppp.NewMessage(id, sppp.TimeoutType, nil)
    rawTimeoutMsg := timeoutMsg.Marshal()

    _, _ = c.Conn.Write(rawTimeoutMsg[:])
}

type Stream struct {
    stream chan sppp.Message
    errors chan error
}

func (s Stream) Close() {
    close(s.stream)
    close(s.errors)
}

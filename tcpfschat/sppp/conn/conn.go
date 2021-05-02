package conn

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
        Conn:           c,
        textMsgChan:    make(chan sppp.Message, 1024),
        newStreamsChan: make(chan Stream, 1024),
        streams:        make(map[int64]Stream),
        msgReadTimeout: 5 * time.Second,
        txtBuffer:      NewTextMessageBuffer(),
    }

    conn.startReading()
    conn.startTimeoutsHandling()

    return conn
}

type Conn struct {
    net.Conn

    textMsgChan   chan sppp.Message
    msgReadTimeout time.Duration
    txtBuffer *TextMessageBuffer

    streamReadTimeout time.Duration
    streamsMutex   sync.Mutex
    newStreamsChan chan Stream
    streams        map[int64]Stream
}

func (c *Conn) ReadMsg() (sppp.Message, error) {
    return <- c.textMsgChan, nil
}

func (c *Conn) ReadStream() (chan []byte, chan error) {
    chunks := make(chan []byte)
    errs := make(chan error)

    stream := <- c.newStreamsChan

    go func() {
        loop:
        for {
            select {
            case m, ok := <- stream.stream:
                if !ok {
                    break loop
                }

                chunks <- m.Content
            case err, ok := <- stream.errors:
                if ok {
                    errs <- err
                }

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
                badMsg := sppp.NewMessage(0, sppp.ErrorType, nil).Marshal()
                _, _ = c.Conn.Write(badMsg[:])
                continue
            }

            switch msg.Type {
            case sppp.TextType:
                c.handleTextMessage(msg)
            case sppp.StreamType:
                c.handleStreamMessage(msg)
            case sppp.EndType:
                c.handleMessageEnd(msg)
            }
        }
    }()
}

func (c *Conn) startTimeoutsHandling() {
    go func() {
        txtMsgTimeouts := c.txtBuffer.Timeouts()
        for id := range txtMsgTimeouts {
            c.writeTimeout(id)
        }
    }()
}

func (c *Conn) handleTextMessage(msg sppp.Message) {
    c.txtBuffer.Message(msg, c.msgReadTimeout)
}

func (c *Conn) handleStreamMessage(msg sppp.Message) {
    c.streamsMutex.Lock()
    defer c.streamsMutex.Unlock()

    s, ok := c.streams[msg.ID]
    if ok {
        s.stream <- msg
    } else {
        s := Stream{
            msgID:  msg.ID,
            stream: make(chan sppp.Message, 1024),
            errors: make(chan error),
            sig: make(chan struct{}, 1024),
        }
        c.deleteStreamAfterTimeout(s)
        s.stream <- msg
        c.newStreamsChan <- s

        c.streams[msg.ID] = s
    }
}

func (c *Conn) handleMessageEnd(msg sppp.Message) {
    m := c.txtBuffer.EndMessage(msg)
    if !m.Empty() {
        c.textMsgChan <- m
        return
    }

    c.removeStream(msg.ID)
}

func (c *Conn) deleteStreamAfterTimeout(s Stream) {
    if c.streamReadTimeout == 0 {
        return
    }

    go func() {
        for {
            select {
            case <- time.After(c.streamReadTimeout):
                s.errors <- TimeoutError
                c.writeTimeout(s.msgID)
                c.removeStream(s.msgID)
            case _, ok := <- s.sig:
                if !ok {
                    return
                }
            }
        }
    }()
}

func (c *Conn) removeStream(msgID int64) {
    c.streamsMutex.Lock()
    defer c.streamsMutex.Unlock()

    s, ok := c.streams[msgID]
    if ok {
        s.Close()
        delete(c.streams, msgID)
    }
}

func (c *Conn) writeTimeout(id int64) {
    timeoutMsg := sppp.NewMessage(id, sppp.TimeoutType, nil)
    rawTimeoutMsg := timeoutMsg.Marshal()

    _, _ = c.Conn.Write(rawTimeoutMsg[:])
}

type Stream struct {
    msgID  int64
    stream chan sppp.Message
    errors chan error
    sig    chan struct{}
}

func (s Stream) Close() {
    close(s.stream)
    close(s.errors)
    close(s.sig)
}

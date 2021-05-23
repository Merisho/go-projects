package sppp

import (
    "errors"
    "github.com/merisho/tcp-fs-chat/internal/csdt"
    "log"
    "math/rand"
    "net"
    "sync"
    "time"
)

var (
    TimeoutError = errors.New("timeout")
)

func NewConn(c net.Conn) *Conn {
    conn := &Conn{
        Conn:               c,
        textMsgChan:        make(chan Message, 1024),
        mainErrChan:        make(chan error),
        newReadStreamsChan: make(chan *readStream, 128),
        readStreams:        make(map[uint64]*readStream),
        msgReadTimeout:     5 * time.Second,
        txtBuffer:          NewTextMessageBuffer(),
        rand:               rand.New(rand.NewSource(time.Now().Unix())),
        finishErr:          csdt.NewValue(),
    }

    conn.startReading()
    conn.startTimeoutsHandling()

    return conn
}

type Conn struct {
    net.Conn

    textMsgChan    chan Message
    msgReadTimeout time.Duration
    txtBuffer *TextMessageBuffer

    streamReadTimeout time.Duration
    streamsMutex       sync.Mutex
    newReadStreamsChan chan *readStream
    readStreams        map[uint64]*readStream
    rand               *rand.Rand

    mainErrChan chan error
    finishErr   *csdt.Value
}

func (c *Conn) ReadMsg() (Message, error) {
    finishErr := c.finishErr.Value()
    if finishErr != nil {
        return Message{}, finishErr.(error)
    }

    select {
    case m := <- c.textMsgChan:
        return m, nil
    case <- c.mainErrChan:
        return Message{}, c.finishErr.Value().(error)
    }
}

func (c *Conn) MsgCount() int {
    return len(c.textMsgChan)
}

func (c *Conn) ReadStream() (ReadStream, error) {
    finishErr := c.finishErr.Value()
    if finishErr != nil {
        return nil, finishErr.(error)
    }

    select {
    case s := <- c.newReadStreamsChan:
        return s, nil
    case <- c.mainErrChan:
        return nil, c.finishErr.Value().(error)
    }
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
                c.finishErr.SetValue(err)
                close(c.mainErrChan)
                go c.closeReadStreams()
                return
            }

            msg, err := UnmarshalMessage(b)
            if err != nil {
                badMsg := NewMessage(0, ErrorType, nil).Marshal()
                _, _ = c.Conn.Write(badMsg[:])
                continue
            }

            switch msg.Type {
            case TextType:
                c.handleTextMessage(msg)
            case StreamType:
                c.handleStreamMessage(msg)
            case EndType:
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

func (c *Conn) handleTextMessage(msg Message) {
    c.txtBuffer.Message(msg, c.msgReadTimeout)
}

func (c *Conn) handleStreamMessage(msg Message) {
    c.streamsMutex.Lock()
    defer c.streamsMutex.Unlock()

    s, ok := c.readStreams[msg.ID]
    if ok {
        s.feed(msg)
    } else {
        s := newReadStream(msg.ID, msg.Content, c.streamReadTimeout)

        c.deleteReadStreamAfterTimeout(s)
        c.newReadStreamsChan <- s

        c.readStreams[msg.ID] = s
    }
}

func (c *Conn) handleMessageEnd(msg Message) {
    m := c.txtBuffer.EndMessage(msg)
    if !m.Empty() {
        c.textMsgChan <- m
        return
    }

    c.removeReadStream(msg.ID)
}

func (c *Conn) deleteReadStreamAfterTimeout(s *readStream) {
    go func() {
        _, ok := <- s.ReadTimeoutWait()
        if !ok {
            return
        }
        
        c.writeTimeout(s.msgID)
        c.removeReadStream(s.msgID)
    }()
}

func (c *Conn) removeReadStream(msgID uint64) {
    c.streamsMutex.Lock()
    defer c.streamsMutex.Unlock()

    s, ok := c.readStreams[msgID]
    if ok {
        err := s.Close()
        if err != nil {
            log.Fatalf("Could not close stream: %s", err)
        }

        delete(c.readStreams, msgID)
    }
}

func (c *Conn) writeTimeout(id uint64) {
    timeoutMsg := NewMessage(id, TimeoutType, nil)
    rawTimeoutMsg := timeoutMsg.Marshal()

    _, _ = c.Conn.Write(rawTimeoutMsg[:])
}

func (c *Conn) WriteMsg(rawMsg []byte) error {
    id := c.rand.Uint64()
    msgs := SplitIntoMessages(id, TextType, rawMsg)
    msgs = append(msgs, NewMessage(id, EndType, nil))

    for _, m := range msgs {
        rawMsg := m.Marshal()
        _, err := c.Write(rawMsg[:])
        if err != nil {
            return err
        }
    }

    return nil
}

func (c *Conn) WriteStream(meta []byte) (WriteStream, error) {
    id := c.rand.Uint64()

    s := newWriteStream(id, c)

    err := s.WriteData(meta)
    if err != nil {
        return nil, err
    }

    return s, nil
}

func (c *Conn) closeReadStreams() {
    c.streamsMutex.Lock()
    defer c.streamsMutex.Unlock()

    for _, s := range c.readStreams {
        _ = s.Close()
    }

    c.readStreams = nil
}

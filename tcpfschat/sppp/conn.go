package sppp

import (
    "errors"
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
        Conn:           c,
        textMsgChan:    make(chan Message, 1024),
        newStreamsChan: make(chan *Stream, 128),
        streams:        make(map[int64]*Stream),
        msgReadTimeout: 5 * time.Second,
        txtBuffer:      NewTextMessageBuffer(),
        rand:           rand.New(rand.NewSource(time.Now().Unix())),
    }

    conn.startReading()
    conn.startTimeoutsHandling()

    return conn
}

type Conn struct {
    net.Conn

    textMsgChan   chan Message
    msgReadTimeout time.Duration
    txtBuffer *TextMessageBuffer

    streamReadTimeout time.Duration
    streamsMutex   sync.Mutex
    newStreamsChan chan *Stream
    streams        map[int64]*Stream
    rand           *rand.Rand
}

func (c *Conn) ReadMsg() (Message, error) {
    return <- c.textMsgChan, nil
}

func (c *Conn) ReadStream() ReadStream {
    return <- c.newStreamsChan
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

    s, ok := c.streams[msg.ID]
    if ok {
        s.feed(msg)
    } else {
        s := NewStream(msg.ID, c.streamReadTimeout)
        s.feed(msg)

        c.deleteStreamAfterTimeout(s)
        c.newStreamsChan <- s

        c.streams[msg.ID] = s
    }
}

func (c *Conn) handleMessageEnd(msg Message) {
    m := c.txtBuffer.EndMessage(msg)
    if !m.Empty() {
        c.textMsgChan <- m
        return
    }

    c.removeStream(msg.ID)
}

func (c *Conn) deleteStreamAfterTimeout(s *Stream) {
    go func() {
        _, ok := <- s.ReadTimeoutWait()
        if !ok {
            return
        }
        
        c.writeTimeout(s.msgID)
        c.removeStream(s.msgID)
    }()
}

func (c *Conn) removeStream(msgID int64) {
    c.streamsMutex.Lock()
    defer c.streamsMutex.Unlock()

    s, ok := c.streams[msgID]
    if ok {
        err := s.Close()
        if err != nil {
            log.Fatalf("Could not close stream: %s", err)
        }

        delete(c.streams, msgID)
    }
}

func (c *Conn) writeTimeout(id int64) {
    timeoutMsg := NewMessage(id, TimeoutType, nil)
    rawTimeoutMsg := timeoutMsg.Marshal()

    _, _ = c.Conn.Write(rawTimeoutMsg[:])
}

func (c *Conn) WriteMsg(rawMsg []byte) error {
    id := c.rand.Int63()
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
    id := c.rand.Int63()

    s := NewStream(id, c.streamReadTimeout)
    s.write = func(message Message) error {
        raw := message.Marshal()
        _, err := c.Write(raw[:])
        return err
    }

    err := s.WriteData(meta)
    if err != nil {
        return nil, err
    }

    return s, nil
}

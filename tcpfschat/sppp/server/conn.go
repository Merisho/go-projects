package server

import (
    "github.com/merisho/tcp-fs-chat/sppp"
    "net"
    "sync"
    "time"
)

func NewConn(c net.Conn) *Conn {
    conn := &Conn{
        Conn:          c,
        textMsgBuffer: make(map[int64]sppp.Message),
        textMsgChan:   make(chan sppp.Message, 1024),
        msgReadTimeout: 5 * time.Second,
    }

    conn.startReading()

    return conn
}

type Conn struct {
    net.Conn
    textMsgBufferMutex sync.Mutex
    textMsgBuffer map[int64]sppp.Message
    textMsgChan   chan sppp.Message
    msgReadTimeout time.Duration
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
            case sppp.MsgEndType:
                c.handleMessageEnd(msg)
            }
        }
    }()
}

func (c *Conn) ReadMsg() (sppp.Message, error) {
    return <- c.textMsgChan, nil
}

func (c *Conn) SetMessageReadTimeout(d time.Duration) {
    c.msgReadTimeout = d
}

func (c *Conn) handleTextMessage(msg sppp.Message) {
    c.textMsgBufferMutex.Lock()
    defer c.textMsgBufferMutex.Unlock()

    m, ok := c.textMsgBuffer[msg.ID]
    if ok {
        m.Content = append(m.Content, msg.Content[:msg.Size]...)
        c.textMsgBuffer[msg.ID] = m
    } else {
        c.deleteAfterTimeout(msg)
        c.textMsgBuffer[msg.ID] = msg
    }
}

func (c *Conn) handleMessageEnd(msg sppp.Message) {
    c.textMsgBufferMutex.Lock()
    defer c.textMsgBufferMutex.Unlock()

    m, ok := c.textMsgBuffer[msg.ID]
    if ok {
        delete(c.textMsgBuffer, msg.ID)
        c.textMsgChan <- m
    }
}

func (c *Conn) deleteAfterTimeout(msg sppp.Message) {
    go func() {
        time.Sleep(c.msgReadTimeout)

        c.textMsgBufferMutex.Lock()
        defer c.textMsgBufferMutex.Unlock()

        timeoutMsg := sppp.NewMessage(msg.ID, sppp.TimeoutType, nil)
        rawTimeoutMsg := timeoutMsg.Marshal()

        _, _ = c.Conn.Write(rawTimeoutMsg[:])
        delete(c.textMsgBuffer, msg.ID)

    }()
}

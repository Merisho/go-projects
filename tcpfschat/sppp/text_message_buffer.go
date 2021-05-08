package sppp

import (
    "sync"
    "time"
)

func NewTextMessageBuffer() *TextMessageBuffer {
    return &TextMessageBuffer{
        buffer: make(map[int64]Message),
        timeouts: make(chan int64),
    }
}

type TextMessageBuffer struct {
    buffer map[int64]Message
    mutex sync.Mutex
    timeouts chan int64
}

func (b *TextMessageBuffer) Message(msg Message, timeout time.Duration) {
    b.mutex.Lock()
    defer b.mutex.Unlock()

    m, ok := b.buffer[msg.ID]
    if ok {
        m.Content = append(m.Content, msg.Content[:msg.Size]...)
        m.Size += msg.Size
        b.buffer[msg.ID] = m
    } else {
        b.deleteBufferOnTimeout(msg, timeout)
        b.buffer[msg.ID] = msg
    }
}

func (b *TextMessageBuffer) EndMessage(msg Message) Message {
    b.mutex.Lock()
    m, ok := b.buffer[msg.ID]
    b.mutex.Unlock()
    if !ok {
        return Message{}
    }

    b.deleteBuffer(msg.ID)

    return m
}

func (b *TextMessageBuffer) Timeouts() <- chan int64 {
    return b.timeouts
}

func (b *TextMessageBuffer) deleteBufferOnTimeout(m Message, timeout time.Duration) {
    go func() {
        time.Sleep(timeout)

        if b.deleteBuffer(m.ID) {
            b.timeouts <- m.ID
        }
    }()
}

func (b *TextMessageBuffer) deleteBuffer(id int64) (deleted bool) {
    b.mutex.Lock()
    defer b.mutex.Unlock()

    if _, ok := b.buffer[id]; !ok {
        return false
    }

    delete(b.buffer, id)
    return true
}

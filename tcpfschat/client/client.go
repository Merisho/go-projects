package client

import (
	"github.com/merisho/tcp-fs-chat/client/messagebuffer"
	"github.com/merisho/tcp-fs-chat/internal/chaterrors"
	"net"
	"time"
)

const (
	maxMessageSize = 16 * 1024
)

func New(conn net.Conn, id []byte) Client {
	c := Client{
		conn: conn,
		id: id,
		receive: make(chan string, 1024),
	}

	<- c.readMessages()

	return c
}

type Client struct {
	conn      net.Conn
	id        []byte
	receive   chan string
}

func (c *Client) readMessages() chan struct{} {
	ready := make(chan struct{})

	msgBuffer := messagebuffer.MessageBuffer{}
	go func() {
		close(ready)
		for {
			b := make([]byte, maxMessageSize)
			n, err := c.conn.Read(b)
			if !chaterrors.IsTemporary(err) {
				close(c.receive)
				return
			}

			msgs := msgBuffer.Messages(b[:n])
			for _, msg := range msgs {
				if len(msg) > 0 {
					c.receive <- string(msg)
				}
			}
		}
	}()

	return ready
}

func (c *Client) Receive() chan string {
	return c.receive
}

func (c *Client) Send(msg string) error {
	if len(msg) < maxMessageSize {
		return c.send([]byte(msg))
	}

	for start := 0; start < len(msg); start += maxMessageSize {
		err := c.send([]byte(msg[start:start + maxMessageSize]))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) send(msg []byte) error {
	_ = c.conn.SetWriteDeadline(time.Now().Add(1 * time.Second))

	m := make([]byte, len(msg))
	copy(m, msg)
	m = append(m, 0)

	_, err := c.conn.Write(m)
	return err
}

func (c *Client) ID() []byte {
	return c.id
}

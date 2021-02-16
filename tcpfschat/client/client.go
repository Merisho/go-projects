package client

import (
	"bytes"
	"github.com/merisho/tcp-fs-chat/internal/chaterrors"
	"net"
	"time"
)

const (
	maxMessageSize = 8
)

func New(conn net.Conn, id []byte) Client {
	c := Client{
		conn: conn,
		id: id,
		receive: make(chan string),
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

	go func() {
		close(ready)
		for {
			b := make([]byte, maxMessageSize)
			n, err := c.conn.Read(b)
			if !chaterrors.IsTemporary(err) {
				close(c.receive)
				return
			}

			msgs := bytes.Split(b[:n], []byte{0})
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

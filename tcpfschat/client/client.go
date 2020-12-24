package client

import (
	"bytes"
	"github.com/merisho/tcp-fs-chat/internal/chaterrors"
	"io"
)

func New(conn io.ReadWriteCloser, id []byte) Client {
	c := Client{
		conn: conn,
		id: id,
		receive: make(chan string),
	}

	<- c.readMessages()

	return c
}

type Client struct {
	conn      io.ReadWriteCloser
	id        []byte
	receive   chan string
}

func (c *Client) readMessages() chan struct{} {
	ready := make(chan struct{})

	go func() {
		close(ready)
		for {
			b := make([]byte, 1024 * 1024)
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
	_, err := c.conn.Write(append([]byte(msg), 0))
	return err
}

func (c *Client) ID() []byte {
	return c.id
}

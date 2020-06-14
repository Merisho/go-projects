package client

import (
	"fmt"
	"net"
	"strings"
)

func Connect(host string, port uint16) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

type Client struct {
	conn net.Conn
}

func (c *Client) Receive() chan string {
	msgs := make(chan string, 4096)

	go func() {
		for {
			b := make([]byte, 1024)
			_, _ = c.conn.Read(b)
			msgs <- strings.TrimRight(string(b), "\x00")
		}
	}()

	return msgs
}

func (c *Client) Send(msg string) {
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
}

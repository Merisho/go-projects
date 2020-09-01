package client

import (
	"io"
	"log"
	"time"
)

const (
	sendMessageCmd commandType = iota
	addReceiverCmd
	removeReceiverCmd
)

func New(conn io.ReadWriteCloser) Client {
	c := Client{
		conn: conn,
		receiversCommands: make(chan command),
	}

	c.commandHandlers = map[commandType]func(command){
		sendMessageCmd: c.sendMessage,
		addReceiverCmd: c.addReceiver,
		removeReceiverCmd: c.removeReceiver,
	}

	<- c.handleReceiversCommands()
	<- c.readMessages()

	return c
}

type command struct {
	cmdType commandType
	payload interface{}
}

type commandType int

type Client struct {
	conn      io.ReadWriteCloser
	receivers []chan string
	receiversCommands chan command
	commandHandlers map[commandType]func(command)
}

func (c *Client) handleReceiversCommands() chan struct{} {
	ready := make(chan struct{})

	go func() {
		close(ready)
		for cmd := range c.receiversCommands {
			c.commandHandlers[cmd.cmdType](cmd)
		}
	}()

	return ready
}

func (c *Client) removeReceiver(cmd command) {
	r := cmd.payload.(chan string)
	var i int
	for i = range c.receivers {
		if c.receivers[i] == r {
			break
		}
	}
	close(r)
	c.receivers = append(c.receivers[:i], c.receivers[i + 1:]...)
}

func (c *Client) addReceiver(cmd command) {
	r := cmd.payload.(chan string)
	c.receivers = append(c.receivers, r)
	r <- "ok"
}

func (c *Client) sendMessage(cmd command) {
	msg := cmd.payload.(string)
	for _, r := range c.receivers {
		select {
		case r <- msg:
		default:
			c.retryMsg(msg, r)
		}
	}
}

func (c *Client) retryMsg(msg string, r chan string) {
	go func() {
		select {
		case r <- msg:
		case <- time.After(10 * time.Millisecond):
			c.receiversCommands <- command{removeReceiverCmd, r}
		}
	}()
}

func (c *Client) readMessages() chan struct{} {
	ready := make(chan struct{})

	go func() {
		close(ready)
		for {
			b := make([]byte, 8192)
			n, err := c.conn.Read(b)
			if err != nil {
				log.Println(err)
				continue
			}

			c.receiversCommands <- command{sendMessageCmd, b[:n]}
		}
	}()

	return ready
}

func (c *Client) Receive() chan string {
	r := make(chan string)
	c.receiversCommands <- command{addReceiverCmd, r}
	<-r
	return r
}

func (c *Client) Send(msg string) error {
	_, err := c.conn.Write([]byte(msg))
	return err
}

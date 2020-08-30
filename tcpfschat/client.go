package main

import (
	"bufio"
	"fmt"
	"github.com/merisho/tcp-fs-chat/client"
	"os"
)

func main() {
	c, err := client.ConnectTCP("localhost", 1337)
	if err != nil {
		panic(err)
	}

	go func() {
		msgs := c.Receive()
		for msg := range msgs {
			fmt.Println(msg)
		}
	}()

	r := bufio.NewReader(os.Stdin)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			panic(err)
		}

		err = c.Send(string(b))
		if err != nil {
			panic(err)
		}
	}
}

func readLine() string {
	r := bufio.NewReader(os.Stdin)
	b, _, err := r.ReadLine()
	if err != nil {
		panic(err)
	}

	return string(b)
}

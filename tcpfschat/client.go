package main

import (
	"bufio"
	"fmt"
	"github.com/merisho/tcp-fs-chat/client"
	"os"
)

type obj struct {}

func main() {
	c, err := client.Connect("localhost", 1337)
	if err != nil {
		panic(err)
	}

	username := readLine()
	password := readLine()
	err = c.Auth(username, password)
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

		c.Send(string(b))
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

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

	go func() {
		r := c.Receive()
		for msg := range r {
			fmt.Println(msg)
		}
	}()

	r := bufio.NewReader(os.Stdin)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			panic(err)
		}

		fmt.Println(b)

		c.Send(string(b))
	}
}

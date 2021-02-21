package main

import (
	"bufio"
	"fmt"
	"github.com/merisho/tcp-fs-chat/client"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("host argument is missing")
	}

	host := os.Args[1]

	port := uint16(1337)
	if len(os.Args) >= 3 {
		p, err := strconv.ParseUint(os.Args[2], 10, 16)
		if err != nil {
			log.Fatal("port is not a number")
		}

		port = uint16(p)
	}

	c, err := client.ConnectTCP(host, port)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) >= 4 {
		c.SetName(os.Args[3])
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
			log.Fatal(err)
		}

		err = c.Send(string(b))
		if err != nil {
			log.Fatal(err)
		}
	}
}

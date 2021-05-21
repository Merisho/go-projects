package main

import (
	"bufio"
	"fmt"
	"github.com/merisho/tcp-fs-chat/sppp"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("host argument is missing")
	}

	host := os.Args[1]
	port := getPort()

	d, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}

	c := sppp.NewConn(d)

	go func() {
		for {
			msg, err := c.ReadMsg()
			if err != nil {
				if err == io.EOF {
					log.Println("Connection closed")
					os.Exit(0)
				}

				log.Fatal(err)
			}

			now := time.Now().Format(time.RFC822Z)
			fmt.Println(now, string(msg.Content))
		}
	}()

	go func() {
		for {
			s, err := c.ReadStream()
			if err != nil {
				if err == io.EOF {
					log.Println("Connection closed")
					os.Exit(0)
				}

				log.Fatal(err)
			}

			go handleStream(s)
		}
	}()

	fmt.Println("Your nickname:")
	processOutgoingMessages(c)
}

func processOutgoingMessages(c *sppp.Conn) {
	r := bufio.NewReader(os.Stdin)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		if string(b[:2]) == "/f" {
			go sendFile(c, string(b[2:]))
		} else {
			err = c.WriteMsg(b)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func sendFile(c *sppp.Conn, filePath string) {
	filePath = strings.TrimSpace(filePath)

	fileName := filepath.Base(filePath)
	fmt.Println("Sending file", fileName)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	ws, err := c.WriteStream([]byte(fileName))
	if err != nil {
		log.Fatalf("Could not send file: %s", err)
	}
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Fatalf("Could not send write stream: %s", err)
		}
	}()

	b := make([]byte, 8192)
	n, err := f.Read(b)
	for err == nil {
		err = ws.WriteData(b[:n])
		if err != nil {
			log.Fatal(err)
		}

		n, err = f.Read(b)
	}

	if err != io.EOF {
		log.Fatal(err)
	}
}

func handleStream(s sppp.ReadStream) {
	meta, err := s.ReadData()
	if err != nil {
		log.Fatal(err)
	}

	fileName := string(meta)
	fmt.Println("Client accepting:", fileName)

	f, err := os.OpenFile(fileName, os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	b, err := s.ReadData()
	for err == nil {
		_, err = f.Write(b)
		if err != nil {
			log.Fatal(err)
		}

		b, err = s.ReadData()
	}

	if err != io.EOF {
		log.Fatal(err)
	}

	fmt.Println("Download finished:", fileName)
}

func getPort() uint16 {
	port := uint16(1337)
	if len(os.Args) >= 3 {
		p, err := strconv.ParseUint(os.Args[2], 10, 16)
		if err != nil {
			log.Fatal("port is not a number")
		}

		port = uint16(p)
	}

	return port
}

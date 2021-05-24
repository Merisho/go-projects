package main

import (
	"bufio"
	"fmt"
	"github.com/merisho/tcp-fs-chat/sppp"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	host := "localhost"
	if len(os.Args) >= 2 {
		host = os.Args[1]
	}

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
					fmt.Println("Connection closed")
					os.Exit(0)
				}

				log.Fatal(err)
			}

			printMsg(string(msg))
		}
	}()

	go func() {
		for {
			s, err := c.ReadStream()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Connection closed")
					os.Exit(0)
				}

				log.Fatal(err)
			}

			go handleStream(s)
		}
	}()

	fmt.Printf("To send a file:\n\t/f <full path>\n\n")

	fmt.Println("Your nickname:")
	go processOutgoingMessages(c)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)

	<- s

	err = c.Close()
	if err != nil {
		log.Println(err)
	}
}

func processOutgoingMessages(c *sppp.Conn) {
	r := bufio.NewReader(os.Stdin)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			printServiceMsg("Could not read stdin:", err.Error())
		}

		if string(b[:2]) == "/f" {
			go sendFile(c, string(b[2:]))
		} else {
			err = c.WriteMsg(b)
			if err != nil {
				if err == io.ErrClosedPipe {
					return
				}

				log.Fatal(err)
			}
		}
	}
}

func sendFile(c *sppp.Conn, filePath string) {
	filePath = strings.TrimSpace(filePath)

	fileName := filepath.Base(filePath)
	printServiceMsg("Sending file", fileName)

	f, err := os.Open(filePath)
	if err != nil {
		printServiceMsg("Cannot open the file:", filePath, err.Error())
	}
	defer f.Close()

	ws, err := c.WriteStream([]byte(fileName))
	if err != nil {
		printServiceMsg("Could not send file:", err.Error())
		return
	}
	defer func() {
		err := ws.Close()
		if err != nil {
			printServiceMsg("Could not close write stream:", err.Error())
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
		printServiceMsg("Could not send file successfully:", err.Error())
	}
}

func handleStream(s sppp.ReadStream) {
	fileName := string(s.Meta())
	printServiceMsg("Client accepting:", fileName)

	f, err := os.OpenFile(fileName, os.O_CREATE, 0777)
	if err != nil {
		printServiceMsg("Cannot write to file:", fileName, err.Error())
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			printServiceMsg("Error writing to a file:", fileName, err.Error())
		}
	}()

	b, err := s.ReadData()
	for err == nil {
		_, err = f.Write(b)
		if err != nil {
			printServiceMsg("Error writing to a file:", fileName, err.Error())
			return
		}

		b, err = s.ReadData()
	}

	if err != io.EOF {
		printServiceMsg("Error downloading a file:", fileName, err.Error())
	}

	printServiceMsg("Download finished:", fileName)
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

func printServiceMsg(m ...string) {
	printMsg(append([]string{"========== "}, m...)...)
}

func printMsg(m ...string) {
	fmt.Printf("[%s] %s\n", time.Now().Format(time.Kitchen), strings.Join(m, " "))
}

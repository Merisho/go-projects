package main

import (
	"fmt"
	"github.com/merisho/tcp-fs-chat/sppp"
	"io"
	"log"
	"net"
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

	d, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}

	c := sppp.NewConn(d)


	go func() {
		for {
			s := c.ReadStream()

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
		}
	}()

	s, err := c.WriteStream([]byte("file.mp4"))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("d:\\Memories\\Feeling Good by Me.mp4")
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 2048)

	n, err := f.Read(b)
	for err == nil {
		err = s.WriteData(b[:n])
		if err != nil {
			break
		}

		n, err = f.Read(b)
	}

	if err != io.EOF {
		log.Fatal(err)
	}

	err = s.Close()
	if err != nil {
		log.Fatal(err)
	}

	//r := bufio.NewReader(os.Stdin)
	//for {
	//	b, _, err := r.ReadLine()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	err = c.WriteMsg(b)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
}

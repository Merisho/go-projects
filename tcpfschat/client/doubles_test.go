package client_test

import (
	"net"
	"strconv"
	"sync"
	"sync/atomic"
)

var port = portGen()

func StartFakeServer() (uint16, chan string) {
	p := port()

	var mu sync.Mutex
	var conns []net.Conn
	go func() {
		ln, err := net.Listen("tcp", ":" + strconv.FormatUint(uint64(p), 10))
		if err != nil {
			panic(err)
		}

		for {
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}

			mu.Lock()
			conns = append(conns, conn)
			mu.Unlock()
		}
	}()

	msgs := make(chan string)
	go func() {
		for msg := range msgs {
			mu.Lock()
			for _, c := range conns {
				_, err := c.Write([]byte(msg))
				if err != nil {
					panic(err)
				}
			}
			mu.Unlock()
		}
	}()

	return p, msgs
}

func portGen() func() uint16 {
	base := uint32(1336)
	return func() uint16 {
		if base == 65535 {
			panic("port out of range")
		}

		return uint16(atomic.AddUint32(&base, 1))
	}
}

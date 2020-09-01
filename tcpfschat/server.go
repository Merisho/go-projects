package main

import (
	"github.com/merisho/tcp-fs-chat/server"
	"log"
)

func main() {
	_, err := server.ServeTCP(1337)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}

package main

import "github.com/merisho/tcp-fs-chat/server"

func main() {
	_, err := server.Serve(1337)
	if err != nil {
		panic(err)
	}

	select {}
}

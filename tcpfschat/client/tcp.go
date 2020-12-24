package client

import (
    "fmt"
    "net"
)

func ConnectTCP(host string, port uint16) (Client, error) {
    addr := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return Client{}, err
    }

    id := make([]byte, 128)
    n, err := conn.Read(id)
    if err != nil {
        return Client{}, err
    }

    return New(conn, id[:n]), nil
}
